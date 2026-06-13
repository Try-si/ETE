package ETEHelper

import (
	tiled "github.com/lafriks/go-tiled"
)

func LoadTMX(path string, propertiesForHeight string) (*MapTMX, error) {
	tiledMap, err := tiled.LoadFile(path)
	if err != nil {
		return nil, err
	}
	return convertTMXToInternal(tiledMap, propertiesForHeight), nil
}

func convertTMXToInternal(tmx *tiled.Map, propertiesForHeight string) *MapTMX {
	layers := make(map[int]Layer, len(tmx.Layers))
	tilesets := make([]Tileset, len(tmx.Tilesets))

	for _, layer := range tmx.Layers {
		layers[int(layer.ID)] = Layer{
			PropertiesForHeight: layer.Properties.GetFloat(propertiesForHeight),
			Tiles:               ListToGridYWX(layer.Tiles, tmx.Width, tmx.Height, [2]int{tmx.Width / 2, tmx.Height / 2}),
		}
	}

	for i, tileset := range tmx.Tilesets {
		tilesets[i] = Tileset{
			FirstGID: tileset.FirstGID,
			Source:   tileset.Source,
		}
	}

	return &MapTMX{
		Layers:   layers,
		Tilesets: tilesets,
	}
}

func (m *MapTMX) GetTiles() map[int]map[[2]int]*tiled.LayerTile {
	tiles := make(map[int]map[[2]int]*tiled.LayerTile)
	for l, layer := range m.Layers {
		tiles[l] = layer.Tiles
	}
	return tiles
}
