package ETEHelper

import "github.com/lafriks/go-tiled"

type MapTMX struct {
	Layers   map[int]Layer
	Tilesets []Tileset
}

type Layer struct {
	ID                  uint32
	PropertiesForHeight float64
	Tiles               map[[2]int]*tiled.LayerTile
	Width, Height       int
}

type Tileset struct {
	FirstGID uint32
	Source   string
}
