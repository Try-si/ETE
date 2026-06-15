# Description

ETM (Ebien Tiled Moteur) est un helper & framework de base pour les jeux utilisant le moteur Ebiten et les maps Tiled.

## Installation

```bash
go get github.com/Try-si/ETM
```

## Utilisation

exemple d'architécture :

```
 Maps/
    Maps/ *
        overworld.tmx // pas obligatoire (c'est un exemple)
        vos map tiled (.tmx)
    Json/ *
        overworld.json // pas obligatoire (c'est un exemple)
        vos json de config de map
    TileMaps/
        vos tilemaps tiled (.tsx)
    Tiles.json
    Elements.json
    Maps.json
 Textures/
    Player_idle_1.png // pas obligatoire (c'est un exemple)
    Player_idle_2.png // pas obligatoire (c'est un exemple)
    Player_idle_3.png // pas obligatoire (c'est un exemple)
    Water_idle_1.png // pas obligatoire (c'est un exemple)
    Water_idle_2.png // pas obligatoire (c'est un exemple)
    Water_idle_3.png // pas obligatoire (c'est un exemple)
    vos textures
 Animations/
    Player.json // pas obligatoire (c'est un exemple)
    Water.json // pas obligatoire (c'est un exemple)
    Animations.json
    vos animations (.json)
config.json
main.go

* = les nom des fichiers sont les mêmes
```

main.go :

```go
package main

import (
    "github.com/Try-si/ETM"
)

func main() {
    ETM.Init(func(deltaTime float64) error {
        // votre code ici
        return nil
    }, "config.json")
    
    ETM.GameLoop()
}
```

config.json :

```json
{
    "ScreenWidth": 800,          // largeur de l'écran
    "ScreenHeight": 600,         // hauteur de l'écran
    "Title": "Test",             // titre de la fenêtre
    "Map": "Overworld",          // map actuelle/de base

    "SpritePath": "Textures",    // chemin vers les sprites
    "MapsPath": "Maps/Maps.json", // chemin vers les maps
    "AnimationsPath": "Animations/Animations.json" // chemin vers les animations
}
```

Maps.json : 

```json
{
    "Maps": ["Overworld"],
    "JsonMap": "Json",
    "TiledMap": "Maps",
    "Elements": "Elements.json",
    "Tiles": "Tiles.json",
    "Parrallax": true,
    "ParrallaxFactor": 1.0
}
```

Elements.json :

```json
{
    "Player": {
        "Size": [32, 32],
        "Box": [0, 0, 0, 0], // witdh, height (si il est == a 0 alors c'est un cercle et witdh = rayon), box pos x, box pos y
        "Tags": ["Player"]
    }
}
```

Tiles.json :

```json
{
    "0":{// is a tile id in map of tiled of grass
        "Animation": "nil",// cela peut etre nil
        "Box": [32,32,0,0],
        "Tags": ["Grass"]
    },
    "1":{// is a tile id in map of tiled of water
        "Animation": "Water_idle",
        "Box": [32,32,0,0],
        "Tags": ["Water"]
    }
}
```

Animations.json :

```json
["Player", "Water"]
```

Exemple map json (overworld.json) :

```json
{
    "Map": "Overworld",         // nom de la map
    "CellSize": 1,              // taille de la cellule en unités
    "Unité": 32,                // taille d'une unité en pixels
    "PropertiesForHeight": "Height", // nom de la propriété dans tiled qui definit la hauteur du layer
    "Cam": {
        "Z": 1.0,
        "Offset": [0.0, 0.0]
    },

    "Elements": {
        "Player": {
            "Animation": "Player_idle",// cela ne peut pas etre nil
            "Name": "Player",  // nom de l'élément dans Elements.json
            "Pos": [0.0, 0.0], // position de l'élément
            "Rotation": 0,     // rotation de l'élément
            "Height": 1,        // hauteur de l'élément
            "MetaData": {
                "Nom de la variable": "valeur de la variable" // toujours une string
            }
        }
    }       // éléments dans la map
}
```
Exemple animation json (Player.json) :

```json
{
    "Player_idle": {
        "Frames": [
            {"Frame": "Player_idle_1.png", "Duration": "rand", "Box": [32, 32, 0, 0]},
            {"Frame": "Player_idle_2.png", "Duration": "rand", "Box": [32, 32, 0, 0]},
            {"Frame": "Player_idle_3.png", "Duration": "rand", "Box": [32, 32, 0, 0]}
        ],
        "Speed": 1
    }
}
```

Exemple animation json (Water.json) :

```json
{
    "Water_idle": {
        "Frames": [
            {"Frame": "Water_idle_1.png", "Duration": 1, "Box": [32, 32, 0, 0]},
            {"Frame": "Water_idle_2.png", "Duration": 1, "Box": [32, 32, 0, 0]},
            {"Frame": "Water_idle_3.png", "Duration": 1, "Box": [32, 32, 0, 0]}
        ],
        "Speed": 1
    }
}
```
