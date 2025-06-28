package renamer

/*
#cgo CXXFLAGS: -std=c++11 -I${SRCDIR} -I${SRCDIR}/../include -I${SRCDIR}/../lief/include
#cgo LDFLAGS: -L${SRCDIR}/.. -L${SRCDIR}/../lief/lib -lrenamer -llief -ldl -Wl,-rpath,${SRCDIR}/../lief/lib
#include <stdlib.h>
#include "renamer.h"
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"unsafe"
)

func SymbolRenaming(tmpFolder string) {
	fmt.Println("[+] Reading Mach-O binary")

	binPath := filepath.Join(tmpFolder, "Payload", "MyBankApp.app", "MyBankApp")
	
	// Convert Go string to C string
	cBinPath := C.CString(binPath)
	defer C.free(unsafe.Pointer(cBinPath))
	
	// Call the C function
	C.SymbolRenaming(cBinPath)
}