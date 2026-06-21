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

	for _, layer := range Map.GetSpriteByOrderYZX() {
		height := layer.Height

		if float32(height) < Map.Cam.Z {
			continue
		}

		dist := Map.Cam.Z - float32(height)

		if dist == 0 {
			dist = 0.0001
		}
		if dist < 0 {
			dist = -dist
		}

		for _, entry := range layer.Sprites {
			if entry.Visible {
				Box := entry.Box
				img := entry.Img

				// === CALCUL DE LA TAILLE (plus loin = plus petit) ===
				spriteWidth := float32(Box[0]) * unit / dist
				spriteHeight := float32(Box[1]) * unit / dist

				// === CALCUL DE LA POSITION MONDE ===
				worldX := Box[4] * unit
				worldY := Box[5] * unit

				centerX := float32(g.Config.ScreenWidth) / 2
				centerY := float32(g.Config.ScreenHeight) / 2

				camOffsetX := Map.Cam.Offset[0] * unit
				camOffsetY := Map.Cam.Offset[1] * unit

				elemOffsetX := Box[2] * unit
				elemOffsetY := Box[3] * unit

				// === POSITION FINALE AVEC PARALLAX ET INVERSION DES AXES ===
				posX := centerX - (worldX+camOffsetX+elemOffsetX+centerX)/dist + centerX/dist
				posY := centerY - (worldY+camOffsetY+elemOffsetY+centerY)/dist + centerY/dist

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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.Config.AdaptativeSize {
		return outsideWidth, outsideHeight
	}
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
