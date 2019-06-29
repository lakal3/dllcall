cmake_minimum_required (VERSION 3.0)
project (hello)

if(NOT CMAKE_DEBUG_POSTFIX)
  set(CMAKE_DEBUG_POSTFIX d)
endif()

add_library(greeting SHARED "greeting.cpp" "greeting.h")

install(TARGETS greeting DESTINATION bin)
