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
					Tiles[x][y].Neighbours = append(Tiles[x][y].Neighbours, &Tiles[n[0]][n[1]])
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
				//check for lightnings first
				if rand.Intn(100) <= config.Lightnings {
					//if rand.Intn(100) <= 30 {
					x := rand.Intn(w)
					y := rand.Intn(h)

					t := &g.Tiles[x][y]
					if t.Status == tree {
						t.Status = hitByLightning
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
									prob = 3500
									//logger.Println(strconv.Itoa(x) + " Nachbarn")
								case 3, 4:
									//prob = 8000
									prob = 2800
									//logger.Println(strconv.Itoa(x) + " Nachbarn")
								case 2:
									//prob = 4000
									prob = 1400
									//logger.Println(strconv.Itoa(x) + " Nachbarn")
								case 1:
									prob = 700
								}
								if rand.Intn(100000) <= calcMapProb(prob) {
									// 1i%% of cells trees start growing
									//img := SubImage{t.X, t.Y, imgTree}
									img := gd.Img.Tree
									t.SubImages[0].Image = img
									t.Status = tree
								}
							case tree:
								var firecounter int
								for _, n := range t.Neighbours {
									if n.Status == fireFull {
										//	if n.fireDuration < config.Fireduration-5 {
										firecounter += 3
										//	}
									}
									if rand.Intn(100) < firecounter {

										t.Status = fireSmall
									}
								}

							case hitByLightning:
								//img := SubImage{t.X, t.Y, imgFire}
								t.SubImages[fireSmall].Image = gd.Img.FireSmall[0]
								t.Status = fireSmall
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
								}

							}
						}
					}
				}
			}
		} else {
			time.Sleep(time.Second * 1)
		}

	}
	logger.Println("Stop")
	//logger.ende
}

/*
	else {
		for i, c1 := range cells {
			for j := range c1 {
				c := &cells[i][j]
				var tmpcell = cell{c.x, c.y, nil, rune(' '), 5, 5, false, rune(' '), c.style.Normal(), c.style.Normal()}
				//c := &cells[id]

				switch c.Rune {

				case r.wasteland:

					switch {
					case c.wastelandDuration == config.WastelandDuration-50:
						tmpcell.Rune = r.wasteland
						tmpcell.style = tmpcell.style.Foreground(tcell.ColorBlack)
						toChange = append(toChange, tmpcell)
						c.wastelandDuration -= 1
					case c.wastelandDuration > 0:

						c.wastelandDuration -= 1
					default:
						if rand.Intn(100) < 2 {
							tmpcell.Rune = r.empty
							toChange = append(toChange, tmpcell)
						}
					}

				case r.burningTree:
					if c.fireDuration > 0 {
						c.fireDuration -= 1
					} else {
						// nothing grows for wastlandDuration
						tmpcell.Rune = r.wasteland
						tmpcell.wastelandDuration = config.WastelandDuration
						tmpcell.style = tmpcell.style.Foreground(tcell.ColorFireBrick)
						toChange = append(toChange, tmpcell)
					}
				case r.lightning:
					tmpcell.Rune = c.lastRune
					tmpcell.style = c.laststyle

					for _, n := range c.Neighbours {

						n.style = n.laststyle
						toChange = append(toChange, *n)
					}
					if c.lastRune == r.tree {
						if rand.Intn(100) < config.LightningStartsFire {

							tmpcell.Rune = r.burningTree
							tmpcell.style = tmpcell.style.Foreground(tcell.ColorRed)
							tmpcell.fireDuration = config.Fireduration
						}
					}

					toChange = append(toChange, tmpcell)

				case r.tree:
					// if a Neighbour burns, I may burn
					var firecounter int
					for _, n := range c.Neighbours {
						if n.Rune == r.burningTree {
							if n.fireDuration < config.Fireduration-5 {
								firecounter += 3
							}
						}
						if rand.Intn(100) < firecounter {
							tmpcell.Rune = r.burningTree
							tmpcell.fireDuration = config.Fireduration
							tmpcell.style = tmpcell.style.Foreground(tcell.ColorRed)
							toChange = append(toChange, tmpcell)

						}
					}

				case r.empty:

					prob := config.CreateNewTree
					//prob := 400
					var x int
					for _, n := range c.Neighbours {
						if n.Rune == r.tree {
							x++
							//prob = prob + 3000
						}
					}
					switch x {
					case 5, 6, 7, 8:
						prob = 500000
						//logger.Println(strconv.Itoa(x) + " Nachbarn")
					case 3, 4:
						prob = 12000
						//logger.Println(strconv.Itoa(x) + " Nachbarn")
					case 2:
						prob = 5000
						//logger.Println(strconv.Itoa(x) + " Nachbarn")
					case 1:
						prob = 1000
					}
					if rand.Intn(1000000) <= prob {
						// 1i%% of cells trees start growing
						tmpcell.Rune = r.tree
						tmpcell.style = tmpcell.style.Foreground(tcell.ColorSpringGreen)
						toChange = append(toChange, tmpcell)
					}
				}
			}
			//			logger.Printf("%v %p \n", &cells[1][2].Neighbours, cells[1][2].Neighbours)
		}
	}
	for _, c := range toChange {
		switch c.Rune {

		case r.wasteland:
			cells[c.x][c.y].Rune = r.wasteland
			cells[c.x][c.y].wastelandDuration = c.wastelandDuration
			cells[c.x][c.y].lastRune = c.lastRune
			cells[c.x][c.y].style = c.style
			cells[c.x][c.y].laststyle = c.laststyle
			scn.SetContent(c.x, c.y, c.Rune, []rune(""), c.style)
			//scn.SetContent(c.x, c.y, c.Rune, []rune(""), tcell.StyleDefault.Foreground(tcell.ColorBlack))
		case r.burningTree:
			cells[c.x][c.y].Rune = r.burningTree
			cells[c.x][c.y].fireDuration = c.fireDuration
			cells[c.x][c.y].style = c.style
			cells[c.x][c.y].laststyle = c.laststyle
			//scn.SetContent(c.x, c.y, c.Rune, []rune(""), tcell.StyleDefault.Foreground(tcell.ColorRed))
			scn.SetContent(c.x, c.y, c.Rune, []rune(""), c.style)
		case r.empty:
			cells[c.x][c.y].Rune = r.empty
			scn.SetContent(c.x, c.y, c.Rune, []rune(""), tcell.StyleDefault)
		case r.lightning:
			cells[c.x][c.y].Rune = r.lightning
			cells[c.x][c.y].lastRune = c.lastRune
			cells[c.x][c.y].laststyle = c.laststyle
			//scn.SetContent(c.x, c.y, c.Rune, []rune(""), tcell.StyleDefault.Foreground(tcell.ColorYellow))
			scn.SetContent(c.x, c.y, c.Rune, []rune(""), c.style)
		case r.tree:
			cells[c.x][c.y].Rune = r.tree
			cells[c.x][c.y].style = c.style
			cells[c.x][c.y].laststyle = c.laststyle
			//scn.SetContent(c.x, c.y, c.Rune, []rune(""), tcell.StyleDefault.Foreground(tcell.ColorSpringGreen))
			scn.SetContent(c.x, c.y, c.Rune, []rune(""), c.style)
		}

	}
	scn.Show()

	//time.Sleep(time.Millisecond * 50)
*/
