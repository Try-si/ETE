package ETECore

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
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

	// 2. Charger le TMX avec le parser personnalisé
	tiledMap, err := ETEHelper.LoadTMX(tmxPath)
	if err != nil {
		panic(err) // Ou gestion d'erreur propre
	}

	// 3. Charger les Tilesets référencés
	tilesetsImages := []*ebiten.Image{}

	for _, tsRef := range tiledMap.Tilesets {
		// Construire le chemin correct pour le TSX
		// Le fichier TMX est dans Maps/Maps/, le source est ../Tileset/Sol.tsx
		// Donc le chemin absolu est Maps/Tileset/Sol.tsx
		tsxAbsPath := dir + "/Tileset/" + tsRef.Source[len("../Tileset/"):]
		tsxData, err := ETEHelper.LoadTSX(tsxAbsPath)
		if err != nil {
			panic(fmt.Sprintf("Erreur tileset %s: %v", tsRef.Source, err))
		}

		// Construire le chemin correct pour l'image
		// L'image est ../../Textures/tileset_1.png relatif au TSX
		// Donc le chemin absolu est Textures/tileset_1.png
		imgPath := dir + "/../Textures/" + tsxData.Image.Source[len("../../Textures/"):]
		fullImg := ETEHelper.LoadImage(imgPath)

		for _, tile := range ETEHelper.SliceImageByGrid(fullImg, (tsxData.TileWidth+tsxData.TileHeight)/2) {
			tilesetsImages = append(tilesetsImages, ebiten.NewImageFromImage(tile))
		}
	}

	// 4. Charger les sprites
	for i, img := range tilesetsImages {
		mc.G.GetGame().Sprites[strconv.Itoa(i)] = img
	}

	// 5. Créer la map
	internalTiles := ETEHelper.ConvertTMXToInternal(tiledMap, jsonMaps.PropertiesForHeight)
	fmt.Printf("ConvertTMXToInternal returned %d tiles\n", len(internalTiles))
	if len(internalTiles) > 0 {
		fmt.Printf("First tile: Height=%d, X=%d, Y=%d, GID=%d\n",
			internalTiles[0].Height, internalTiles[0].X, internalTiles[0].Y, internalTiles[0].GID)
	}

	// Initialiser les éléments avec G
	elements := make(map[string]*Element)
	for k, v := range jsonMaps.Elements {
		v.G = mc.G
		elements[k] = v
	}

	resultat := Map{
		Map: MapData{
			Tiles: internalTilesToTiles(internalTiles, mc.G.GetGame()),
		},
		CellSize: jsonMaps.CellSize,
		Unité:    jsonMaps.Unite,
		Cam:      jsonMaps.Cam,
		Elements: elements, // ✅ Maintenant G est initialisé
	}

	return &resultat
}

func (g *Game) InitTile() {
	dir = strings.Join(strings.Split(g.Config.MapsPath, "/")[:len(strings.Split(g.Config.MapsPath, "/"))-1], "/")

	g.Tiles = make(map[string]*Tile)
	for i, j := range ETEHelper.JsonToStruct[map[string]*Tile](dir + "/" + g.MapConfig.Tiles) {
		tileId, _ := strconv.Atoi(i)
		tileId += 1
		g.Tiles[strconv.Itoa(tileId)] = j
		println("Loaded tile: " + strconv.Itoa(tileId))
	}
}

func (m *Map) GetSpriteByOrderYZX() map[int]map[[9]float32]*ebiten.Image { // [witdh/radius, height, xOffset, yOffset, xPos, yPos, xSize, ySize, rotation]
	resultat := make(map[int]map[[9]float32]*ebiten.Image)

	for k, v := range m.GetTileByLayer() {
		if resultat[k] == nil {
			resultat[k] = make(map[[9]float32]*ebiten.Image)
		}
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

	for _, es := range m.GetElementByLayer() {
		for _, e := range es {
			if resultat[e.Z] == nil {
				resultat[e.Z] = make(map[[9]float32]*ebiten.Image)
			}
			resultat[e.Z][[9]float32{
				e.Box[0], e.Box[1],
				e.Box[2], e.Box[3],
				e.Pos[0], e.Pos[1],
				float32(e.Size[0]), float32(e.Size[1]),
				e.Rotation,
			}] = e.GetSprite()
		}
	}

	// Sort keys
	keys := make([]int, 0, len(resultat))
	for k := range resultat {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Create new map with sorted keys
	sortedResult := make(map[int]map[[9]float32]*ebiten.Image)
	for _, k := range keys {
		sortedResult[k] = resultat[k]
	}

	return sortedResult
}

func (m *Map) GetElementByLayer() map[int][]Element {
	lay := make(map[int][]Element)
	for _, v := range m.Elements {
		(v).Animation = v.G.GetGame().Elements[v.Name].Animation
		(v).Size = v.G.GetGame().Elements[v.Name].Size
		(v).Box = v.G.GetGame().Elements[v.Name].Box
		(v).PushFrame()
		lay[v.Z] = append(lay[v.Z], *v)
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
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].Unité
				lay[la][[9]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0}] = img
			}
		}
	}

	return lay
}

func internalTilesToTiles(internalTiles []ETEHelper.TileInstance, G IForGame) map[int]map[[2]int]*TileElement {
	fmt.Printf("internalTilesToTiles called with %d tiles\n", len(internalTiles))
	result := make(map[int]map[[2]int]*TileElement)
	for _, t := range internalTiles {
		height := t.Height
		gid := int(t.GID)
		if result[height] == nil {
			result[height] = make(map[[2]int]*TileElement)
		}
		if G.GetGame().Tiles[strconv.Itoa(gid)] == nil {
			fmt.Printf("Tile %d not found\n", gid)
			continue
		}
		result[height][[2]int{t.X, t.Y}] = &TileElement{
			Id:   int(t.GID),
			Game: G,
		}
		result[height][[2]int{t.X, t.Y}].PushFrame()
	}
	fmt.Printf("internalTilesToTiles returning %d layers\n", len(result))
	return result
}
