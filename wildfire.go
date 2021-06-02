package main

import (
	//	"farni.com/assets"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//	"sync"
)

//var MapFileName string
var gd GameData

type GameData struct {
	ScreenWidth  int
	ScreenHeight int
	TileWidth    int
	TileHeight   int
	MapFileName  string
	TreeFileName string
	FireFileName string
	Img          GameImages
}

// map.png, tree16x16.png
func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 1024 / 2
	g.ScreenHeight = 768 / 2
	g.TileWidth = 2
	g.TileHeight = 2
	g.MapFileName = "assets/map1024x768.png"
	g.TreeFileName = "assets/tree4x4.png"
	g.FireFileName = "assets/fire/fire.png"
}

/*
// map.png, tree8x8.png
func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 1024 / 8
	g.ScreenHeight = 768 / 8
	g.TileWidth = 8
	g.TileHeight = 8
	g.MapFileName = "assets/map1024x768.png"
	g.TreeFileName = "assets/tree8x8.png"
	g.FireFileName = "assets/fire/fire.png"
}
// map.png, tree16x16.png
func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 1024 / 4
	g.ScreenHeight = 768 / 4
	g.TileWidth = 4
	g.TileHeight = 4
	g.MapFileName = "assets/map1024x768.png"
	g.TreeFileName = "assets/tree16x16.png"
	g.FireFileName = "assets/fire/fire.png"
}
// map.png, tree4x4.png
func (g *GameData) init() {
	//g := GameData{
	g.ScreenWidth = 1024 / 4
	g.ScreenHeight = 768 / 4
	g.TileWidth = 4
	g.TileHeight = 4
	g.MapFileName = "assets/map1024x768.png"
	g.TreeFileName = "assets/tree4x4.png"
	g.FireFileName = "assets/fire/fire4x4.png"
}
// map64x48.png, tree4x4.png
func (g *GameData) init() {
	g.ScreenWidth = 64 / 4
	g.ScreenHeight = 48 / 4
	g.TileWidth = 4
	g.TileHeight = 4
	g.MapFileName = "assets/map64x48.png"
	g.TreeFileName = "assets/tree4x4.png"
	g.FireFileName = "assets/fire/fire4x4.png"
}
// map64*48.png, tree8x8.png
func (g *GameData) init() {
	g.ScreenWidth = 8
	g.ScreenHeight = 6
	g.TileWidth = 8
	g.TileHeight = 8
	g.MapFileName = "assets/map64x48.png"
	g.TreeFileName = "assets/tree8x8.png"
	g.FireFileName = "assets/fire/fire.png"
}
// map64x48.png, tree16x16.png
func (g *GameData) init() {
	g.ScreenWidth = 64 / 8
	g.ScreenHeight = 48 / 8
	g.TileWidth = 2
	g.TileHeight = 2
	g.MapFileName = "assets/map64x48.png"
	g.TreeFileName = "assets/tree4x4.png"
	g.FireFileName = "assets/fire/fire.png"
}
*/

type Point struct {
	X int
	Y int
}

type GameImages struct {
	Tree      *ebiten.Image
	FireSmall []*ebiten.Image
	FireFull  []*ebiten.Image
}

func (g *GameData) LoadSubImages() {
	g.Img = GameImages{}

	g.Img.Tree, _, _ = ebitenutil.NewImageFromFile(g.TreeFileName)

	imgFire, _, _ := ebitenutil.NewImageFromFile(g.FireFileName)
	bounds := imgFire.Bounds()
	/*
		"imgFire contains 4 images, 2 per row"
		"row 0 contains to 2 images representing fireSmall"
		"row 1 contains to 2 images representing fireFull"
	*/
	p := Point{bounds.Max.X, bounds.Max.Y}
	g.Img.FireSmall = append(g.Img.FireSmall, imgFire.SubImage(image.Rect(0, 0, p.X/2-1, p.Y/2-1)).(*ebiten.Image))
	g.Img.FireSmall = append(g.Img.FireSmall, imgFire.SubImage(image.Rect(p.X/2, 0, p.X-1, p.Y/2-1)).(*ebiten.Image))
	g.Img.FireFull = append(g.Img.FireFull, imgFire.SubImage(image.Rect(0, p.Y/2-1, p.X/2, p.Y-1)).(*ebiten.Image))
	g.Img.FireFull = append(g.Img.FireFull, imgFire.SubImage(image.Rect(p.X/2, p.Y/2, p.X-1, p.Y-1)).(*ebiten.Image))

	//return img
	//imgTree, _, _ = ebitenutil.NewImageFromFile(gd.TreeFileName)
}

