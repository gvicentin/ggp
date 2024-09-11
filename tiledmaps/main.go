package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const (
	mapPath = "map.tmx"

	screenTitle  = "Rendering map"
	screenWidth  = 640
	screenHeight = 480

	scale = 5
)

type Game struct {
	mapImage *ebiten.Image

	cameraX float64
	cameraY float64
}

func (g *Game) init() {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFile(mapPath)
	if err != nil {
		log.Fatalf("error loading map: %s", err.Error())
	}

	fmt.Println(gameMap)

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		log.Fatalf("map unsupported for rendering: %s", err.Error())
	}

	// Render the map to the output image.
	err = renderer.RenderVisibleLayers()
	if err != nil {
		log.Fatalf("layer unsupported for rendering: %s", err.Error())
	}

	var buff []byte
	buffer := bytes.NewBuffer(buff)

	renderer.SaveAsPng(buffer)

	im, err := png.Decode(buffer)

	g.mapImage, _ = ebiten.NewImageFromImage(im, ebiten.FilterNearest)
}

func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cameraY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cameraY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraX += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraX -= 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x25, 0x13, 0x1a, 0xff})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.cameraX, g.cameraY)
	op.GeoM.Scale(scale, scale)
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(g.mapImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.init()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(screenTitle)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
