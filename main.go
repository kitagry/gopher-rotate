package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	rmascot "github.com/hajimehoshi/ebiten/v2/examples/resources/images/mascot"
)

const (
	width  = 200
	height = 200
)

var (
	gopher1 *ebiten.Image
	gopher2 *ebiten.Image
	gopher3 *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img1, _, err := image.Decode(bytes.NewReader(rmascot.Out01_png))
	if err != nil {
		log.Fatal(err)
	}
	gopher1 = ebiten.NewImageFromImage(img1)

	img2, _, err := image.Decode(bytes.NewReader(rmascot.Out02_png))
	if err != nil {
		log.Fatal(err)
	}
	gopher2 = ebiten.NewImageFromImage(img2)

	img3, _, err := image.Decode(bytes.NewReader(rmascot.Out03_png))
	if err != nil {
		log.Fatal(err)
	}
	gopher3 = ebiten.NewImageFromImage(img3)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type mascot struct {
	x16  int
	y16  int
	vx16 int
	vy16 int

	ground int
	count  int
}

func (m *mascot) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (m *mascot) Update() error {
	m.count++

	sw, sh := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowPosition(m.x16/16, m.y16/16+sh-height)

	if m.vx16 == 0 && m.x16 == 0 {
		m.vx16 = 64
	}
	m.x16 += m.vx16
	// if m.x16/16 > sw-width && m.vx16 > 0 {
	// 	m.vx16 = -64
	// }
	// if m.x16 <= 0 && m.vx16 < 0 {
	// 	m.vx16 = 64
	// }

	// Accelarate the mascot in the Y direction.
	// m.vy16 += 8
	m.y16 += m.vy16

	// If the mascot is on the ground, stop it in the Y direction.
	// if m.y16 >= 0 {
	// 	m.y16 = 0
	// 	m.vy16 = 0
	// }

	switch m.ground % 4 {
	case 0:
		if m.x16/16 > sw-width && m.vx16 > 0 {
			m.ground++
			m.vx16, m.vy16 = m.vy16, -m.vx16
		}
	case 1:
		if -m.y16/16 > sh-height && m.vy16 < 0 {
			m.ground++
			m.vx16, m.vy16 = m.vy16, m.vx16
		}
	case 2:
		if m.x16/16 < 0 && m.vx16 < 0 {
			m.ground++
			m.vx16, m.vy16 = m.vy16, -m.vx16
		}
	case 3:
		if m.y16/16 > 0 && m.vy16 > 0 {
			m.ground++
			m.vx16, m.vy16 = m.vy16, m.vx16
		}
	}

	// If the mascto is on the ground, cause an action in random.
	// if rand.Intn(60) == 0 && m.y16 == 0 {
	// 	switch rand.Intn(2) {
	// 	case 0:
	// 		// Jump.
	// 		m.vy16 = -240
	// 	case 1:
	// 		// Turn.
	// 		m.vx16 = -m.vx16
	// 	}
	// }
	return nil
}

func (m *mascot) Draw(screen *ebiten.Image) {
	img := gopher1
	switch (m.count / 3) % 4 {
	case 0:
		img = gopher1
	case 1, 3:
		img = gopher2
	case 2:
		img = gopher3
	}
	op := &ebiten.DrawImageOptions{}
	w, h := img.Size()
	theta := -float64(m.ground) * 90 * math.Pi * 2 / 360
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2) // move to center
	op.GeoM.Rotate(theta)
	tx := math.Abs(float64(w)/2*math.Sin(theta) + float64(h)/2*math.Cos(theta))
	ty := math.Abs(float64(w)/2*math.Cos(theta) + float64(h)/2*math.Sin(theta))
	op.GeoM.Translate(tx, ty)
	// op.GeoM.Translate(float64(h)/2, float64(w)/2)
	if m.vx16 < 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(width, 0)
	}
	screen.DrawImage(img, op)
}

func main() {
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(width, height)
	if err := ebiten.RunGame(&mascot{}); err != nil {
		log.Fatal(err)
	}
}
