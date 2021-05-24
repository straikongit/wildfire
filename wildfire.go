package main

import (
	//	"farni.com/assets"
	//	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//	"image"
)

type Game struct {
	Tiles [][]Tile
}

var MapFileName string
var gd GameData

type GameData struct {
	ScreenWidth  int
	ScreenHeight int
	TileWidth    int
	TileHeight   int
}

func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 4 * 64
	g.ScreenHeight = 4 * 48
	g.TileWidth = 4
	g.TileHeight = 4
	//return g
}

/*
func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 64
	g.ScreenHeight = 48
	g.TileWidth = 16
	g.TileHeight = 16
	//return g
}
func (g *GameData) init() {
	g.ScreenWidth = 2
	g.ScreenHeight = 2
	g.TileWidth = 2
	g.TileHeight = 2
	//return g
}
*/

type Tile struct {
	PixelX     int
	PixelY     int
	Properties map[string]bool
	Image      *ebiten.Image
}

/*
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
*/
var TileProperties = make([][]map[string]bool, 0)

func CreateTiles(g *Game) {
	//allTiles := make([][]Tile, gd.ScreenWidth)
	g.Tiles = make([][]Tile, gd.ScreenWidth)

	var m Map
	img := m.Load(MapFileName)

	m.TileWidth = gd.TileWidth
	m.TileHeight = gd.TileHeight
	m.MakeSubImages(img)
	/*
		f, err := os.Open("assets/map.png")

		if err != nil {
			log.Fatal(err)
		}

		img, err := png.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	*/
	for x := range g.Tiles {

		tiles := make([]Tile, gd.ScreenHeight)
		for y := range tiles {
			/*
				si, ok := (img).(SubImager)
				if !ok {
					fmt.Println(": img does not support SubImage()")
					//			log.Fatal(err)
				}
				pointX := image.Point{x * gd.TileWidth, y * gd.TileHeight}
				pointY := image.Point{x*gd.TileWidth + gd.TileWidth, y*gd.TileHeight + gd.TileHeight}
				subImg := si.SubImage(image.Rectangle{pointX, pointY})
			*/
			tile := Tile{
				PixelX:     x * gd.TileWidth,
				PixelY:     y * gd.TileHeight,
				Properties: map[string]bool{},
			}
			tile.Image = m.SubImages[x][y]
			tiles[y] = tile
			//g.Tiles[x][y] = tile
		}

		g.Tiles[x] = tiles
	}

	m.MakeTileProperties(&g, img)
	//return allTiles
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).

func (g *Game) Draw(screen *ebiten.Image) {
	//img1, _, _ := ebitenutil.NewImageFromFile("assets/tree1.png")
	img1, _, _ := ebitenutil.NewImageFromFile("assets/treesmall.png")
	for x := range g.Tiles {

		for y := range g.Tiles[x] {

			tile := g.Tiles[x][y]
			//Pixels can only be viewed in the main loop. So we do it here
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
			if tile.Properties["isForest"] && !tile.Properties["isWater"] {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				screen.DrawImage(img1, op)

			}
			/*
				if x == 0 { // && x < 11 && y == 20 {
					op = &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
					screen.DrawImage(img1, op)
				}
				if y == 0 {
					op = &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
					screen.DrawImage(img2, op)
				}
			*/
		}

	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gd.TileWidth * gd.ScreenWidth, gd.TileHeight * gd.ScreenHeight
}

func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}
func main() {
	MapFileName = "assets/map.png"
	//MapFileName = "assets/treesmall.png"
	Printme()
	gd.init()

	g := &Game{}
	CreateTiles(g)

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(100, 100)
	//ebiten.SetWindowSize(1024, 768)
	//ebiten.SetWindowSize(4, 4)
	ebiten.SetWindowTitle("Wildfire")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
