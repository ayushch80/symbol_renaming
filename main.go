package main

import (
	"fmt"

	"symbol_renaming/renamer"
	"symbol_renaming/utils"
)

func main() {
	fmt.Println("[*] SYMBOL RENAMING")

	inputIPAPath := "/Users/brluser/Desktop/MyBankApp.ipa"

	fmt.Println("[+] Unzipping the IPA")
	tmpFolder := utils.Unzip(inputIPAPath)

	renamer.SymbolRenaming(tmpFolder)

	deleteFolder := true
	if deleteFolder {
		// utils.DeleteFolder(tmpFolder)
	}
}
