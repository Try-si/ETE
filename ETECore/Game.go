package ETECore

import (
	"image/color"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Update() error {
	return g.UpdateFunc(float32(ebiten.ActualFPS()))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for height, L := range g.Maps[g.Config.Map].GetSpriteByOrderYZX() {

		if float32(height) < g.Maps[g.Config.Map].Cam.Z {
			continue
		}

		for Box, img := range L { // Box = [witdh/radius, height, xOffset, yOffset, xPos, yPos, xSize, ySize, rotation]
			unité := float64(g.Maps[g.Config.Map].Unité)
			if img == nil {
				continue
			}
			zoom := g.Maps[g.Config.Map].Cam.Z
			screenCenterX := float64(g.Config.ScreenWidth) / 2
			screenCenterY := float64(g.Config.ScreenHeight) / 2

			// Reproduire la séquence de transformations du rendu normal
			baseX := float64(Box[4])*unité + screenCenterX
			baseY := float64(Box[5])*unité + float64(g.Config.ScreenHeight) // Calculer le facteur de parallaxe
			parallaxFactor := float32(1.0)
			if g.MapConfig.Parrallax && g.Maps[g.Config.Map].Cam.Z > 0 {
				// Les layers plus "hauts" (height élevé) bougent moins
				// Les layers plus "bas" (height faible) bougent plus
				parallaxFactor = 1.0 - (float32(height)/g.Maps[g.Config.Map].Cam.Z)*g.MapConfig.ParrallaxFactor
				if parallaxFactor < 0 {
					parallaxFactor = 0
				}
			}

			// Appliquer offset caméra avec parallaxe
			baseX -= float64(g.Maps[g.Config.Map].Cam.Offset[0]) * unité * float64(parallaxFactor)
			baseY += float64(g.Maps[g.Config.Map].Cam.Offset[1]) * unité * float64(parallaxFactor)

			// Zoom vers le centre
			posX := -(baseX-screenCenterX)/float64(zoom) + screenCenterX
			posY := -(baseY-screenCenterY)/float64(zoom) + screenCenterY

			// xOffset/yOffset avec zoom
			if Box[2] != 0 {
				posX += float64(Box[2]) * unité / float64(zoom)
			}
			if Box[3] != 0 {
				posY += float64(Box[3]) * unité / float64(zoom)
			}
			if g.Debug { // si le mode debug est activé

				whith := float64(Box[0]) / float64(zoom)
				height := float64(Box[1]) / float64(zoom)

				if whith == 0 && height == 0 { // si la hitbox n'est pas définie
					continue
				} else if height == 0 { // si la hitbox est un cercle
					//Draw circle
					drawCircle(screen, float32(posX), float32(posY), float32(whith), ETEHelper.ImgMoyenne(*img))
				} else { // si la hitbox est un rectangle
					//Draw rectangle
					drawRect(screen, float32(posX), float32(posY), float32(whith), float32(height), ETEHelper.ImgMoyenne(*img))
				}

				drawText(screen, float32(posX), float32(posY), ETEHelper.GetKey(g.Sprites, img))
			} else {

				opts := &ebiten.DrawImageOptions{}

				if Box[6] != 0 && Box[7] != 0 { // si la taille est définie
					opts.GeoM.Scale(float64(-Box[6])/float64(img.Bounds().Dx()), -float64(Box[7])/float64(img.Bounds().Dy()))
					// scale with element size : element.Size = taille en unité, * g.Maps[g.Conf.Map].Unité = mettre taille en pixels, / img.Bounds().Dx() = scale
				} else {
					opts.GeoM.Scale(-unité, -unité)
				}

				opts.GeoM.Rotate(float64(Box[8])) // rotate

				opts.GeoM.Scale(float64(g.Maps[g.Config.Map].Cam.Z), float64(g.Maps[g.Config.Map].Cam.Z))

				opts.GeoM.Translate(posX, posY)
				screen.DrawImage(img, opts) // dessiner l'image
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.Config.ScreenWidth), int(g.Config.ScreenHeight)
}

func drawRect(screen *ebiten.Image, x, y, width, height float32, clr color.Color) {
	vector.FillRect(screen, x-width/2, y-height/2, width, height, clr, false)
	vector.StrokeRect(screen, x-width/2, y-height/2, width, height, 1, clr, false)
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
func drawText(screen *ebiten.Image, x, y float32, text string) {
	ebitenutil.DebugPrintAt(screen, text, int(x), int(y))
}
