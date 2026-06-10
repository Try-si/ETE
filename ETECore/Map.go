package ETECore

import (
	"strings"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
	tiled "github.com/lafriks/go-tiled"
)

func (g *Game) InitMap() {
	g.Maps = make(map[string]*Map)
	for _, mapName := range g.MapConfig.Maps {
		g.Maps[mapName] = g.MapConfig.LoadMap(mapName)
	}
}

func (mc *MapConfig) LoadMap(mapName string) *Map {
	dir := strings.Join(strings.Split(mc.G.GetGame().Config.MapsPath, "/")[:len(strings.Split(mc.G.GetGame().Config.MapsPath, "/"))-1], "/")

	tiledMap, err := tiled.LoadFile(dir + "/" + mc.TiledMap + "/" + mapName + ".tmx")
	if err != nil {
		panic(err)
	}

	jsonMaps := ETEHelper.JsonToStruct[JsonMap](dir + "/" + mc.JsonMap + "/" + mapName + ".json")

	tilesets := []*ebiten.Image{}

	for _, tileset := range tiledMap.Tilesets {
		for _, tile := range ETEHelper.SliceImageByGrid(ETEHelper.LoadImage(tileset.Image.Source), (tileset.TileWidth+tileset.TileHeight)/2) {
			tilesets = append(tilesets, ebiten.NewImageFromImage(tile))
		}
	}

	resultat := Map{
		Map: MapData{
			Tileset: tilesets,
		},
		CellSize: jsonMaps.CellSize,
		Unité:    jsonMaps.Unite,
		Cam:      jsonMaps.Cam,
		Elements: jsonMaps.Elements,
	}

	return &resultat
}

func (g *Game) InitTile() {
	for i, j := range ETEHelper.JsonToStruct[map[string]Tile](g.MapConfig.Tiles) {
		g.Tiles[i] = &j
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
			if isAnimated {
				b := tile.Game.GetGame().Animations[tile.Game.GetGame().Tiles[string(tile.Id)].Animation].Frames[tile.Frame].Box
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].CellSize
				lay[la][[9]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0}] = img
			} else {
				b := tile.Game.GetGame().Tiles[string(tile.Id)].Box
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].CellSize
				lay[la][[9]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0}] = img
			}
		}
	}

	return lay
}
