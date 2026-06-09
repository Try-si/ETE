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

	Debug      bool
	UpdateFunc func(float32) error
}

type Config struct {
	ScreenWidth  int    `json:"ScreenWidth"`  // largeur de l'écran
	ScreenHeight int    `json:"ScreenHeight"` // hauteur de l'écran
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
}

type Map struct {
	Map      MapData
	CellSize int
	Unité    float32
	Cam      Camera
	Elements map[string]*Element
	G        IForGame
}

type MapData struct {
	Tiles   map[int]map[[2]int]TileElement
	Tileset []*ebiten.Image
}

type TileElement struct {
	Id    int
	Frame int
	Game  IForGame
}

type Tile struct {
	Animation string   `json:"Animation"`
	Collision bool     `json:"Collision"`
	Box       [4]int   `json:"Box"`
	Tags      []string `json:"Tags"`
}

type Element struct {
	Animation string     `json:"Animation"`
	Size      [2]int     `json:"Size"`
	Box       [4]float32 `json:"Box"`
	Tags      []string   `json:"Tags"`

	Name         string            `json:"Name"`
	Pos          [2]float32        `json:"Pos"`
	Rotation     float32           `json:"Rotation"`
	Layer, Frame int               `json:"Layer"`
	MetaData     map[string]string `json:"MetaData"`
	G            IForGame
}

type JsonMap struct {
	Map      string              `json:"Map"`
	CellSize int                 `json:"CellSize"`
	Unite    float32             `json:"Unite"`
	Cam      Camera              `json:"Cam"`
	Elements map[string]*Element `json:"Elements"`
}

type Camera struct {
	Zoom   float32    `json:"Zoom"`
	Offset [2]float32 `json:"Offset"`
}

// image

type Animation struct {
	Frames []Frame `json:"Frames"`
	Speed  int     `json:"Speed"`
}

type Frame struct {
	Frame    string `json:"Frame"`
	Duration int    `json:"Duration"`
	Box      [4]int `json:"Box"`
}
