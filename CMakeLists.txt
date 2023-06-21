
cmake_minimum_required(VERSION 3.14)
project(cppcourse)


# try to prevent in-source builds
if("${PROJECT_SOURCE_DIR}" STREQUAL "${PROJECT_BINARY_DIR}")
    message(FATAL_ERROR "In-source builds are not allowed. Please create a separate 'build' directory and build in there.")
endif()

# set default c++ standard for compile targets
set(CMAKE_CXX_STANDARD 20)

find_package(doctest REQUIRED)

if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE "Debug")
endif()

#------- tweaks for an improved build experience -------

# so language servers like clangd know what the buildsystem does
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# windows needs "exported" symbols for linking
set(CMAKE_WINDOWS_EXPORT_ALL_SYMBOLS ON)

# display many warnings during compilation
# they usually help you find some bugs
if(CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang|AppleClang")
    add_compile_options(-Wall -Wextra -Wconversion -pedantic -Wfatal-errors)
elseif(CMAKE_CXX_COMPILER_ID MATCHES "MSVC")
    add_compile_options(/W3 /WX)

    # so the dll files are alongside the executables
    # and the windows loader is happy (no more error 0xc0000135)
    set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
endif()



#------- homework code inclusion -------

add_subdirectory("data_layer")
