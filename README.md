# gosg [![GoDoc](https://godoc.org/github.com/fcvarela/gosg?status.svg)](https://godoc.org/github.com/fcvarela/gosg)

Package gosg is a lightweight screnegraph/scenetree based rendering toolkit. It provides a set of tools to build standalone windowed 3D applications for MacOS, Linux and Windows.

Applications are built by providing a 'ClientApplication' object which knows how to build a set of scenes and handle scene transitions.

Scenes are built by specifying a graph of nodes, cameras, lights and some configuration for types of perspective to use, how to clear the screen, etc.

## Features
This package provides minimal (but improving over time) wrappers for an immediate-mode GUI system, a physics engine and rendering system. There is currently only one rendering backend and it is based on OpenGL 3.3 (Core). OpenGL is currently stuck at 4.1 on OSX and 3.3 for iGPUs on Linux. This will be kept as up-to-date as possible as support for more recent GL releases is added to either of these OSs. A Vulkan rendering backend is in the works too, but don't expect it to work on OSX as there is no official support from Apple or any of the GPU vendors.

## Warning
This package is experimental and expect it to break at any time. It will also change dramatically when it moves to zero-dependency subpackages and separate scenegraph + submission packages.

## Contributing
Contributions in the form of features, improvements or even architectural changes are most welcome and should be submitted in the form of pull-requests following a discussion ticket/issue.

