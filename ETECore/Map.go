package ETECore

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
)

var dir string
var fid int

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
		tsxAbsPath := dir + tsRef.Source
		tsxAbsPath = strings.ReplaceAll(tsxAbsPath, "..", "")
		tsxData, err := ETEHelper.LoadTSX(tsxAbsPath)
		if err != nil {
			panic(fmt.Sprintf("Erreur tileset %s: %v", tsRef.Source, err))
		}

		// Construire le chemin correct pour l'image
		// L'image est ../../Textures/tileset_1.png relatif au TSX
		// Donc le chemin absolu est Textures/tileset_1.png
		imgPath := tsxData.Image.Source
		imgPath = strings.ReplaceAll(imgPath, "..", "")
		for strings.Contains(imgPath, "//") {
			imgPath = strings.ReplaceAll(imgPath, "//", "")
		}
		fullImg := ETEHelper.LoadImage(imgPath)

		for _, tile := range ETEHelper.SliceImageByGrid(fullImg, (tsxData.TileWidth+tsxData.TileHeight)/2) {
			tilesetsImages = append(tilesetsImages, ebiten.NewImageFromImage(tile))
		}
	}

	// 4. Charger les sprites
	for i, img := range tilesetsImages {
		mc.G.GetGame().Sprites[strconv.Itoa(fid+i)] = img
	}
	fid += len(tilesetsImages)

	// 5. Créer la map
	internalTiles := ETEHelper.ConvertTMXToInternal(tiledMap, jsonMaps.PropertiesForHeight)

	// Initialiser les éléments avec G
	elements := make(map[string]*Element)
	for k, v := range jsonMaps.Elements {
		//fmt.Println("Pré Element "+k+"Z:", jsonMaps.Elements[k].Z)
		//fmt.Println("Element "+k+" Z:", v.Z)
		v.G = mc.G
		v.Rand = time.Now().Hour() + time.Now().Minute() + time.Now().Second() + time.Now().Nanosecond()
		v.Pos[1] = -v.Pos[1]
		elements[k] = v
	}

	if jsonMaps.Cam.DebZ == 0 {
		jsonMaps.Cam.DebZ = jsonMaps.Cam.Z
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

var resulta map[int]bool

func (m *Map) GetSpriteByOrderYZX() []HeightLayer {
	type SpriteData struct {
		Height  float32
		Box     [9]float32
		Img     *ebiten.Image
		Visible bool
		Paralax bool
		Rand    int
	}

	var result []SpriteData

	// Ajouter les tiles
	for k, v := range m.GetTileByLayer() {
		for pos, tile := range v {
			result = append(result, SpriteData{
				Height: k,
				Box: [9]float32{
					float32(pos[0]), float32(pos[1]),
					float32(pos[2]), float32(pos[3]),
					float32(pos[4]), float32(pos[5]),
					float32(pos[6]), float32(pos[7]),
					float32(pos[8]),
				},
				Img:     tile,
				Visible: true,
				Rand:    pos[9],
			})
		}
	}

	// Ajouter les éléments
	for _, es := range m.GetElementByLayer() {
		for _, e := range es {
			result = append(result, SpriteData{
				Height: e.Z,
				Box: [9]float32{
					e.Box[0], e.Box[1],
					e.Box[2], e.Box[3],
					e.Pos[0], e.Pos[1],
					float32(e.Size[0]), float32(e.Size[1]),
					e.Rotation,
				},
				Img:     e.GetSprite(),
				Visible: e.Visible,
				Rand:    e.Rand,
				Paralax: e.Parallax,
			})
		}
	}

	// Trier par hauteur (du plus profond au plus proche)
	sort.Slice(result, func(i, j int) bool {
		//fmt.Println("Height i: ", result[i].Height, " Height j: ", result[j].Height)
		if result[i].Height != result[j].Height {
			return result[i].Height > result[j].Height // profond -> proche, inchangé
		}
		if result[i].Box[5] != result[j].Box[5] {
			return result[i].Box[5] > result[j].Box[5] // Y croissant (Box[4]=X, Box[5]=Y pour tiles et éléments)
		}
		if result[i].Box[4] != result[j].Box[4] {
			return result[i].Box[4] > result[j].Box[4] // X croissant en tie-break
		}
		return result[i].Rand > result[j].Rand
	})

	// Regrouper par hauteur réelle, SANS utiliser Box comme clé
	var resultat []HeightLayer
	var currentLayer []SpriteEntry
	first := true
	var currentHeight float32

	for _, sprite := range result {
		if first || sprite.Height != currentHeight {
			if !first {
				resultat = append(resultat, HeightLayer{
					Height:  currentHeight,
					Sprites: currentLayer,
				})
			}
			currentHeight = sprite.Height
			currentLayer = make([]SpriteEntry, 0) // Réinitialiser à chaque changement de hauteur
			first = false
			//fmt.Println("Height: ", currentHeight)
		}
		currentLayer = append(currentLayer, SpriteEntry{
			Box:     sprite.Box,
			Img:     sprite.Img,
			Visible: sprite.Visible,
			Paralax: sprite.Paralax,
		})
	}
	if !first {
		resultat = append(resultat, HeightLayer{
			Height:  currentHeight,
			Sprites: currentLayer,
		})
	}

	sort.Slice(resultat, func(i, j int) bool {
		return resultat[i].Height > resultat[j].Height
	})

	return resultat
}

type SpriteEntry struct {
	Box     [9]float32
	Img     *ebiten.Image
	Visible bool
	Paralax bool
}

type HeightLayer struct {
	Height  float32
	Sprites []SpriteEntry
}

func (m *Map) GetElementByLayer() map[float32][]Element {
	lay := make(map[float32][]Element)
	for _, v := range m.Elements {
		//fmt.Println("Element "+n+" Z: ", v.Z)
		(v).Animation = v.G.GetGame().Elements[v.Name].Animation
		(v).Size = v.G.GetGame().Elements[v.Name].Size
		(v).Box = v.G.GetGame().Elements[v.Name].Box
		(v).PushFrame()
		lay[v.Z] = append(lay[v.Z], *v)
	}

	return lay
}

func (m *Map) GetTileByLayer() map[float32]map[[10]int]*ebiten.Image {
	lay := make(map[float32]map[[10]int]*ebiten.Image)
	for la, tiles := range m.Map.Tiles {
		lay[la] = make(map[[10]int]*ebiten.Image)
		for pos, tile := range tiles {
			//fmt.Println("Tile "+strconv.Itoa(int(tile.Id))+" layer Z : ", la)
			tile.PushFrame()
			img, isAnimated := tile.GetSprite()
			if img == nil {
				img = ebiten.NewImage(32, 32)
				img.Fill(color.RGBA{255, 6, 181, 255})
			}
			if isAnimated {
				b := tile.Game.GetGame().Animations[tile.Game.GetGame().Tiles[strconv.Itoa(tile.Id)].Animation].Frames[tile.Frame-1].Box
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].CellSize
				lay[la][[10]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0, tile.Rand}] = img
			} else {
				var b [4]int
				if _, exist := tile.Game.GetGame().Tiles[strconv.Itoa(int(tile.Id))]; !exist {
					fmt.Printf("Tile box: %d not exist\n", tile.Id)
					b = [4]int{0, 0, 0, 0}
				} else {
					b = tile.Game.GetGame().Tiles[strconv.Itoa(int(tile.Id))].Box
				}
				cs := tile.Game.GetGame().Maps[tile.Game.GetGame().Config.Map].Unité
				lay[la][[10]int{b[0], b[1], b[2], b[3], pos[0], pos[1], cs, cs, 0, tile.Rand}] = img
			}
		}
	}

	return lay
}

func internalTilesToTiles(internalTiles []ETEHelper.TileInstance, G IForGame) map[float32]map[[2]int]*TileElement {
	fmt.Printf("internalTilesToTiles called with %d tiles\n", len(internalTiles))
	result := make(map[float32]map[[2]int]*TileElement)
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
			Rand: time.Now().Hour() + time.Now().Minute() + time.Now().Second() + time.Now().Nanosecond(),
		}
		result[height][[2]int{t.X, t.Y}].PushFrame()
	}
	fmt.Printf("internalTilesToTiles returning %d layers\n", len(result))
	return result
}
