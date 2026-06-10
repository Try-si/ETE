package ETECore

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) InitMap() {
	g.Maps = make(map[string]*Map)
	for _, mapName := range g.MapConfig.Maps {
		g.Maps[mapName] = g.MapConfig.LoadMap(mapName)
	}
}

func (mc *MapConfig) LoadMap(mapName string) *Map {
	dir := strings.Join(strings.Split(mc.G.GetGame().Config.MapsPath, "/")[:len(strings.Split(mc.G.GetGame().Config.MapsPath, "/"))-1], "/")

	tmxPath := dir + "/" + mc.TiledMap + "/" + mapName + ".tmx"

	// 1. Charger le TMX avec VOTRE code
	tiledMap, err := ETEHelper.LoadTMX(tmxPath)
	if err != nil {
		panic(err) // Ou gestion d'erreur propre
	}

	// 2. Charger les Tilesets référencés
	tilesetsImages := []*ebiten.Image{}

	// On itère sur les références de tilesets dans le TMX
	for _, tsRef := range tiledMap.Tilesets {
		// Construire le chemin du .tsx
		// Attention : le chemin dans 'Source' est relatif au .tmx
		tsxAbsPath := filepath.Join(filepath.Dir(tmxPath), tsRef.Source)

		tsxData, err := ETEHelper.LoadTSX(tsxAbsPath)
		if err != nil {
			panic(fmt.Sprintf("Erreur tileset %s: %v", tsRef.Source, err))
		}

		// Charger l'image et la découper (votre logique existante)
		imgPath := filepath.Join(filepath.Dir(tsxAbsPath), tsxData.Image.Source)
		fullImg := ETEHelper.LoadImage(imgPath)

		// Découpage
		for _, tile := range ETEHelper.SliceImageByGrid(fullImg, (tsxData.TileWidth+tsxData.TileHeight)/2) {
			tilesetsImages = append(tilesetsImages, ebiten.NewImageFromImage(tile))
		}

		// NOTE: Ici, vous devez faire attention au FirstGID.
		// Si vous avez plusieurs tilesets, l'ID 5 du tileset A n'est pas l'index 5 de votre slice.
		// Il faut une logique de mapping : GID -> (TilesetIndex, LocalIndex)
		// Exemple : if gid >= tsRef.FirstGID && gid < tsRef.FirstGID + tsxData.TileCount ...
	}

	jsonMaps := ETEHelper.JsonToStruct[JsonMap](dir + "/" + mc.JsonMap + "/" + mapName + ".json")

	resultat := Map{
		Map: MapData{
			Tileset: tilesetsImages,
			// Vous devrez probablement stocker les données brutes des calques ici
			// pour pouvoir les dessiner en tenant compte des chunks négatifs.
			// Layers: tiledMap.Layers,
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
