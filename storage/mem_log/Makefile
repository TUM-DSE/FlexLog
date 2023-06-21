libstorage.so: Log.o cLog.o
	g++ Log.o cLog.o -o libstorage.so -shared -L. -ltbb

Log.o: Log.cpp Log.hpp
	g++ Log.cpp -c -fPIC -shared -o Log.o
	
cLog.o: cLog.cpp Log.h
	g++ cLog.cpp -c -fPIC -shared -o cLog.o
	
clean: 
	rm libstorage.so Log.o cLog.o
