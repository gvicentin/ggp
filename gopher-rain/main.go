package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenTitle  = "Gopher Rain"
	screenWidth  = 640
	screenHeight = 480

	// Ground
	groundHeight = 20
	groundY      = screenHeight - groundHeight

	// Gopher
	gopherPixelsWidth  = 14
	gopherPixelsHeight = 14

	gopherScale  = 3.5
	gopherWidth  = gopherPixelsWidth * gopherScale
	gopherHeight = gopherPixelsHeight * gopherScale

	gopherSpeed = 450 // pixels per second

	// Coins
	maxCoins = 10

	coinWidth  = 54
	coinHeight = 54

	coinSpawnY    = -100
	coinSpawnXMax = screenWidth

	coinStartCooldown = 1.5

	livesScale = 2
)

type coin struct {
	active bool
	x, y   float64
}

func (c *coin) spawn() {
	c.x = float64(rand.Intn(coinSpawnXMax - coinWidth))
	c.y = coinSpawnY
}

func testCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x2 < x1+w1 && y1 < y2+h2 && y2 < y1+h1
}

type Game struct {
	// Graphics
	groundImage *ebiten.Image
	gopherImage *ebiten.Image
	coinImage   *ebiten.Image
	livesImage  *ebiten.Image
	font        *text.GoTextFaceSource

	// Gopher
	gopherX, gopherY  float64
	gopherFacingRight bool

	// Coins
	coins        [maxCoins]coin
	coinCooldown float64

	// Game state
	score int
	lives int
}

func (g *Game) Init() error {
	// Ground
	g.groundImage = ebiten.NewImage(screenWidth, groundHeight)
	g.groundImage.Fill(color.RGBA{0x93, 0xb6, 0x5f, 0xff})

	// Gopher
	gopherImage, _, err := ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		return err
	}

	g.gopherImage = gopherImage

	// Coin
	coinImage, _, err := ebitenutil.NewImageFromFile("gopher-coin.png")
	if err != nil {
		return err
	}

	g.coinImage = coinImage

	gopherImageRect := image.Rect(0, 0, gopherPixelsWidth, gopherPixelsHeight)
	g.livesImage = g.gopherImage.SubImage(gopherImageRect).(*ebiten.Image)

	// Font
	file, err := os.Open("pressstart2p.ttf")
	if err != nil {
		return err
	}

	font, err := text.NewGoTextFaceSource(file)
	if err != nil {
		return err
	}

	g.font = font

	g.reset()

	return nil
}

func (g *Game) reset() {
	g.gopherX = screenWidth / 2
	g.gopherY = groundY - gopherHeight

	for i := 0; i < maxCoins; i++ {
		g.coins[i].active = false
		g.coins[i].spawn()
	}

	g.coins[0].active = true
	g.coinCooldown = coinStartCooldown

	g.score = 0
	g.lives = 3
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// Quit the game when escape is pressed
		return ebiten.Termination
	}

	var x float64
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x -= 1
		g.gopherFacingRight = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x += 1
		g.gopherFacingRight = true
	}

	dt := 1.0 / float64(ebiten.TPS())

	// Update gopher position
	g.gopherX += x * gopherSpeed * dt

	if g.gopherX < 0 {
		g.gopherX = 0
	}
	if g.gopherX+gopherWidth > screenWidth {
		g.gopherX = screenWidth - gopherWidth
	}

	// Update coins
	for i := 0; i < maxCoins; i++ {
		if !g.coins[i].active {
			continue
		}

		g.coins[i].y += 100 * dt

		if g.coins[i].y > screenHeight {
			g.coins[i].active = false
			g.lives--
		}
	}

	if g.lives < 0 {
		g.reset()
		return nil
	}

	g.coinCooldown -= dt
	if g.coinCooldown <= 0 {
		spawned := false

		for i := 0; i < maxCoins; i++ {
			if !g.coins[i].active {
				spawned = true
				g.coins[i].spawn()
				g.coins[i].active = true
				g.coinCooldown = coinStartCooldown
				break
			}
		}

		if !spawned {
			panic("No inactive coins available")
		}
	}

	// Collision detection
	for i := 0; i < maxCoins; i++ {
		if !g.coins[i].active {
			continue
		}

		if testCollision(g.gopherX, groundY-gopherHeight, gopherWidth, gopherHeight, g.coins[i].x, g.coins[i].y, coinWidth, coinHeight) {
			g.coins[i].active = false
			g.score++
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	// Draw ground
	groundOpts := ebiten.DrawImageOptions{}
	groundOpts.GeoM.Translate(0, screenHeight-groundHeight)
	screen.DrawImage(g.groundImage, &groundOpts)

	// Draw gopher
	gopherImageRect := image.Rect(0, 0, gopherPixelsWidth, gopherPixelsHeight)
	gopherImage := g.gopherImage.SubImage(gopherImageRect).(*ebiten.Image)

	xFlip := 1.0
	xAdd := 0.0
	if g.gopherFacingRight {
		xFlip = -1.0
		xAdd = gopherWidth
	}
	gopherOpts := ebiten.DrawImageOptions{}
	gopherOpts.GeoM.Scale(xFlip*gopherScale, gopherScale)
	gopherOpts.GeoM.Translate(g.gopherX+xAdd, groundY-gopherHeight)

	screen.DrawImage(gopherImage, &gopherOpts)

	// Coin
	for i := 0; i < maxCoins; i++ {
		if !g.coins[i].active {
			continue
		}

		coinOpts := ebiten.DrawImageOptions{}
		coinOpts.GeoM.Translate(g.coins[i].x, g.coins[i].y)
		screen.DrawImage(g.coinImage, &coinOpts)
	}

	// Score
	scoreTextOpts := &text.DrawOptions{}
	scoreTextOpts.GeoM.Translate(10, 10)
	text.Draw(screen, fmt.Sprintf("%03d", g.score), &text.GoTextFace{
		Source: g.font,
		Size:   24,
	}, scoreTextOpts)

	// Lives
	for i := 0; i < g.lives; i++ {
		x := screenWidth - (i+1)*(gopherPixelsWidth*livesScale) - 10
		liveImgOpts := ebiten.DrawImageOptions{}
		liveImgOpts.GeoM.Scale(livesScale, livesScale)
		liveImgOpts.GeoM.Translate(float64(x), 10)
		screen.DrawImage(g.livesImage, &liveImgOpts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.Init()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(screenTitle)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
