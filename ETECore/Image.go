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
	for _, p := range g.Config.SpritePath {
		for _, fileName := range ETEHelper.GetAllFilesInDir(p) {
			g.Sprites[fileName] = ebiten.NewImageFromImage(ETEHelper.LoadImage(p + "/" + fileName))
			println("Loaded sprite: " + fileName)
		}
	}
}

func (g *Game) InitAnimations() {
	for _, p := range g.Config.AnimationsPath {
		Animations := ETEHelper.JsonToStruct[[]string](p)

		animDir := strings.Join(strings.Split(p, "/")[:len(strings.Split(p, "/"))-1], "/")

		g.Animations = make(map[string]*Animation)
		for _, animationName := range Animations {
			anim := ETEHelper.JsonToStruct[map[string]*Animation](animDir + "/" + animationName + ".json")
			for k, v := range anim {
				g.Animations[k] = v
			}
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

	if t.Game == nil || t.Game.GetGame() == nil || t.Game.GetGame().Animations == nil {
		return
	}

	if t.Frame > 0 {
		if float32(t.FFrame) >= t.Game.GetGame().Animations[tileDef.Animation].Frames[t.Frame-1].Duration/(t.Game.GetGame().Animations[tileDef.Animation].Speed*0.025) {
			t.Frame++
			t.FFrame = 0
		} else {
			t.FFrame++
		}
	} else {
		if float32(t.FFrame) >= t.Game.GetGame().Animations[tileDef.Animation].Frames[len(t.Game.GetGame().Animations[tileDef.Animation].Frames)-1].Duration/(t.Game.GetGame().Animations[tileDef.Animation].Speed*0.025) {
			t.Frame++
			t.FFrame = 0
		} else {
			t.FFrame++
		}
	}

	if t.Frame > len(t.Game.GetGame().Animations[tileDef.Animation].Frames) {
		t.Frame = 1
	}
}

func (t *Element) PushFrame() {
	if t.G == nil || t.G.GetGame() == nil || t.G.GetGame().Animations == nil {
		return
	}

	tileDef, exists := t.G.GetGame().Elements[t.Name]
	if !exists {
		fmt.Println("Element : " + t.Name + " not exist")
		return // Si l'élément n'est pas défini dans Elements.json, on ignore l'animation
	}

	if tileDef.Animation == "nil" {
		return
	}

	if t.Frame > 0 {
		if float32(t.FFrame) >= t.G.GetGame().Animations[tileDef.Animation].Frames[t.Frame-1].Duration/(t.G.GetGame().Animations[tileDef.Animation].Speed*0.025) {
			t.Frame++
			t.FFrame = 0
		} else {
			t.FFrame++
		}
	} else {
		if float32(t.FFrame) >= t.G.GetGame().Animations[tileDef.Animation].Frames[len(t.G.GetGame().Animations[tileDef.Animation].Frames)-1].Duration/(t.G.GetGame().Animations[tileDef.Animation].Speed*0.025) {
			t.Frame++
			t.FFrame = 0
		} else {
			t.FFrame++
		}
	}

	if t.Frame > len(t.G.GetGame().Animations[tileDef.Animation].Frames) {
		t.Frame = 1
	}
}

func (e *Element) GetSprite() *ebiten.Image {
	if e.G == nil {
		return nil
	}
	game := e.G.GetGame()
	if game == nil {
		fmt.Println("Game not found")
		return nil
	}
	if game.Animations == nil {
		fmt.Println("Animations not found")
		return nil
	}
	if game.Sprites == nil {
		fmt.Println("Sprites not found")
		return nil
	}
	if game.Animations[e.Animation] == nil {
		fmt.Printf("Animation %s not found\n", e.Animation)
		return nil
	}
	if game.Animations[e.Animation].Frames == nil {
		fmt.Printf("Animation frames not found for %s\n", e.Animation)
		return nil
	}

	spriteKey := ""
	if e.Frame-1 < 0 {
		spriteKey = e.G.GetGame().Animations[e.Animation].Frames[len(e.G.GetGame().Animations[e.Animation].Frames)-1].Frame
	} else {
		spriteKey = e.G.GetGame().Animations[e.Animation].Frames[e.Frame-1].Frame
	}
	sprite, exists := game.Sprites[spriteKey]
	if !exists {
		fmt.Println("Sprite not found for animation %s\n", e.Animation)
		return nil
	}
	if e.Frame == 0 {
		return sprite
	}

	// Cas animé : ajouter des nil checks
	if game.Elements == nil {
		return sprite
	}
	tileDef, exists := game.Elements[spriteKey]
	if !exists {
		return sprite
	}

	if game.Animations == nil {
		return sprite
	}
	anim, exists := game.Animations[tileDef.Animation]
	if !exists || anim == nil {
		return sprite
	}

	if e.Frame-1 >= len(anim.Frames) {
		return sprite
	}
	frame := anim.Frames[e.Frame-1]
	/*if frame == nil {
		return sprite
	}*/

	frameName := frame.Frame

	if float32(e.FFrame) < frame.Duration*(anim.Speed*0.025) {
		if e.Frame-2 >= 0 && e.Frame-2 < len(anim.Frames) {
			frameName = anim.Frames[len(anim.Frames)-1].Frame
		}
		e.FFrame = 0
	}

	frameSprite, exists := game.Sprites[frameName]
	if !exists {
		return sprite
	}

	return frameSprite
}

func (t *TileElement) GetSprite() (*ebiten.Image, bool) {
	if t.Game == nil {
		return nil, false
	}
	game := t.Game.GetGame()
	if game == nil || game.Sprites == nil {
		return nil, false
	}

	spriteKey := strconv.Itoa(int(t.Id - 1))
	sprite, exists := game.Sprites[spriteKey]
	if !exists {
		fmt.Printf("Sprite not found for GID %d\n", t.Id-1)
		return nil, false
	}
	if t.Frame == 0 {
		return sprite, false
	}

	// Cas animé : ajouter des nil checks
	if game.Tiles == nil {
		return sprite, false
	}
	tileDef, exists := game.Tiles[strconv.Itoa(int(t.Id))]
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

	frameName := frame.Frame

	if float32(t.FFrame) < frame.Duration*(anim.Speed*0.025) {
		if t.Frame-2 >= 0 && t.Frame-2 < len(anim.Frames) {
			frameName = anim.Frames[len(anim.Frames)-1].Frame
		}
		t.FFrame = 0
	}

	frameSprite, exists := game.Sprites[frameName]
	if !exists {
		return sprite, false
	}

	return frameSprite, true
}
