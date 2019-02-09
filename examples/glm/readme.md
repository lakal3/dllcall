# GLM example

Use GLM library to multiply vectors with 4x4 matrix

# Building C++ library

Download glm library from [glm releases](https://github.com/g-truc/glm/tags)

Run CMAKE from glmcpp directory and set GLM_INCLUDE_DIR to directory where you installed glm headers

Build C++ project and copy or link generated shared library to a location where operating system can locate it.

# Building go executable

Just use go build in glmgo directory  
  