package ETECore

import "github.com/hajimehoshi/ebiten/v2"

// All

type Game struct {
	Config    Config
	MapConfig MapConfig

	Elements   map[string]Element
	Maps       map[string]*Map
	Sprites    map[string]*ebiten.Image
	Tiles      map[string]*Tile
	Animations map[string]*Animation

	Debug, Quite bool
	UpdateFunc   func(float32) error
}

func (g *Game) GetGame() *Game {
	return g
}

type Config struct {
	ScreenWidth  int    `json:"ScreenWidth"`  // largeur de l'écran
	ScreenHeight int    `json:"ScreenHeight"` // hauteur de l'écran
	Resizeable   bool   `json:"Resizeable"`   // si la fenêtre peut être redimensionnée
	Title        string `json:"Title"`        // titre de la fenêtre
	Map          string `json:"Map"`          // map actuelle/de base

	SpritePath     string `json:"SpritePath"`     // chemin vers les sprites
	MapsPath       string `json:"MapsPath"`       // chemin vers les maps
	AnimationsPath string `json:"AnimationsPath"` // chemin vers les animations
}

// Map

type MapConfig struct {
	Maps            []string `json:"Maps"`
	JsonMap         string   `json:"JsonMap"`
	TiledMap        string   `json:"TiledMap"`
	Elements        string   `json:"Elements"`
	Tiles           string   `json:"Tiles"`
	Parrallax       bool     `json:"Parrallax"`
	ParrallaxFactor float32  `json:"ParrallaxFactor"`
	G               IForGame
}

type Map struct {
	Map      MapData
	CellSize int
	Unité    int
	Cam      Camera
	Elements map[string]*Element
}

type MapData struct {
	Tiles map[int]map[[2]int]*TileElement
}

type TileElement struct {
	Id     int
	Frame  int
	FFrame int
	Game   IForGame
}

type Tile struct {
	Animation string   `json:"Animation"`
	Box       [4]int   `json:"Box"`
	Tags      []string `json:"Tags"`
}

type Element struct {
	Animation string     `json:"Animation"`
	Size      [2]int     `json:"Size"`
	Box       [4]float32 `json:"Box"`
	Tags      []string   `json:"Tags"`

	Name             string            `json:"Name"`
	Pos              [2]float32        `json:"Pos"`
	Rotation         float32           `json:"Rotation"`
	Z, Frame, FFrame int               `json:"Height"`
	MetaData         map[string]string `json:"MetaData"`
	G                IForGame
}

type JsonMap struct {
	Map                 string              `json:"Map"`
	CellSize            int                 `json:"CellSize"`
	Unite               int                 `json:"Unite"`
	PropertiesForHeight string              `json:"PropertiesForHeight"`
	Cam                 Camera              `json:"Cam"`
	Elements            map[string]*Element `json:"Elements"`
}

type Camera struct {
	Z      float32    `json:"Z"`
	Offset [2]float32 `json:"Offset"`
}

// image

type Animation struct {
	Frames []Frame `json:"Frames"`
	Speed  float32 `json:"Speed"`
}

type Frame struct {
	Frame    string  `json:"Frame"`
	Duration float32 `json:"Duration"`
	Box      [4]int  `json:"Box"`
}
