package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	dst := image.NewRGBA(image.Rect(0, 0, 900, 420))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.RGBA{R: 30, G: 30, B: 30, A: 255}}, image.Point{}, draw.Src)

	panel := image.Rect(40, 40, 860, 380)
	draw.Draw(dst, panel, &image.Uniform{C: color.RGBA{R: 19, G: 52, B: 87, A: 255}}, image.Point{}, draw.Src)

	addText(dst, 60, 80, "POST https://localhost:8080/auth/login", color.White)
	addText(dst, 60, 120, "Body:", color.RGBA{R: 200, G: 220, B: 255, A: 255})
	addText(dst, 100, 150, `{`, color.White)
	addText(dst, 120, 180, `"email": "user@example.com",`, color.White)
	addText(dst, 120, 210, `"password": "secret123"`, color.White)
	addText(dst, 100, 240, `}`, color.White)

	addText(dst, 60, 280, "Response 200 OK:", color.RGBA{R: 200, G: 220, B: 255, A: 255})
	addText(dst, 100, 310, `{ \"token\": \"<jwt>\" }`, color.White)

	if err := os.MkdirAll("docs", 0o755); err != nil {
		log.Fatalf("create docs dir: %v", err)
	}

	file, err := os.Create("docs/postman_example.png")
	if err != nil {
		log.Fatalf("create file: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, dst); err != nil {
		log.Fatalf("encode png: %v", err)
	}

	log.Println("Generated docs/postman_example.png")
}

func addText(img draw.Image, x, y int, label string, col color.Color) {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
