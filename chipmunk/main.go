package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jakecoffman/cp"
)

const (
	screenWidth  = 600
	screenHeight = 480
)

var (
	dot = ebiten.NewImage(1, 1)
)

func init() {
	dot.Fill(color.White)
}

type Game struct {
	space *cp.Space

	playerBody  *cp.Body
	playerShape *cp.Shape

	lastJumpState bool
}

func NewGame() *Game {
	space := cp.NewSpace()
	space.Iterations = 10
	space.SetGravity(cp.Vector{X: 0, Y: 100})

	walls := []cp.Vector{
		{X: 0, Y: screenHeight - 20}, {X: screenWidth, Y: screenHeight - 20},
	}

	for wall := 0; wall < len(walls)-1; wall++ {
		shape := space.AddShape(cp.NewSegment(space.StaticBody, walls[wall], walls[wall+1], 0))
		shape.SetElasticity(0)
		shape.SetFriction(0)
	}

	// player
	playerBody := space.AddBody(cp.NewBody(1, cp.INFINITY))
	playerBody.SetPosition(cp.Vector{X: 100, Y: 225})
	playerBody.SetVelocityUpdateFunc(playerUpdateVelocity)

	playerShape := space.AddShape(cp.NewBox(playerBody, 20, 20, 0))
	playerShape.SetElasticity(0)
	playerShape.SetFriction(0)
	playerShape.SetCollisionType(1)

	return &Game{
		space:       space,
		playerBody:  playerBody,
		playerShape: playerShape,
	}
}

func playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	var targetX float64
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		targetX = -400
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		targetX = 400
	}

	body.UpdateVelocity(gravity, damping, dt)

	velocity := body.Velocity()
	body.SetVelocityVector(cp.Vector{X: targetX, Y: velocity.Y})
}

func (g *Game) Update() error {
	jumpState := ebiten.IsKeyPressed(ebiten.KeySpace)

	if jumpState && !g.lastJumpState {
		g.playerBody.SetVelocityVector(cp.Vector{X: 0, Y: -100})
	}

	g.space.Step(1.0 / float64(ebiten.TPS()))

	g.lastJumpState = jumpState

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	vector.DrawFilledRect(screen, 0, screenHeight-10, screenWidth, 10, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)

	playerPos := g.playerBody.Position()
	vector.DrawFilledRect(screen, float32(playerPos.X), float32(playerPos.Y), 20, 20, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
