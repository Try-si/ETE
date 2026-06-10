package ETEHelper

import "encoding/xml"

// --- Fichier .TMX (La Carte) ---

type TMXMap struct {
	XMLName      xml.Name `xml:"map"`
	Version      string   `xml:"version,attr"`
	TiledVersion string   `xml:"tiledversion,attr"`
	Orientation  string   `xml:"orientation,attr"`
	RenderOrder  string   `xml:"renderorder,attr"`
	Width        int      `xml:"width,attr"`
	Height       int      `xml:"height,attr"`
	TileWidth    int      `xml:"tilewidth,attr"`
	TileHeight   int      `xml:"tileheight,attr"`
	Infinite     bool     `xml:"infinite,attr"` // true si carte infinie
	NextLayerID  int      `xml:"nextlayerid,attr"`
	NextObjectID int      `xml:"nextobjectid,attr"`

	// Les tuiles sont référencées, pas définies ici
	Tilesets []TMXTilesetRef `xml:"tileset"`

	// Les calques (Layers) contiennent les données
	Layers []TMXLayer `xml:"layer"`

	// Groups (optionnel, si vous utilisez des Group Layers)
	Groups []TMXGroup `xml:"group"`
}

type TMXTilesetRef struct {
	FirstGID int    `xml:"firstgid,attr"` // CRUCIAL: Offset des IDs
	Source   string `xml:"source,attr"`   // Chemin vers le .tsx
	// Si le tileset est intégré directement (pas de source), on aurait ici la structure complète
}

type TMXLayer struct {
	ID         int          `xml:"id,attr"`
	Name       string       `xml:"name,attr"`
	Width      int          `xml:"width,attr"`
	Height     int          `xml:"height,attr"`
	Data       TMXLayerData `xml:"data"`
	Chunks     []TMXChunk   `xml:"chunk"` // Pour les cartes infinies
	Properties []Property   `xml:"properties>property"`
}

// Structure pour un Chunk (Carte Infinie)
// NOTE: X et Y sont des INT (signés) pour accepter les valeurs négatives (-16, -32, etc.)
type TMXChunk struct {
	X      int          `xml:"x,attr"`
	Y      int          `xml:"y,attr"`
	Width  int          `xml:"width,attr"`
	Height int          `xml:"height,attr"`
	Data   TMXLayerData `xml:"data"`
}

type TMXLayerData struct {
	Encoding string `xml:"encoding,attr"` // "csv", "base64"
	// Pour CSV, on va parser le contenu texte manuellement ou via un champ chardata
	RawContent string `xml:",chardata"`
	// Si Base64, on utilisera un champ []byte
	// Content []byte `xml:",chardata"`
}

// --- Fichier .TSX (Le Tileset) ---

type TSXTileset struct {
	XMLName    xml.Name `xml:"tileset"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`

	// L'image source
	Image TSXImage `xml:"image"`

	// Définitions spécifiques de tuiles (collisions, animations, props)
	Tiles []TSXTile `xml:"tile"`
}

type TSXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Trans  string `xml:"trans,attr"` // Couleur de transparence (hex)
}

type TSXTile struct {
	ID         int        `xml:"id,attr"`
	Type       string     `xml:"type,attr"`
	Properties []Property `xml:"properties>property"`
	// On peut ajouter ObjectGroup (collisions) et Animation ici si besoin
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"` // "int", "float", "bool", "string"
}

// Groupe de calques (si utilisé)
type TMXGroup struct {
	Name   string     `xml:"name,attr"`
	Layers []TMXLayer `xml:"layer"`
}
