package main

import (
	//  "farni.com/assets"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"image"
	//	"image/color"

	"image/png"
	"math/rand"
	"os"
	"time"
	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//var imgTree *ebiten.Image
//var imgFire *ebiten.Image

func Printme() {
	fmt.Println("Printme")
}

type Map struct {
	Image      *ebiten.Image
	SubImages  [][]*ebiten.Image
	Width      int
	Height     int
	TileWidth  int
	TileHeight int
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func (m *Map) Load(filename string) image.Image {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ebitImage := ebiten.NewImageFromImage(img)
	m.Image = ebitImage
	bounds := img.Bounds()
	m.Width, m.Height = bounds.Max.X, bounds.Max.Y
	return img

}
/*
func (m *Map) MakeMapSubImages(img image.Image) {
	log.Println(m.Width, m.TileWidth, m.TileHeight)
	si, ok := (img).(SubImager)
	if !ok {
		fmt.Println(": img does not support SubImage()")
	}
	for x := 0; x < m.Width; x = x + m.TileWidth {
		var sub = make([]*ebiten.Image, 0) //m.Height/TileHeight)
		for y := 0; y < m.Height; y = y + m.TileHeight {

			pointX := image.Point{x, y}
			pointY := image.Point{x + m.TileWidth, y + m.TileHeight}
			subImg := si.SubImage(image.Rectangle{pointX, pointY})
			i := ebiten.NewImageFromImage(subImg)
			sub = append(sub, i)
		}
		m.SubImages = append(m.SubImages, sub)
	}
}
*/


func (m *Map) MakeTileProperties(g **Game, img image.Image) {
	//tiles := g.Tiles
	rand.Seed(time.Now().Unix())
	var xx Game
	xx = **g
	//fmt.Println(xx)
	size := img.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {

			//pixel := img.At(x, y)
			r, g, b, _ := img.At(x, y).RGBA()
			//log.Println(r, g, b, a)
			//c := color.RGBAModel.Convert(pixel).(color.RGBA)
			if b > g && b > r {
				xx.Tiles[int(x/gd.TileWidth)][int(y/gd.TileHeight)].Properties["isWater"] = true
				rand.Intn(gd.TileWidth)

			} else if g > r && g > b {
				xx.Tiles[int(x/gd.TileWidth)][int(y/gd.TileHeight)].Properties["isForest"] = true
				xx.Tiles[int(x/gd.TileWidth)][int(y/gd.TileHeight)].OffsetX = rand.Intn(gd.TileWidth)
				xx.Tiles[int(x/gd.TileWidth)][int(y/gd.TileHeight)].OffsetY = rand.Intn(gd.TileHeight)
			}
		}
	}
}

/*
	si, _ := (img).(SubImager)
	var tp = make([][]map[string]bool, 0)

	for x := 0; x < m.Width; x = x + m.TileWidth {
		var pp = make([]map[string]bool, 0)
		for y := 0; y < m.Height; y = y + m.TileHeight {

			pointX := image.Point{x, y}
			pointY := image.Point{x + m.TileWidth, y + m.TileHeight}
			//fmt.Println(pointX,pointY)
			subImg := si.SubImage(image.Rectangle{pointX, pointY})
			p := map[string]bool{}
			for xx := 0; xx < m.TileWidth; xx++ {
				for yy := 0; yy < m.TileHeight; yy++ {
					pixel := subImg.At(xx, yy)
					c := color.RGBAModel.Convert(pixel).(color.RGBA)
					if c.G > 0 && c.G < 255 {
						//	if c.G > 0 {
						fmt.Println(xx, yy)
						fmt.Println(c)
						p["isForest"] = true

					} else {
						p["isWater"] = true
						//fmt.Println(pixel)
					}
				}
			}
			//prop := GetTileProperties(subImg)
			pp = append(pp, p)
		}
		tp = append(tp, pp)
	}

	return tp
*/
/*
func GetTileProperties(img image.Image) map[string]bool {

	size := img.Bounds().Size()
	//rect := image.Rect(0, 0, size.X, size.Y)
	//	wImg := image.NewRGBA(rect)

	p := map[string]bool{}
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			c := color.RGBAModel.Convert(pixel).(color.RGBA)
			if c.G > 0 && c.G < 255 {
				//	if c.G > 0 {
				fmt.Println(y, x)
				fmt.Println(c)
				p["isForest"] = true

			} else {
				p["isWater"] = true
				//fmt.Println(pixel)
			}
		}
	}
	return p

}
*/
