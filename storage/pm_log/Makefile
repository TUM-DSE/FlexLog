CCX		= g++
CXX_FLAGS	= -Wall -g -O3
WORKING_DIR	= /home/${USER}/PMLog-bench
INCLUDE_PATH	= -I${WORKING_DIR}/pmdk/local/usr/include -I${WORKING_DIR}/libpmemobj-cpp/include
LD_FLAGS	= -L${WORKING_DIR}/pmdk/local/lib -L.
LIBS		= -lfmt -lstorage -lpmemobj -ltbb -lpmem -pthread -lgflags



SOURCES		= PMLog.cpp cPMLog.cpp LogCache.cpp
#HEADERS_H 	= $(shell echo *.h)
#HEADERS_HPP 	= $(shell echo *.hpp)

OBJECTS		= $(SOURCES:.cpp=.o)
NAMES		= $(SOURCES:.cpp=)

.PHONY = all clean 
all: clean $(NAMES) pmem_bench main

libstorage.so : $(NAMES)
	 $(CXX) ${CXX_FLAGS} ${INCLUDE_PATH} -shared -fPIC -o $@ ${OBJECTS} $(LD_FLAGS) ${LIBS}

$(NAMES): $(SOURCES)
	$(CXX) $(CXX_FLAGS) ${INCLUDE_PATH} -c $@.cpp -o $@.o $(LD_FLAGS) $(LIBS)

pmem_bench: libstorage.so
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) $@.cpp -o $@ ${LD_FLAGS} ${LIBS}

main: libstorage.so
	$(CXX) $(CXX_FLAGS) $(INCLUDE_PATH) $@.cpp -o $@ ${LD_FLAGS} ${LIBS}

clean: 
	-rm -f libstorage.so PMLog.o cPMLog.o LogCache.o main.o main pmem_bench pmem_bench.o
	-rm -f /mnt/pmem0/dimitra/PMLog3*
