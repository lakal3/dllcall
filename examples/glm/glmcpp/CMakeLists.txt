cmake_minimum_required (VERSION 3.0)
project (glmcpp)

SET(GLM_INCLUDE_DIR "" CACHE PATH "GLM include directory")
if (NOT GLM_INCLUDE_DIR)
  message(SEND_ERROR "Set GLM_INCLUDE_DIR to GLM root")
endif()

message("GLM path ${GLM_INCLUDE_DIR}")

if(NOT CMAKE_DEBUG_POSTFIX)
  set(CMAKE_DEBUG_POSTFIX d)
endif()

include_directories(${GLM_INCLUDE_DIR})

add_library(glmcpp SHARED "ifimpl.cpp")

install(TARGETS glmcpp DESTINATION bin)
