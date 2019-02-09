cmake_minimum_required (VERSION 3.0)
project (docdb)

if(NOT CMAKE_DEBUG_POSTFIX)
  set(CMAKE_DEBUG_POSTFIX d)
endif()

add_library(docdb SHARED "docdb.cpp")

