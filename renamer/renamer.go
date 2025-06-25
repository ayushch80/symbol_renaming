package renamer

import (
	"fmt"
	"debug/macho"
	"path/filepath"
	"log"
	"os"
)

func SymbolRenaming(tmpFolder string) {
	fmt.Println("[+] Reading Mach-O binary")

	binPath := filepath.Join(tmpFolder, "Payload", "MyBankApp.app", "MyBankApp")

	machoBinary, err := macho.Open(binPath)
	if err != nil {
		log.Fatal("[-] Failed while reading the binary :", err.Error())
	}

	_, err = os.ReadFile(binPath)
	// machoData, err := os.ReadFile(binPath)
	if err != nil {
		log.Fatal("[-] Failed while reading the binary :", err.Error())
	}

	syms := machoBinary.Symtab.Syms
	
	fmt.Println(machoBinary.Symtab.Syms)
}