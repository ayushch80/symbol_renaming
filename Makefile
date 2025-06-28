CXX       = g++
# Include project root, local headers, and manual LIEF include
CXXFLAGS  = -std=c++11 -fPIC -I. -Iinclude -Ilief/include
# Link against local LIEF library and set runtime path
LDFLAGS   = -shared -Llief/lib -Wl,-rpath,$(PWD)/lief/lib
LIBS      = -llief -ldl

all: librenamer.so main

# Build the shared library, linking against LIEF
librenamer.so: renamer/renamer.o
	$(CXX) $(LDFLAGS) -o $@ renamer/renamer.o $(LIBS)

# Compile C++ renamer object
renamer/renamer.o: renamer/renamer.cpp renamer/renamer.h
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Build Go main; ensure cgo sees the right include paths
main: librenamer.so main.go
	CGO_CXXFLAGS="-std=c++11 -I. -Iinclude -Ilief/include" \
	CGO_LDFLAGS="-Llief/lib -llief -ldl -Wl,-rpath,$(PWD)/lief/lib" \
	go build -o main main.go

clean:
	rm -f renamer/renamer.o librenamer.so main