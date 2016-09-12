#include "modelimport.h"
#include <iostream>

#ifdef __cplusplus
extern "C" {
#endif

size_t get_mesh_count(struct aiScene *scene) {
    return scene->mNumMeshes;
}

size_t get_vertex_count(struct aiScene *scene, int mesh_idx) {
    return scene->mMeshes[0]->mNumVertices;
}

void get_mesh_maps(struct aiScene *scene, int mesh_idx, char *diffuse, char *normal) {
    if (scene->mMeshes[mesh_idx]->mMaterialIndex > scene->mNumMaterials) {
        return;
    }

    struct aiMaterial *mat = scene->mMaterials[scene->mMeshes[mesh_idx]->mMaterialIndex];
    struct aiString texturePath;

    if (AI_SUCCESS == aiGetMaterialString(mat, AI_MATKEY_TEXTURE_DIFFUSE(0), &texturePath)) {
        strcpy(diffuse, texturePath.data);
    } else {
        strcpy(diffuse, "");
    }

    if (AI_SUCCESS == aiGetMaterialString(mat, AI_MATKEY_TEXTURE_HEIGHT(0), &texturePath)) {
        strcpy(normal, texturePath.data);
    } else {
        strcpy(normal, "");
    }
}

void get_mesh_name(struct aiScene *scene, int mesh_idx, char *name) {
    std::cerr << scene->mMeshes[mesh_idx]->mName.C_Str() << std::endl;
    strcpy(name, scene->mMeshes[mesh_idx]->mName.data);
}

int get_mesh_wrapmode(struct aiScene *scene, int mesh_idx) {
    if (scene->mMeshes[mesh_idx]->mMaterialIndex > scene->mNumMaterials) {
        return -1;
    }

    struct aiMaterial *mat = scene->mMaterials[scene->mMeshes[mesh_idx]->mMaterialIndex];
    int mode;
    unsigned int max = 1;
    if (AI_SUCCESS == aiGetMaterialIntegerArray(mat, AI_MATKEY_MAPPINGMODE_U_DIFFUSE(0), &mode, &max)) {
        return mode;
    }

    return -1;
}

float get_mesh_opacity(struct aiScene *scene, int mesh_idx) {
    if (scene->mMeshes[mesh_idx]->mMaterialIndex > scene->mNumMaterials) {
        return 1.0;
    }

    float opacity;
    struct aiMaterial *mat = scene->mMaterials[scene->mMeshes[mesh_idx]->mMaterialIndex];
    if (AI_SUCCESS == aiGetMaterialFloat(mat, AI_MATKEY_OPACITY, (float *)&opacity)) {
        return opacity;
    }

    return 1.0;
}

float *get_positions(struct aiScene *scene, int mesh_idx) {
    return (float *)scene->mMeshes[mesh_idx]->mVertices;
}

float *get_normals(struct aiScene *scene, int mesh_idx) {
    return (float *)scene->mMeshes[mesh_idx]->mNormals;
}

float *get_tangents(struct aiScene *scene, int mesh_idx) {
    return (float *)scene->mMeshes[mesh_idx]->mTangents;
}

float *get_bitangents(struct aiScene *scene, int mesh_idx) {
    return (float *)scene->mMeshes[mesh_idx]->mBitangents;
}

float *get_texturecoords(struct aiScene *scene, int mesh_idx) {
    return (float *)scene->mMeshes[mesh_idx]->mTextureCoords[0];
}

unsigned int get_indexcount(struct aiScene *scene, int mesh_idx) {
    return scene->mMeshes[mesh_idx]->mNumFaces * 3;
}

unsigned int *get_indices(struct aiScene *scene, int mesh_idx) {
    struct aiMesh *mesh = scene->mMeshes[mesh_idx];

    size_t indexcount = mesh->mNumFaces * 3;

    unsigned int *buf = (unsigned int *)malloc(sizeof(unsigned int) * indexcount);
    for (size_t f=0; f<mesh->mNumFaces; f++) {
        memcpy(&buf[f*3], mesh->mFaces[f].mIndices, 3*sizeof(unsigned int));
    }
    return buf;
}

#ifdef __cplusplus
}
#endif
