// package main

// import (
// 	"fmt"

// 	"github.com/yeqown/go-qrcode/v2"
// 	"github.com/yeqown/go-qrcode/writer/standard"
// )

// func main() {
// 	qrcode.WithEncodingMode(qrcode.EncModeAlphanumeric)
// 	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
// 	if err != nil {
// 		fmt.Printf("could not generate QRCode: %v", err)
// 		return
// 	}

// 	w, err := standard.New("repo-qrcode.jpeg")
// 	if err != nil {
// 		fmt.Printf("standard.New failed: %v", err)
// 		return
// 	}

// 	// save file
// 	if err = qrc.Save(w); err != nil {
// 		fmt.Printf("could not save image: %v", err)
// 	}
// }

// package main

// import (
// 	"github.com/yeqown/go-qrcode/v2"
// 	"github.com/yeqown/go-qrcode/writer/terminal"
// )

// func main() {
// 	qrc, _ := qrcode.New("fuck you")

// 	w := terminal.New()

// 	if err := qrc.Save(w); err != nil {
// 		panic(err)
// 	}
// }

// package main

// import (
// 	"github.com/yeqown/go-qrcode/v2"
// 	"github.com/yeqown/go-qrcode/writer/standard"
// )

// type smallerCircle struct {
// 	smallerPercent float64
// }

// func (sc *smallerCircle) DrawFinder(ctx *standard.DrawContext) {
// 	backup := sc.smallerPercent
// 	sc.smallerPercent = 1.0
// 	sc.Draw(ctx)
// 	sc.smallerPercent = backup
// }

// func newShape(radiusPercent float64) standard.IShape {
// 	return &smallerCircle{smallerPercent: radiusPercent}
// }

// func (sc *smallerCircle) Draw(ctx *standard.DrawContext) {
// 	w, h := ctx.Edge()
// 	x, y := ctx.UpperLeft()
// 	color := ctx.Color()

// 	// choose a proper radius values
// 	radius := w / 2
// 	r2 := h / 2
// 	if r2 <= radius {
// 		radius = r2
// 	}

// 	// 80 percent smaller
// 	radius = int(float64(radius) * sc.smallerPercent)

// 	cx, cy := x+float64(w)/2.0, y+float64(h)/2.0 // get center point
// 	ctx.DrawCircle(cx, cy, float64(radius))
// 	ctx.SetColor(color)
// 	ctx.Fill()

// }

// func main() {
// 	shape := newShape(0.7)
// 	qrc, err := qrcode.New("with-custom-shape")
// 	// qrc, err := qrcode.New("with-custom-shape", qrcode.WithCircleShape())
// 	if err != nil {
// 		panic(err)
// 	}

// 	w, err := standard.New("./smaller.png", standard.WithCustomShape(shape))
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = qrc.Save(w)
// 	if err != nil {
// 		panic(err)
// 	}
// }

package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/yeqown/go-qrcode/writer/standard"

	"github.com/yeqown/go-qrcode/v2"
)

const resultDir = "./qrCodes"

func main() {
	link := flag.String("linkedin", "https://linkedin.com/", "Linkedin profile link")
	logo := flag.String("logo", "https://github.com/PPerminov/qrLogo/raw/main/linkeindLogo.png", "logo to use")
	name := flag.String("name", "linkedinQRCode", "name for the file")
	flag.Parse()

	qrc, err := qrcode.NewWith(*link,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(resultDir, os.ModeDir)
	if err != nil {
		panic(err)
	}
	w, err := standard.New(
		fmt.Sprintf("%s/%s.png", resultDir, *name),
		standard.WithHalftone(downloadAndSaveOrReturnPath(*logo)),
		standard.WithBgColor(color.RGBA{G: 0, B: 0, R: 0}),
		standard.WithFgColor(color.RGBA{G: 255, B: 255, R: 255}),
	)
	if err != nil {
		panic(err)
	}

	if err = qrc.Save(w); err != nil {
		panic(err)
	}
}

func downloadAndSaveOrReturnPath(path string) string {
	u, err := url.Parse(path)
	if err != nil || u.Host == "" {
		return path
	}
	data, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}
	defer data.Body.Close()
	file, err := os.CreateTemp("", "logoImage*.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, data.Body)
	if err != nil {
		panic(err)
	}
	return file.Name()
}
