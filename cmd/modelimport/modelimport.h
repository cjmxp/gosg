#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <assimp/cimport.h>
#include <assimp/scene.h>
#include <assimp/postprocess.h>

#ifdef __cplusplus
extern "C" {
#endif

size_t get_mesh_count(struct aiScene *scene);
size_t get_vertex_count(struct aiScene *scene, int mesh_idx);
void get_mesh_maps(struct aiScene *scene, int mesh_idx, char *diffuse, char *normal);
void get_mesh_name(struct aiScene *scene, int mesh_idx, char *name);
int get_mesh_wrapmode(struct aiScene *scene, int mesh_idx);
float *get_positions(struct aiScene *scene, int mesh_idx);
float *get_normals(struct aiScene *scene, int mesh_idx);
float *get_tangents(struct aiScene *scene, int mesh_idx);
float *get_bitangents(struct aiScene *scene, int mesh_idx);
float *get_texturecoords(struct aiScene *scene, int mesh_idx);
unsigned int get_indexcount(struct aiScene *scene, int mesh_idx);
unsigned int *get_indices(struct aiScene *scene, int mesh_idx);

#ifdef __cplusplus
}
#endif
