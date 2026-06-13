package ETE

import (
	"github.com/Try-si/ETE/ETECore"
	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Game *ETECore.Game
)

func Init(updateFunc func(float32) error, config string) {
	Game = &ETECore.Game{
		UpdateFunc: updateFunc,
	}

	if len(config) >= 5 && config[len(config)-5:] != ".json" {
		Game.Config = ETEHelper.StringToStruct[ETECore.Config](config)
	} else {
		Game.Config = ETEHelper.JsonToStruct[ETECore.Config](config)
	}

	if len(Game.Config.MapsPath) >= 5 && Game.Config.MapsPath[len(Game.Config.MapsPath)-5:] != ".json" {
		Game.MapConfig = ETEHelper.StringToStruct[ETECore.MapConfig](Game.Config.MapsPath)
	} else {
		Game.MapConfig = ETEHelper.JsonToStruct[ETECore.MapConfig](Game.Config.MapsPath)
	}
	Game.MapConfig.G = Game

	Game.InitSprites()
	Game.InitMap()
	Game.InitTile()
	Game.InitAnimations()
	Game.InitElements()
}

func GameLoop() {
	ebiten.SetWindowSize(int(Game.Config.ScreenWidth), int(Game.Config.ScreenHeight))
	ebiten.SetWindowTitle(Game.Config.Title)
	if err := ebiten.RunGame(Game); err != nil {
		panic(err)
	}
}

func GetGame() *ETECore.Game {
	return Game
}
