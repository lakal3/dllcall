cmake_minimum_required (VERSION 3.0)
project (fibonlib)

if(NOT CMAKE_DEBUG_POSTFIX)
  set(CMAKE_DEBUG_POSTFIX d)
endif()

add_library(fibonlib SHARED "fibon.cpp" "fibon_if.h")

install(TARGETS fibonlib DESTINATION bin)
