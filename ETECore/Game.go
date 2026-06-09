package ETECore

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Update() error {
	return g.UpdateFunc(float32(ebiten.ActualFPS()))
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, L := range g.Maps[g.Config.Map].GetSpriteByOrderYZX() {
		for Box, img := range L { // Box = [witdh/radius, height, xOffset, yOffset, xPos, yPos, xSize, ySize, rotation]
			if g.Debug { // si le mode debug est activé
				posX := (Box[4] - g.Maps[g.Config.Map].Cam.Offset[0]) * g.Maps[g.Config.Map].Unité // calculer la position x en pixels
				posY := (Box[5] + g.Maps[g.Config.Map].Cam.Offset[1]) * g.Maps[g.Config.Map].Unité // calculer la position y en pixels

				whith := Box[0] * g.Maps[g.Config.Map].Unité  // obtenir la largeur de la hitbox en pixels
				height := Box[1] * g.Maps[g.Config.Map].Unité // obtenir la hauteur de la hitbox en pixels

				if Box[2] != 0 { // si xOffset est différent de 0
					posX += Box[2] * g.Maps[g.Config.Map].Unité
				}
				if Box[3] != 0 { // si yOffset est différent de 0
					posY += Box[3] * g.Maps[g.Config.Map].Unité
				}

				if whith == 0 && height == 0 { // si la hitbox n'est pas définie
					continue
				} else if height == 0 { // si la hitbox est un cercle
					//Draw circle
					drawCircle(screen, posX, posY, whith, color.RGBA{255, 255, 255, 128})
				} else { // si la hitbox est un rectangle
					//Draw rectangle
					drawRect(screen, posX, posY, whith, height, color.RGBA{255, 255, 255, 128})
				}
			}

			opts := &ebiten.DrawImageOptions{}

			// 1. Centrer sur l'origine (avant scale)
			width := float32(img.Bounds().Dx())  // largeur de l'image en pixels
			height := float32(img.Bounds().Dy()) // hauteur de l'image en pixels
			if Box[6] != 0 && Box[7] != 0 {      // si la taille est définie
				width = Box[6] * g.Maps[g.Config.Map].Unité  // largeur en pixels
				height = Box[7] * g.Maps[g.Config.Map].Unité // hauteur en pixels
			}
			opts.GeoM.Translate(float64(-width/2), float64(-height/2)) // centrer sur l'origine

			// 2. Scale (avec zoom)
			if Box[6] != 0 && Box[7] != 0 { // si la taille est définie
				opts.GeoM.Scale(float64(Box[6]*g.Maps[g.Config.Map].Unité)/float64(img.Bounds().Dx()), float64(Box[7]*g.Maps[g.Config.Map].Unité)/float64(img.Bounds().Dy()))
				// scale with element size : element.Size = taille en unité, * g.Maps[g.Conf.Map].Unité = mettre taille en pixels, / img.Bounds().Dx() = scale
			} else {
				opts.GeoM.Scale(float64(g.Maps[g.Config.Map].Unité), float64(g.Maps[g.Config.Map].Unité))
			}

			// 3. Rotate
			opts.GeoM.Rotate(float64(Box[8])) // rotate

			// 4. Translate vers la position finale (sans zoom dans la translation)
			opts.GeoM.Translate(float64(Box[4])*float64(g.Maps[g.Config.Map].Unité), float64(Box[5])*float64(g.Maps[g.Config.Map].Unité)) // translate

			// 5. Camera offset (avec zoom)
			opts.GeoM.Scale(float64(g.Maps[g.Config.Map].Cam.Zoom), float64(g.Maps[g.Config.Map].Cam.Zoom))                                                                                       // Zoom
			opts.GeoM.Translate(float64(g.Maps[g.Config.Map].Cam.Offset[0])*float64(g.Maps[g.Config.Map].Unité), float64(g.Maps[g.Config.Map].Cam.Offset[1])*float64(g.Maps[g.Config.Map].Unité)) // Center

			screen.DrawImage(img, opts) // dessiner l'image
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.Config.ScreenWidth), int(g.Config.ScreenHeight)
}

func drawRect(screen *ebiten.Image, x, y, width, height float32, clr color.Color) {
	vector.FillRect(screen, x-width/2, y-height/2, width, height, clr, false)
}
func drawCircle(screen *ebiten.Image, x, y, radius float32, clr color.Color) { // dessiner un cercle
	centerX, centerY := int(x), int(y) // centre du cercle
	r := int(radius)                   // rayon du cercle

	for dy := -r; dy <= r; dy++ { // itérer sur tous les pixels du cercle
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r { // vérifier si le pixel est dans le cercle
				screen.Set(centerX+dx, centerY+dy, clr) // dessiner le pixel
			}
		}
	}
}
