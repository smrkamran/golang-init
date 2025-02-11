package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	dirUp           = Point{x: 0, y: -1}
	dirDown         = Point{x: 0, y: 1}
	dirLeft         = Point{x: -1, y: 0}
	dirRight        = Point{x: 1, y: 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	gameSpeed    = time.Second / 6
	screenWidth  = 800
	screenHeight = 600
	gridSize     = 20
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.lastUpdate = time.Now()

		g.snake = []Point{{
			x: screenWidth / gridSize / 2,
			y: screenHeight / gridSize / 2,
		}}
		g.direction = Point{x: 1, y: 0}
		g.gameOver = false
	}

	if g.gameOver {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.direction = dirUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.direction = dirDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.direction = dirLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.direction = dirRight
	}

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()
	g.UpdateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) UpdateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]

	newHead := Point{head.x + direction.x, head.y + direction.y}

	if g.isCollision(newHead, *snake) {
		g.gameOver = true
		return
	}

	if newHead == g.food {
		*snake = append(
			[]Point{newHead}, *snake...)
		g.spawnFood()
	} else {
		*snake = append(
			[]Point{newHead}, (*snake)[:len(*snake)-1]...)
	}

}

func (g *Game) spawnFood() {
	g.food = Point{
		rand.Intn(screenWidth / gridSize),
		rand.Intn(screenHeight / gridSize),
	}
}

func (g Game) isCollision(newHead Point, snake []Point) bool {
	if newHead.x < 0 || newHead.y < 0 || newHead.x >= screenWidth/gridSize || newHead.y >= screenHeight/gridSize {
		return true
	}
	for _, p := range snake {
		if p == newHead {
			return true
		}
	}

	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake {
		vector.DrawFilledRect(screen, float32(p.x*gridSize), float32(p.y*gridSize), gridSize, gridSize, color.White, true)
	}

	vector.DrawFilledRect(screen, float32(g.food.x*gridSize), float32(g.food.y*gridSize), gridSize, gridSize, color.RGBA{255, 0, 0, 255}, true)

	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		t := "Game Over!"

		w, h := text.Measure(t, face, face.Size)

		op := &text.DrawOptions{}
		op.GeoM.Translate(screenWidth/2-w/2, screenHeight/2-h/2)
		op.ColorScale.ScaleWithColor(color.White)

		text.Draw(screen, t, face, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func StartGame() {
	s, err := text.NewGoTextFaceSource(
		bytes.NewReader(fonts.MPlus1pRegular_ttf),
	)

	if err != nil {
		log.Fatal(err)
	}

	mplusFaceSource = s

	g := &Game{
		snake: []Point{{
			x: screenWidth / gridSize / 2,
			y: screenHeight / gridSize / 2,
		}},
		direction: Point{x: 1, y: 0},
	}

	g.spawnFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func main() {
	StartGame()
}