type Game struct {
	Tiles       [][]Tile
	ActiveTiles map[Point]*Tile
}

type Tile struct {
	X                 int
	Y                 int
	Properties        map[string]bool
	OffsetX           int
	OffsetY           int
	MapImage          *ebiten.Image
	SubImages         []SubImage
	Neighbours        []*Tile
	Status            status
	fireDuration      int
	wastelandDuration int
}
type SubImage struct {
	//X     int
	//Y     int
	Image *ebiten.Image
}

var TileProperties = make([][]map[string]bool, 0)
var fireCounter int
var mutex = &sync.Mutex{}
var m Map

func CreateTiles(g *Game) {
	//allTiles := make([][]Tile, gd.ScreenWidth)
	rand.Seed(time.Now().Unix())

	m.Load(gd.MapFileName)

	m.TileWidth = gd.TileWidth
	m.TileHeight = gd.TileHeight
	//m.MakeMapSubImages(img)
	g.Tiles = make([][]Tile, gd.ScreenWidth)
	for x := range g.Tiles {

		tiles := make([]Tile, gd.ScreenHeight)
		for y := range tiles {
			tile := Tile{
				X:       x * gd.TileWidth,
				Y:       y * gd.TileHeight,
				OffsetX: rand.Intn(gd.TileWidth),
				OffsetY: rand.Intn(gd.TileHeight),
				Status:  empty,
				//Properties: make( Properties[string]bool{},0),
				//Properties: make(map[string]bool),
				SubImages: make([]SubImage, 5),
			}
			tile.Properties = m.GetProperties(tile)

			//	tile.MapImage = m.SubImages[x][y]
			tiles[y] = tile
			//g.Tiles[x][y] = tile
		}

		g.Tiles[x] = tiles
	}

	//m.MakeTileProperties(&g, img)
	//return allTiles
	g.Tiles = getNeighbours(g.Tiles)
}

func leftTouched() bool {
	for _, id := range ebiten.TouchIDs() {
		x, _ := ebiten.TouchPosition(id)
		if x < gd.ScreenWidth/2 {
			return true
		}
	}
	return false
}

/*
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || leftTouched() {
		log.Println("left")
	}

	return nil
}
*/

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	//img1, _, _ := ebitenutil.NewImageFromFile("assets/tree1.png")
	//log.Println(len(g.ActiveTiles))
	screen.DrawImage(ebiten.NewImageFromImage(m.Image), nil)
	var op = &ebiten.DrawImageOptions{}
	mutex.Lock()
	for _, tile := range g.ActiveTiles {

		//for _, tile := range tiles {

		//wg.Wait()
		for _, s := range tile.SubImages {
			if s.Image != nil {

				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.X-tile.OffsetX), float64(tile.Y-tile.OffsetY))
				screen.DrawImage(s.Image, op)
				//s.Image = nil
			}
		}

		//}

	}
	mutex.Unlock()
	if showDebugInfo {
		msg := fmt.Sprintf(`TPS: %0.2f
	FPS: %0.2f
	Num of tiles: %d
	Press Space to pause game`, ebiten.CurrentTPS(), ebiten.CurrentFPS(), len(g.ActiveTiles))
		ebitenutil.DebugPrint(screen, msg)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gd.TileWidth * gd.ScreenWidth, gd.TileHeight * gd.ScreenHeight
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) || leftTouched() {
		pause = !pause
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) || leftTouched() {
		showDebugInfo = !showDebugInfo
	}

	return nil
}
func main() {
	Printme()
	gd.init()
	gd.LoadSubImages()
	g := &Game{}
	CreateTiles(g)
	g.ActiveTiles = make(map[Point]*Tile)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(400, 300)
	//ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Wildfire")
	go updateGame(g)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
