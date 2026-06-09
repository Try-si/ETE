package ETECore

import (
	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
)

type IForGame interface {
	GetGame() *Game
}

func (g *Game) InitSprites() {
	for _, fileName := range ETEHelper.GetAllFilesInDir(g.Config.SpritePath) { // loadImg
		g.Sprites[fileName] = ebiten.NewImageFromImage(ETEHelper.LoadImage(g.Config.SpritePath + "/" + fileName))
	}

	for _, m := range g.Maps { // loadTileset
		for i, tile := range m.Map.Tileset {
			g.Sprites[string(i)] = tile
		}
	}
}

func (g *Game) InitAnimations() {
	g.Animations = ETEHelper.JsonToStruct[map[string]*Animation](g.Config.AnimationsPath)
}

func (g *Game) InitElements() {
	g.Elements = ETEHelper.JsonToStruct[map[string]Element](g.MapConfig.Elements)
}

func (t *TileElement) PushFrame() {
	if t.Game.GetGame().Tiles[string(t.Id)].Animation == "nil" {
		return
	}
	if t.Frame > len(t.Game.GetGame().Animations[t.Game.GetGame().Tiles[string(t.Id)].Animation].Frames) {
		t.Frame = 0
	}
	t.Frame++
}

func (t *Element) PushFrame() {
	if t.Frame >= len(t.G.GetGame().Animations[t.Animation].Frames) {
		t.Frame = 0
	} else {
		t.Frame++
	}
}

func (e *Element) GetSprite() *ebiten.Image {
	return e.G.GetGame().Sprites[e.G.GetGame().Animations[e.Animation].Frames[e.Frame].Frame]
}

func (t *TileElement) GetSprite() (*ebiten.Image, bool) {
	if t.Frame == 0 {
		return t.Game.GetGame().Sprites[string(t.Id)], false
	}
	return t.Game.GetGame().Sprites[t.Game.GetGame().Animations[t.Game.GetGame().Tiles[string(t.Id)].Animation].Frames[t.Frame].Frame], true
}
