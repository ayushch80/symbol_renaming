package renamer

import (
	"debug/macho"
)

type SymbolCategory struct {
	Symbols []macho.Symbol
	Category   string
}
