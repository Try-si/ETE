package ETECore

import (
	"errors"
	"image/color"
	"time"

	"github.com/Try-si/ETE/ETEHelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Update() error {
	now := time.Now()
	if !g.LastTime.IsZero() {
		g.DeltaTime = float32(now.Sub(g.LastTime).Seconds())
	}
	g.LastTime = now

	if g.Quite {
		return errors.New("quit")
	}
	return g.UpdateFunc(float32(ebiten.ActualFPS()))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	Map := g.Maps[g.Config.Map]
	unit := float32(Map.Unité)

	parallaxFactor := g.MapConfig.ParrallaxFactor

	for _, layer := range Map.GetSpriteByOrderYZX() {
		height := layer.Height

		if height < Map.Cam.Z {
			continue
		}

		if Map.Cam.DebZ == 0 {
			Map.Cam.DebZ = Map.Cam.Z
		}
		if Map.Cam.Zoom == 0 {
			Map.Cam.Zoom = 1
		}

		dist := (height - Map.Cam.Z) / Map.Cam.Zoom

		if !g.MapConfig.Parrallax {
			dist = Map.Cam.Z / Map.Cam.Zoom
		}

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
				spriteWidth := float32(-Box[0]) * unit / dist
				spriteHeight := float32(-Box[1]) * unit / dist

				// === CALCUL DE LA POSITION MONDE ===
				worldX := -Box[4] * unit
				worldY := -Box[5] * unit

				centerX := float32(g.Config.ScreenWidth) / 2
				centerY := float32(g.Config.ScreenHeight) / 2

				camOffsetX := -Map.Cam.Offset[0] * unit
				camOffsetY := Map.Cam.Offset[1] * unit

				elemOffsetX := Box[2] * unit
				elemOffsetY := Box[3] * unit

				if g.MapConfig.Parrallax || entry.Paralax {
					elemOffsetX = elemOffsetX * parallaxFactor
					elemOffsetY = elemOffsetY * parallaxFactor
				}

				OffsetX := camOffsetX + elemOffsetX
				OffsetY := camOffsetY + elemOffsetY

				// === POSITION FINALE AVEC INVERSION DES AXES ===
				var posX, posY float32
				posX = centerX - (worldX+OffsetX+centerX)/dist + centerX/dist
				posY = centerY - (worldY+OffsetY+centerY)/dist + centerY/dist

				if img == nil {
					if g.Debug {
						drawRect(screen, posX, posY, spriteWidth, spriteHeight, color.RGBA{255, 6, 181, 255})
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

					imgWidth := float64(img.Bounds().Dx())
					imgHeight := float64(img.Bounds().Dy())

					opts.GeoM.Translate(-imgWidth/2, -imgHeight/2)

					if imgWidth > 0 && imgHeight > 0 {
						opts.GeoM.Scale(float64(spriteWidth)/imgWidth, float64(spriteHeight)/imgHeight)
					}

					rotation := float64(Box[8]) + 3.14159
					opts.GeoM.Rotate(rotation)

					opts.GeoM.Translate(float64(posX), float64(posY))

					screen.DrawImage(img, opts)
				}
				//vector.StrokeRect(screen, posX, posY, spriteWidth, spriteHeight, 1, color.Black, false)
				//screen.Set(int(posX), int(posY), color.White)
				//ebitenutil.DebugPrintAt(screen, strconv.FormatFloat(float64(height), 'f', -1, 64), int(posX), int(posY))
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
