package bullet

// #cgo pkg-config: bullet
// #cgo windows LDFLAGS: -Wl,--allow-multiple-definition
// #include "bulletglue.h"
import "C"
import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/golang/glog"
)

func init() {
	core.SetPhysicsSystem(New())
}

// convenience
func vec3_to_bullet(vec mgl64.Vec3) (out C.plVector3) {
	out[0] = C.plReal(vec.X())
	out[1] = C.plReal(vec.Y())
	out[2] = C.plReal(vec.Z())

	return out
}

func quat_to_bullet(quat mgl64.Quat) (out C.plQuaternion) {
	out[0] = C.plReal(quat.X())
	out[1] = C.plReal(quat.Y())
	out[2] = C.plReal(quat.Z())
	out[3] = C.plReal(quat.W)

	return out
}

func mat4_to_bullet(mat mgl64.Mat4) (out [16]C.plReal) {
	for x := 0; x < 16; x++ {
		out[x] = C.plReal(mat[x])
	}
	return out
}

func mat4_from_bullet(mat [16]C.plReal) (out mgl64.Mat4) {
	for x := 0; x < 16; x++ {
		out[x] = float64(mat[x])
	}
	return out
}

type PhysicsSystem struct {
	sdk   C.plPhysicsSdkHandle
	world C.plDynamicsWorldHandle
}

func New() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (p *PhysicsSystem) Start() {
	glog.Info("Starting")

	// create an sdk handle
	p.sdk = C.plNewBulletSdk()

	// instance a world
	p.world = C.plCreateDynamicsWorld(p.sdk)
	C.plSetGravity(p.world, 0.0, 0.0, 0.0)
}

func (this *PhysicsSystem) Stop() {
	glog.Info("Stopping")

	C.plDeleteDynamicsWorld(this.world)
	C.plDeletePhysicsSdk(this.sdk)
}

func (this *PhysicsSystem) SetGravity(g mgl64.Vec3) {
	vec := vec3_to_bullet(g)
	C.plSetGravity(this.world, vec[0], vec[1], vec[2])
}

// fixme: remove gosg dependencies by passing a RigidBodyVec instead of NodeVec
func (this *PhysicsSystem) Update(dt float64, nodes []*core.Node) {
	for _, n := range nodes {
		n.RigidBody().SetTransform(n.WorldTransform())
	}
	C.plStepSimulation(this.world, C.plReal(dt))
	for _, n := range nodes {
		n.SetWorldTransform(n.RigidBody().GetTransform())
	}
}

func (this *PhysicsSystem) AddRigidBody(rigidBody core.RigidBody) {
	C.plAddRigidBody(this.world, rigidBody.(RigidBody).handle)
}

func (this *PhysicsSystem) RemoveRigidBody(rigidBody core.RigidBody) {
	C.plRemoveRigidBody(this.world, rigidBody.(RigidBody).handle)
}

func (this *PhysicsSystem) CreateRigidBody(mass float32, shape core.CollisionShape) core.RigidBody {
	body := C.plCreateRigidBody(nil, C.float(mass), shape.(CollisionShape).handle)
	r := RigidBody{body}
	return r
}

func (this *PhysicsSystem) DeleteRigidBody(body core.RigidBody) {
	C.plDeleteRigidBody(body.(RigidBody).handle)
}

func (this *PhysicsSystem) NewStaticPlaneShape(normal mgl64.Vec3, constant float64) core.CollisionShape {
	vec := vec3_to_bullet(normal)
	return CollisionShape{C.plNewStaticPlaneShape(&vec[0], C.float(constant))}
}

func (this *PhysicsSystem) NewSphereShape(radius float64) core.CollisionShape {
	return CollisionShape{C.plNewSphereShape(C.plReal(radius))}
}

func (this *PhysicsSystem) NewBoxShape(box mgl64.Vec3) core.CollisionShape {
	vec := vec3_to_bullet(box)
	return CollisionShape{C.plNewBoxShape(vec[0], vec[1], vec[2])}
}

func (this *PhysicsSystem) NewCapsuleShape(radius float64, height float64) core.CollisionShape {
	return CollisionShape{C.plNewCapsuleShape(C.plReal(radius), C.plReal(height))}
}

func (this *PhysicsSystem) NewConeShape(radius float64, height float64) core.CollisionShape {
	return CollisionShape{C.plNewConeShape(C.plReal(radius), C.plReal(height))}
}

func (this *PhysicsSystem) NewCylinderShape(radius float64, height float64) core.CollisionShape {
	return CollisionShape{C.plNewCylinderShape(C.plReal(radius), C.plReal(height))}
}

func (this *PhysicsSystem) NewCompoundShape() core.CollisionShape {
	return CollisionShape{C.plNewCompoundShape()}
}

func (this *PhysicsSystem) NewConvexHullShape() core.CollisionShape {
	return CollisionShape{C.plNewConvexHullShape()}
}

func (this *PhysicsSystem) NewStaticTriangleMeshShape(mesh core.Mesh) core.CollisionShape {
	/*
		bulletMeshInterface := C.plNewMeshInterface()

		// add triangles
		for v := 0; v < len(indices); v += 3 {
			i1 := indices[v+0]
			i2 := indices[v+1]
			i3 := indices[v+2]

			v1 := vec3_to_bullet(positions[i1*3])
			v2 := vec3_to_bullet(positions[i2*3])
			v3 := vec3_to_bullet(positions[i3*3])

			C.plAddTriangle(bulletMeshInterface, &v1[0], &v2[0], &v3[0])
		}

		return CollisionShape{C.plNewStaticTriangleMeshShape(bulletMeshInterface)}
	*/
	return nil
}

func (this *PhysicsSystem) DeleteShape(shape core.CollisionShape) {
	C.plDeleteShape(shape.(CollisionShape).handle)
}

type RigidBody struct {
	handle C.plRigidBodyHandle
}

func (this RigidBody) GetTransform() mgl64.Mat4 {
	mat := mat4_to_bullet(mgl64.Ident4())
	C.plGetOpenGLMatrix(this.handle, &mat[0])
	return mat4_from_bullet(mat)
}

func (this RigidBody) SetTransform(transform mgl64.Mat4) {
	mat := mat4_to_bullet(transform)
	C.plSetOpenGLMatrix(this.handle, &mat[0])
}

func (this RigidBody) ApplyImpulse(impulse mgl64.Vec3, localPoint mgl64.Vec3) {
	i := vec3_to_bullet(impulse)
	p := vec3_to_bullet(localPoint)
	C.plApplyImpulse(this.handle, &i[0], &p[0])
}

type CollisionShape struct {
	handle C.plCollisionShapeHandle
}

func (this CollisionShape) AddChildShape(s core.CollisionShape, p mgl64.Vec3, o mgl64.Quat) {
	vec := vec3_to_bullet(p)
	quat := quat_to_bullet(o)
	C.plAddChildShape(this.handle, s.(CollisionShape).handle, &vec[0], &quat[0])
}

func (this CollisionShape) AddVertex(v mgl64.Vec3) {
	C.plAddVertex(this.handle, C.plReal(v.X()), C.plReal(v.Y()), C.plReal(v.Z()))
}
