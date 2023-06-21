CCX		= g++
CXX_FLAGS	= -Wall -g -O3
WORKING_DIR	= /home/${USER}/PMLog-bench
INCLUDE_PATH	= -I${WORKING_DIR}/boki/deps/rocksdb/include
ROCKSDB_PATH	= ${WORKING_DIR}/boki/deps/rocksdb
LD_FLAGS	= -L${WORKING_DIR}/boki/deps/rocksdb
LIBS		= -lfmt -pthread -lgflags -ldl -lz -lrocksdb



SOURCES		= rocksdb_bench.cpp
OBJECTS		= $(SOURCES:.cpp=.o)
NAMES		= $(SOURCES:.cpp=)

.PHONY = all clean 
all: clean rocksdb_bench

rocksdb_bench: rocksdb
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) $@.cpp -o $@ ${LD_FLAGS} ${LIBS}

rocksdb:
	make -C $(ROCKSDB_PATH) static_lib -j8

clean: 
	# make clean -C $(ROCKSDB_PATH)
	-rm -f rocksdb_bench
