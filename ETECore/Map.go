package ETECore

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

var dir string

func (g *Game) InitMap() {
	g.Maps = make(map[string]*Map)
	for _, mapName := range g.MapConfig.Maps {
		g.Maps[mapName] = g.MapConfig.LoadMap(mapName)
	}
}

func (mc *MapConfig) LoadMap(mapName string) *Map {
	dir = strings.Join(strings.Split(mc.G.GetGame().Config.MapsPath, "/")[:len(strings.Split(mc.G.GetGame().Config.MapsPath, "/"))-1], "/")

	tmxPath := dir + "/" + mc.TiledMap + "/" + mapName + ".tmx"

	// 1. Charger le JSON
	jsonMaps := ETEHelper.JsonToStruct[JsonMap](dir + "/" + mc.JsonMap + "/" + mapName + ".json")

	// 2. Charger le TMX avec VOTRE code
	tiledMap, err := ETEHelper.LoadTMX(tmxPath, jsonMaps.PropertiesForHeight)
	if err != nil {
		panic(err) // Ou gestion d'erreur propre
	}

	// 3. Charger les Tilesets référencés
	for _, tileset := range tiledMap.Tilesets {
		for i, tile := range ETEHelper.SliceImageByGrid(ETEHelper.LoadImage(dir+"/"+mc.TiledMap+tileset.Source), int(jsonMaps.Unite)) {
			mc.G.GetGame().Sprites[strconv.Itoa(int(tileset.FirstGID)+i)] = ebiten.NewImageFromImage(tile)
		}
	}

	// 4. Créer la map
	resultat := Map{
		Map: MapData{
			Tiles: TileToTile(tiledMap.GetTiles(), mc.G), // <--- On injecte les tuiles converties ici
		},
		CellSize: jsonMaps.CellSize,
		Unité:    jsonMaps.Unite,
		Cam:      jsonMaps.Cam,
		Elements: jsonMaps.Elements, // Gardez le JSON pour les objets dynamiques (ennemis, joueur, etc.)
	}

	return &resultat
}

func (g *Game) InitTile() {
	g.Tiles = make(map[string]*Tile)
	for i, j := range ETEHelper.JsonToStruct[map[string]*Tile](dir + "/" + g.MapConfig.Tiles) {
		tileId, _ := strconv.Atoi(i)
		g.Tiles[strconv.Itoa(tileId)] = j
		println("Loaded tile: " + strconv.Itoa(tileId))
	}
}

func (m *Map) GetSpriteByOrderYZX() map[int]map[[9]float32]*ebiten.Image { // [witdh/radius, height, xOffset, yOffset, xPos, yPos, xSize, ySize, rotation]
	resultat := make(map[int]map[[9]float32]*ebiten.Image)

	for _, e := range m.Elements {
		resultat[e.Layer] = map[[9]float32]*ebiten.Image{
			{
				e.Box[0], e.Box[1],
				e.Box[2], e.Box[3],
				float32(e.Pos[0]), float32(e.Pos[1]),
				float32(e.Size[0]), float32(e.Size[1]),
				e.Rotation,
			}: e.GetSprite(),
		}
	}

	for k, v := range m.GetTileByLayer() {
		resultat[k] = make(map[[9]float32]*ebiten.Image)
		for pos, tile := range v {
			resultat[k][[9]float32{
				float32(pos[0]), float32(pos[1]),
				float32(pos[2]), float32(pos[3]),
				float32(pos[4]), float32(pos[5]),
				float32(pos[6]), float32(pos[7]),
				float32(pos[8]),
			}] = tile
		}
	}

	return resultat
}

func (m *Map) GetElementByLayer() map[int][]Element {
	lay := make(map[int][]Element)
	for _, v := range m.Elements {
		v.PushFrame()
		lay[v.Layer] = append(lay[v.Layer], *v)
	}

	return lay
}

func (m *Map) GetTileByLayer() map[int]map[[9]int]*ebiten.Image {
	lay := make(map[int]map[[9]int]*ebiten.Image)
	for la, tiles := range m.Map.Tiles {
		lay[la] = make(map[[9]int]*ebiten.Image)
		for pos, tile := range tiles {
			tile.PushFrame()
			img, isAnimated := tile.GetSprite()
			if img == nil {
				img = ebiten.NewImage(32, 32)
				img.Fill(color.RGBA{255, 6, 181, 255})
			}
			if isAnimated {
				b := tile.Game.GetGame().Animations[tile.Game.GetGame().Tiles[strconv.Itoa(tile.Id)].Animation].Frames[tile.Frame-1].Box
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].CellSize
				lay[la][[9]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0}] = img
			} else {
				var b [4]int
				if _, exist := tile.Game.GetGame().Tiles[strconv.Itoa(int(tile.Id))]; !exist {
					fmt.Printf("Tile box: %d not exist\n", tile.Id)
					b = [4]int{0, 0, 0, 0}
				} else {
					b = tile.Game.GetGame().Tiles[strconv.Itoa(int(tile.Id))].Box
				}
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].CellSize
				lay[la][[9]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0}] = img
			}
		}
	}

	return lay
}

func TileToTile(tiles map[int]map[[2]int]*tiled.LayerTile, G IForGame) map[int]map[[2]int]*TileElement {
	result := make(map[int]map[[2]int]*TileElement)
	for layer, layerTiles := range tiles {
		result[layer] = make(map[[2]int]*TileElement)
		for pos, tile := range layerTiles {
			f := 0
			if G.GetGame().Tiles[strconv.Itoa(int(tile.ID))].Animation != "nil" {
				f = 1
			}
			result[layer][[2]int{pos[0], pos[1]}] = &TileElement{
				Id:    int(tile.ID),
				Frame: f,
				Game:  G,
			}
		}
	}
	return result
}
