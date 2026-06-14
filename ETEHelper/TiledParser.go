package ETEHelper

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type TMXMap struct {
	XMLName    xml.Name   `xml:"map"`
	Version    string     `xml:"version,attr"`
	TiledVersion string   `xml:"tiledversion,attr"`
	Orientation string    `xml:"orientation,attr"`
	Width      int        `xml:"width,attr"`
	Height     int        `xml:"height,attr"`
	TileWidth  int        `xml:"tilewidth,attr"`
	TileHeight int        `xml:"tileheight,attr"`
	Infinite   int        `xml:"infinite,attr"`
	Tilesets   []TMXTileset `xml:"tileset"`
	Layers     []TMXLayer `xml:"layer"`
}

type TMXTileset struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type TMXLayer struct {
	ID         int           `xml:"id,attr"`
	Name       string        `xml:"name,attr"`
	Width      int           `xml:"width,attr"`
	Height     int           `xml:"height,attr"`
	Properties TMXProperties `xml:"properties"`
	Data       TMXData       `xml:"data"`
}

type TMXProperties struct {
	Property []TMXProperty `xml:"property"`
}

type TMXProperty struct {
	Name  string  `xml:"name,attr"`
	Value string  `xml:"value,attr"`
	Type  string  `xml:"type,attr"`
}

type TMXData struct {
	Encoding string      `xml:"encoding,attr"`
	Chunks   []TMXChunk `xml:"chunk"`
}

type TMXChunk struct {
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Data   string `xml:",chardata"`
}

type TSXTileset struct {
	XMLName   xml.Name `xml:"tileset"`
	Name      string   `xml:"name,attr"`
	TileWidth int      `xml:"tilewidth,attr"`
	TileHeight int     `xml:"tileheight,attr"`
	TileCount int      `xml:"tilecount,attr"`
	Image     TSXImage `xml:"image"`
}

type TSXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type TileInstance struct {
	LayerID int
	X       int
	Y       int
	GID     uint32
}

func LoadTMX(path string) (*TMXMap, error) {
	data := LoadJson(path)
	var tmx TMXMap
	err := xml.Unmarshal([]byte(data), &tmx)
	if err != nil {
		return nil, fmt.Errorf("error parsing TMX: %v", err)
	}
	return &tmx, nil
}

func LoadTSX(path string) (*TSXTileset, error) {
	data := LoadJson(path)
	var tsx TSXTileset
	err := xml.Unmarshal([]byte(data), &tsx)
	if err != nil {
		return nil, fmt.Errorf("error parsing TSX: %v", err)
	}
	return &tsx, nil
}

func ConvertTMXToInternal(tmx *TMXMap) []TileInstance {
	var tiles []TileInstance

	for _, layer := range tmx.Layers {
		for _, chunk := range layer.Data.Chunks {
			lines := strings.Split(strings.TrimSpace(chunk.Data), "\n")
			for y, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				values := strings.Split(line, ",")
				for x, val := range values {
					val = strings.TrimSpace(val)
					if val == "" {
						continue
					}
					gid, err := strconv.ParseUint(val, 10, 32)
					if err != nil {
						continue
					}
					if gid == 0 {
						continue
					}
					// Calculer la position absolue
					absX := chunk.X + x
					absY := chunk.Y + y
					tiles = append(tiles, TileInstance{
						LayerID: layer.ID,
						X:       absX,
						Y:       absY,
						GID:     uint32(gid),
					})
				}
			}
		}
	}

	return tiles
}
