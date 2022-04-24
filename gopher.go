package gopherrotate

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	rmascot "github.com/hajimehoshi/ebiten/v2/examples/resources/images/mascot"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	width  = 240
	height = 280

	imgWidth  = 200
	imgHeight = 200
)

var (
	gopher1 *ebiten.Image
	gopher2 *ebiten.Image
	gopher3 *ebiten.Image

	fukidashi *ebiten.Image

	mplusNormalFont font.Face
)

//go:embed assets/fukidashi.png
var fukidashiPng []byte

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

	img4, _, err := image.Decode(bytes.NewReader(fukidashiPng))
	if err != nil {
		log.Fatal(err)
	}
	fukidashi = ebiten.NewImageFromImage(img4)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type horizonMascot struct {
	x16  int
	y16  int
	vx16 int
	vy16 int
}

func (h *horizonMascot) Update() error {
	if h.vx16 == 0 && h.x16 == 0 {
		h.vx16 = 64
	}
	h.x16 += h.vx16

	h.y16 += h.vy16
	h.vy16 -= 8

	if h.y16 < 0 {
		h.y16 = 0
	}

	// If the mascto is on the ground, cause an action in random.
	if rand.Intn(60) == 0 {
		switch rand.Intn(2) {
		case 0:
			// Jump.
			if h.y16 == 0 {
				// h.vy16 = 240
			}
		case 1:
			// Turn.
			h.vx16 = -h.vx16
		}
	}
	return nil
}

type Mascot struct {
	h   *horizonMascot
	x16 int
	y16 int

	ground  int
	count   int
	reverse bool

	msg string
}

func (m *Mascot) Say(msg string) {
	m.msg = msg
}

func NewMascot() *Mascot {
	return &Mascot{
		h: &horizonMascot{},
	}
}

func (m *Mascot) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (m *Mascot) Update() error {
	m.count++
	err := m.h.Update()
	if err != nil {
		return err
	}

	sw, sh := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowPosition(m.x16/16, m.y16/16+sh-height)

	m.reverse = m.h.vx16 < 0

	ground := 0
	x16 := m.h.x16
	if x16 > 0 {
		for {
			if ground%2 == 0 && x16/16 > sw-width {
				ground++
				x16 -= (sw - width) * 16
			} else if ground%2 == 1 && x16/16 > sh-height {
				ground++
				x16 -= (sh - height) * 16
			} else {
				break
			}
		}
	} else {
		for {
			if ground < 0 {
				ground += 4
			}
			if ground%2 == 0 && x16/16 < 0 {
				ground--
				x16 += (sh - height) * 16
			} else if ground%2 == 1 && x16/16 < 0 {
				ground--
				x16 += (sw - width) * 16
			} else {
				break
			}
		}
	}
	m.ground = ground

	switch m.ground % 4 {
	case 0:
		m.x16 = x16
		m.y16 = m.h.y16
	case 1:
		m.x16 = (sw-width)*16 - m.h.y16
		m.y16 = -x16
	case 2:
		m.x16 = (sw-width)*16 - x16
		m.y16 = -(sh-height)*16 + m.h.y16
	case 3:
		m.x16 = m.h.y16
		m.y16 = -(sh-height)*16 + x16
	}

	return nil
}

func (m *Mascot) Draw(screen *ebiten.Image) {
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
	switch m.ground % 4 {
	case 0:
		op.GeoM.Translate(0, height-imgHeight)
	case 1:
		op.GeoM.Translate(width-imgWidth, 0)
	}

	if m.reverse {
		if m.ground%2 == 0 {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(imgWidth, 0)
		} else if m.ground%2 == 1 {
			op.GeoM.Scale(1, -1)
			op.GeoM.Translate(0, imgHeight)
		}
	}
	screen.DrawImage(img, op)

	if m.msg != "" {
		fukidashiOp := &ebiten.DrawImageOptions{}
		switch m.ground % 4 {
		case 0:
			fukidashiOp.GeoM.Translate(float64(w)*3/5, float64(h))
		case 1:
			fukidashiOp.GeoM.Translate(float64(w)/2, float64(h))
		case 2:
			fukidashiOp.GeoM.Translate(float64(w)/2, float64(h)*5/4)
		case 3:
			fukidashiOp.GeoM.Translate(float64(w)/2, float64(h))
		}
		fukidashiOp.GeoM.Scale(0.8, 0.8)
		screen.DrawImage(fukidashi, fukidashiOp)
		text.Draw(fukidashi, m.msg, mplusNormalFont, 40, height-220, color.Black)
	}
}

func (m *Mascot) Run() error {
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(width, height)
	if err := ebiten.RunGame(m); err != nil {
		return err
	}
	return nil
}
