#include <iostream>
#include <fstream>

#include "model.pb.h"

#include <assimp/cimport.h>
#include <assimp/scene.h>
#include <assimp/postprocess.h>

std::string load_texture(std::string fullpath);

int main(int argc, char **argv) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " modelfile" << std::endl;
        return 1;
    }

    // copy modelfile
    std::string modelfile = std::string(argv[1]);
    auto base_directory = modelfile.substr(0, modelfile.find_last_of("/\\"));
    std::cerr << base_directory <<std::endl;

    // set our global options. this will change
    auto opts = 0 |
             aiProcessPreset_TargetRealtime_MaxQuality |
            aiProcess_CalcTangentSpace |
            aiProcess_FixInfacingNormals |
            aiProcess_FlipUVs |
            aiProcess_GenSmoothNormals |
            aiProcess_ImproveCacheLocality |
            aiProcess_TransformUVCoords |
            aiProcess_Triangulate |
            aiProcess_ValidateDataStructure |
            aiProcess_OptimizeGraph |
            aiProcess_OptimizeMeshes;

    auto scene = aiImportFile(modelfile.c_str(), opts);
    if (scene == nullptr) {
        std::cerr << aiGetErrorString() << std::endl;
        return 1;
    }

    // create model
    auto model = new protos::Model();

    // iterate through meshes
    for (int i=0; i<scene->mNumMeshes; i++) {
        auto vertexCount = scene->mMeshes[i]->mNumVertices;

        // create a mesh
        protos::Mesh *mesh = model->add_meshes();

        // set name
        mesh->set_name(scene->mMeshes[i]->mName.data);
        std::cerr << "Name: " << mesh->name() << std::endl;

        // set buffers
        mesh->set_positions(scene->mMeshes[i]->mVertices, sizeof(float)*vertexCount*3);
        mesh->set_normals(scene->mMeshes[i]->mNormals, sizeof(float)*vertexCount*3);
        mesh->set_tangents(scene->mMeshes[i]->mTangents, sizeof(float)*vertexCount*3);
        mesh->set_bitangents(scene->mMeshes[i]->mBitangents, sizeof(float)*vertexCount*3);
        mesh->set_tcoords(scene->mMeshes[i]->mTextureCoords[0], sizeof(float)*vertexCount*3);

        // indices are trickier, we're forcing uint16_t, may need to bump this for complex meshes
        auto indexBuffer = new uint16_t[scene->mMeshes[i]->mNumFaces * 3];
        for (size_t f=0; f<scene->mMeshes[i]->mNumFaces; f++) {
            indexBuffer[f*3+0] = scene->mMeshes[i]->mFaces[f].mIndices[0];
            indexBuffer[f*3+1] = scene->mMeshes[i]->mFaces[f].mIndices[1];
            indexBuffer[f*3+2] = scene->mMeshes[i]->mFaces[f].mIndices[2];
        }
        mesh->set_indices(indexBuffer, scene->mMeshes[i]->mNumFaces*3*sizeof(uint16_t));

        // material and textures
        aiMaterial *mat = scene->mMaterials[scene->mMeshes[i]->mMaterialIndex];

        struct aiString materialName;
        aiGetMaterialString(mat, AI_MATKEY_NAME, &materialName);
        mesh->set_material(materialName.data);

        struct aiString albedoPath, normalPath;
        aiGetMaterialString(mat, AI_MATKEY_TEXTURE_DIFFUSE(0), &albedoPath);
        aiGetMaterialString(mat, AI_MATKEY_TEXTURE_HEIGHT(0), &normalPath);

        if (albedoPath.length > 0) {
            auto fullAlbedoPath = base_directory + "/" + std::string(albedoPath.data);
            mesh->set_albedo_map(load_texture(fullAlbedoPath));
        }

        if (normalPath.length > 0) {
            auto fullNormalPath = base_directory + "/" + std::string(normalPath.data);
            mesh->set_normal_map(load_texture(fullNormalPath));
        }
    }

    // save
    {
        std::fstream output("out.model", std::ios::out | std::ios::trunc | std::ios::binary);
        if (!model->SerializeToOstream(&output)) {
            std::cerr << "Failed to write model." << std::endl;
            return -1;
        }
    }

    google::protobuf::ShutdownProtobufLibrary();
    return 0;
}

std::string load_texture(std::string fullpath) {
    std::ifstream textureFile(fullpath, std::ifstream::binary | std::ifstream::ate);

    std::streamsize textureSize = textureFile.tellg();
    textureFile.seekg(0, std::ios::beg);

    char *buffer = new char[textureSize];
    textureFile.read(buffer, textureSize);
    textureFile.close();

    return std::string(buffer, textureSize);
}