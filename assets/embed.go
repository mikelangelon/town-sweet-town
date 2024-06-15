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

	//go:embed roguelikeChar_transparent.png
	Characters []byte

	//go:embed roguelikeChar_transparent.tsx
	CharactersTSX []byte
)
