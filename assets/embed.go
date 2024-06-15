package assets

import (
	_ "embed"
)

var (
	//go:embed tilemap_packed.png
	TileMapPacked []byte

	//go:embed tilemap_packed.tsx
	TileMapPackedTSX []byte

	//go:embed town1.tmx
	Town1 []byte

	//go:embed town2.tmx
	Town2 []byte

	//go:embed town3.tmx
	Town3 []byte

	//go:embed roguelikeChar_transparent.png
	Characters []byte

	//go:embed roguelikeChar_transparent.tsx
	CharactersTSX []byte
)
