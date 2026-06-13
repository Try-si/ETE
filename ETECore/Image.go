package ETECore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
)

type IForGame interface {
	GetGame() *Game
}

func (g *Game) InitSprites() {
	g.Sprites = make(map[string]*ebiten.Image)
	for _, fileName := range ETEHelper.GetAllFilesInDir(g.Config.SpritePath) {
		g.Sprites[fileName] = ebiten.NewImageFromImage(ETEHelper.LoadImage(g.Config.SpritePath + "/" + fileName))
		println("Loaded sprite: " + fileName)
	}
}

func (g *Game) InitAnimations() {
	Animations := ETEHelper.JsonToStruct[[]string](g.Config.AnimationsPath)

	animDir := strings.Join(strings.Split(g.Config.AnimationsPath, "/")[:len(strings.Split(g.Config.AnimationsPath, "/"))-1], "/")

	g.Animations = make(map[string]*Animation)
	for _, animationName := range Animations {
		anim := ETEHelper.JsonToStruct[map[string]*Animation](animDir + "/" + animationName + ".json")
		for k, v := range anim {
			g.Animations[k] = v
		}
	}
}

func (g *Game) InitElements() {
	g.Elements = ETEHelper.JsonToStruct[map[string]Element](dir + "/" + g.MapConfig.Elements)

	for k, _ := range g.Elements {
		newElement := g.Elements[k]
		newElement.G = g
		g.Elements[k] = newElement
	}
}

func (t *TileElement) PushFrame() {
	tileDef, exists := t.Game.GetGame().Tiles[strconv.Itoa(int(t.Id))]
	if !exists {
		fmt.Println("Tile : " + strconv.Itoa(int(t.Id)) + " not exist")
		return // Si le tile n'est pas défini dans Tiles.json, on ignore l'animation
	}

	if tileDef.Animation == "nil" {
		return
	}

	if t.Frame > len(t.Game.GetGame().Animations[tileDef.Animation].Frames) {
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
	if t.Game == nil {
		return nil, false
	}
	game := t.Game.GetGame()
	if game == nil || game.Sprites == nil {
		return nil, false
	}

	spriteKey := strconv.Itoa(int(t.Id))
	sprite, exists := game.Sprites[spriteKey]
	if !exists {
		fmt.Printf("Sprite not found for GID %d\n", t.Id)
		return nil, false
	}
	if t.Frame == 0 {
		return sprite, false
	}

	// Cas animé : ajouter des nil checks
	if game.Tiles == nil {
		return sprite, false
	}
	tileDef, exists := game.Tiles[spriteKey]
	if !exists || tileDef == nil {
		return sprite, false
	}

	if game.Animations == nil {
		return sprite, false
	}
	anim, exists := game.Animations[tileDef.Animation]
	if !exists || anim == nil {
		return sprite, false
	}

	if t.Frame-1 >= len(anim.Frames) {
		return sprite, false
	}
	frame := anim.Frames[t.Frame-1]
	/*if frame == nil {
		return sprite, false
	}*/

	frameSprite, exists := game.Sprites[frame.Frame]
	if !exists {
		return sprite, false
	}

	return frameSprite, true
}
