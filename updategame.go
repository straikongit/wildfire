package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/gdamore/tcell"

	//"strconv"
	//"config"
	"time"
)

type cell struct {
	x                 int
	y                 int
	Neighbours        []*cell
	Rune              rune
	fireDuration      int
	wastelandDuration int
	hitbyLightning    bool
	lastRune          rune
	style             tcell.Style
	laststyle         tcell.Style
}
type status int

const (
	empty status = iota
	tree
	hitByLightning
	fireSmall
	fireFull
	wasteland
)

func getNeighbours(Tiles [][]Tile) (c [][]Tile) {

	for x, c := range Tiles {
		for y := range c {
			var s [][]int
			s = append(s, []int{x - 1, y})
			s = append(s, []int{x - 1, y - 1})
			s = append(s, []int{x - 1, y + 1})

			s = append(s, []int{x + 1, y})
			s = append(s, []int{x + 1, y - 1})
			s = append(s, []int{x + 1, y + 1})

			s = append(s, []int{x, y + 1})
			s = append(s, []int{x, y - 1})
			for _, n := range s {
				if n[0] >= 0 && n[0] < gd.ScreenWidth && n[1] >= 0 && n[1] < gd.ScreenHeight {
					t := &Tiles[n[0]][n[1]]
					if t.Properties["isForest"] && !t.Properties["isWater"] {
						Tiles[x][y].Neighbours = append(Tiles[x][y].Neighbours, &Tiles[n[0]][n[1]])
					}
				}

			}
		}
	}

	return Tiles
}

var Width, Height int
var pause bool

/*
Wahrscheinlichkeiten sind auf Screen 1024*768 ausgelegt.
Kleinere Karten brauchen größere Wahrscheinlichkeiten
Klappt so Mittel
*/
func calcMapProb(prob int) int {
	p := 1024 / gd.ScreenWidth / 2 * prob
	return int(p)

}

func updateGame(g *Game) {
	// logger
	f, err := os.OpenFile("/home/far/daten/dev/go/src/farhome/wildfire/fire.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "fire", log.LstdFlags)
	logger.Println("Start")
	//logger.ende

	var config Config
	config.init()

	quit := make(chan struct{})
	/*
				go func() {
					for {
						ev := scn.PollEvent()
						logger.Println(ev)
						switch event := ev.(type) {

						case *tcell.EventKey:

							switch event.Key() {

							case tcell.KeyEscape:
								close(quit)
							case tcell.KeyCtrlC:
								close(quit)
							case tcell.KeyCtrlSpace:
								pause = !pause
							case 256: // Space
								pause = !pause
							}
						}

					}
				}()
		loop:
	*/
	w := gd.ScreenWidth
	h := gd.ScreenHeight
	for {

		if !pause {
			select {
			case <-quit:
				break //loop
			case <-time.After(time.Millisecond * config.PausePerRound):
								mutex.Lock()
				x := rand.Intn(w)
				y := rand.Intn(h)
				t := g.Tiles[x][y]
				if t.Properties["isForest"] && !t.Properties["isWater"] {
					if t.Status == empty {
						if rand.Intn(100) <= 1 { //calcMapProb(prob) {
							// 1% of cells trees start growing
							img := gd.Img.Tree
							t.SubImages[0].Image = img
							t.Status = tree
							g.ActiveTiles[Point{t.X, t.Y}] = &t
							for _, n := range t.Neighbours {
								if n.Status == empty {
									g.ActiveTiles[Point{t.X, t.Y}] = n
								}

							}
						}
					}
				}
				for x, t1 := range g.Tiles {
					for y := range t1 {
						t := &g.Tiles[x][y]
						//c := &cells[id]
						if t.Properties["isForest"] && !t.Properties["isWater"] {
							switch t.Status {

							case empty:

								prob := config.CreateNewTree
								//prob := 400
								var count int
								for _, n := range t.Neighbours {
									//if n.Status != empty {
									if n.Status == tree {
										count++
										//prob = prob + 3000
									}
								}
								switch count {
								case 5, 6, 7, 8:
									//prob = 300000
									prob = 35000
								case 3, 4:
									//prob = 8000
									prob = 28000
								case 2:
									//prob = 4000
									prob = 14000
								case 1:
									prob = 7000
								}
								if rand.Intn(100000) <= calcMapProb(prob) {
									// 1i%% of cells trees start growing
									img := gd.Img.Tree
									t.SubImages[0].Image = img
									t.Status = tree
									g.ActiveTiles[Point{t.X, t.Y}] = t
								}
							case tree:
								//check for lightnings first
								if rand.Intn(1000000) <= config.Lightnings {
									//if rand.Intn(100) <= 30 {
									/*x := rand.Intn(w)
									y := rand.Intn(h)

									t := &g.Tiles[x][y]

									*/
									t.Status = fireSmall
								} else {

									var firecounter int
									for _, n := range t.Neighbours {
										if n.Status == fireFull {
											//	if n.fireDuration < config.Fireduration-5 {
											firecounter += 2
											//	}
										}
										if rand.Intn(100) < firecounter {

											t.Status = fireSmall
										}
									}
								}

							case fireSmall:
								t.SubImages[fireSmall].Image = gd.Img.FireSmall[rand.Intn(2)]
								rand.Intn(2)
								t.fireDuration++
								if t.fireDuration > config.FireDurationSmall {
									t.Status = fireFull
								}
							case fireFull:

								t.SubImages[fireFull].Image = gd.Img.FireFull[rand.Intn(2)]
								t.fireDuration++
								if t.fireDuration > config.FireDurationFull {
									t.Status = wasteland
								}
							case wasteland:

								if len(t.SubImages) > 0 {
									// `Add(1) signifies that there is 1 task that we need to wait for
									//	wg.Add(1)
									t.SubImages = nil
									t.SubImages = make([]SubImage, 5)
									// Calling `wg.Done` indicates that we are done with the task we are waiting fo
									//	defer wg.Done()
								}
								t.wastelandDuration++
								if t.wastelandDuration > config.WastelandDuration {
									t.wastelandDuration = 0
									t.fireDuration = 0
									t.Status = empty
									delete(g.ActiveTiles, Point{t.X, t.Y})
								}

							}
						}
					}
				}
								mutex.Unlock()
			}
		} else {
			time.Sleep(time.Second * 1)
		}

	}
	logger.Println("Stop")
	//logger.ende
}
