cmake_minimum_required (VERSION 3.0)
project (msgbox)

if(NOT CMAKE_DEBUG_POSTFIX)
  set(CMAKE_DEBUG_POSTFIX d)
endif()

add_compile_definitions(_CRT_SECURE_NO_WARNINGS)

include_directories(".")

add_library(msgboxlib SHARED "msgbox.cpp" "msgboxif.h")

install(TARGETS msgboxlib DESTINATION bin)
