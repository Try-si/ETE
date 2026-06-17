package ETECore

import (
	"errors"
	"image/color"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Update() error {
	if g.Quite {
		return errors.New("quit")
	}
	return g.UpdateFunc(float32(ebiten.ActualFPS()))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	Map := g.Maps[g.Config.Map]
	unit := float32(Map.Unité)

	for height, L := range Map.GetSpriteByOrderYZX() {
		if float32(height) < Map.Cam.Z {
			continue
		}

		for Box, img := range L {
			zoom := Map.Cam.Z

			// === CALCUL DE LA TAILLE ===
			spriteWidth := float32(Box[0]) * unit * zoom
			spriteHeight := float32(Box[1]) * unit * zoom

			// === CALCUL DE LA POSITION ===
			worldX := Box[4] * unit * zoom
			worldY := Box[5] * unit * zoom

			centerX := float32(g.Config.ScreenWidth) / 2
			centerY := float32(g.Config.ScreenHeight) / 2

			camOffsetX := Map.Cam.Offset[0] * unit * zoom
			camOffsetY := Map.Cam.Offset[1] * unit * zoom

			elemOffsetX := Box[2] * unit * zoom
			elemOffsetY := Box[3] * unit * zoom

			// === INVERSION DES AXES X ET Y ===
			// Votre système : X+ = gauche, Y+ = haut
			// Ebiten : X+ = droite, Y+ = bas
			posX := centerX - (worldX - camOffsetX + elemOffsetX) // X INVERSIÉ
			posY := centerY - (worldY - camOffsetY + elemOffsetY) // Y INVERSIÉ

			posX -= unit
			posY -= unit

			if img == nil {
				if g.Debug {
					drawRect(screen, posX, posY, spriteWidth, spriteHeight, color.RGBA{255, 0, 0, 255})
				}
				continue
			}

			if g.Debug {
				if spriteWidth == 0 {
					drawCircle(screen, posX, posY, spriteHeight, ETEHelper.ImgMoyenne(*img))
				} else {
					drawRect(screen, posX, posY, spriteWidth, spriteHeight, ETEHelper.ImgMoyenne(*img))
				}
				drawText(screen, posX, posY, ETEHelper.GetKey(g.Sprites, img))
			} else {
				opts := &ebiten.DrawImageOptions{}

				opts.GeoM.Translate(-float64(spriteWidth)/2, -float64(spriteHeight)/2)

				imgWidth := float64(img.Bounds().Dx())
				imgHeight := float64(img.Bounds().Dy())
				if imgWidth > 0 && imgHeight > 0 {
					opts.GeoM.Scale(float64(spriteWidth)/imgWidth, float64(spriteHeight)/imgHeight)
				}

				// === ROTATION 180° POUR CORRIGER LES SPRITES À L'ENVERS ===
				rotation := float64(Box[8]) + 3.14159 // + π radians (180°)
				opts.GeoM.Rotate(rotation)

				opts.GeoM.Translate(float64(posX), float64(posY))

				screen.DrawImage(img, opts)
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

func (g *Game) Quit() {

}
