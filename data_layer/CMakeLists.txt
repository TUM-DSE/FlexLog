cmake_minimum_required(VERSION 3.14)

add_library(data_layer replica.cpp order_layer_proxy.cpp pm_storage.cpp memory_storage.cpp)

target_include_directories(data_layer PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

add_executable(run_replica main.cpp)

target_link_libraries(run_replica PUBLIC data_layer)
