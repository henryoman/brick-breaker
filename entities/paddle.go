package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PaddleWidth  = 120
	PaddleHeight = 20
	// Inertia parameters – tweak for desired feel
	PaddleAccel    = 2500.0 // px/s² when key held (reduced for heavier feel)
	PaddleFriction = 2400.0 // px/s² when no key
	PaddleMaxSpeed = 450.0  // px/s terminal velocity (further reduced)

	PaddleY = 500 // Y position (moved DOWN towards bottom)

	// Gameplay area boundaries (not full screen)
	GameAreaLeft   = 10  // 10px left margin
	GameAreaRight  = 710 // 10px + 700px width
	GameAreaTop    = 30  // 30px top (HUD)
	GameAreaBottom = 530 // 30px + 500px height
	GameAreaWidth  = 700 // gameplay area width
	GameAreaHeight = 500 // gameplay area height

	ScreenWidth = 720        // full screen width
	Tick        = 1.0 / 60.0 // fixed timestep (should match ebiten TPS)
)

// Paddle represents the player's paddle
type Paddle struct {
	x  float64 // center position
	vx float64 // horizontal velocity
}

// NewPaddle creates a new paddle at the center of the gameplay area
func NewPaddle() *Paddle {
	return &Paddle{
		x:  GameAreaLeft + GameAreaWidth/2, // center of gameplay area
		vx: 0,
	}
}

// Update applies acceleration, friction, and updates position – gives the paddle inertia.
func (p *Paddle) Update() {
	// 1. Determine acceleration from input
	ax := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		ax = -PaddleAccel
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		// If both keys held ax cancels to 0 → friction only
		ax = +PaddleAccel
	}

	// 2. If no input apply friction opposite to current velocity
	if ax == 0 {
		if p.vx > 0 {
			ax = -PaddleFriction
			if p.vx+ax*Tick < 0 {
				ax = -p.vx / Tick // prevent overshoot through zero
			}
		} else if p.vx < 0 {
			ax = +PaddleFriction
			if p.vx+ax*Tick > 0 {
				ax = -p.vx / Tick
			}
		}
	}

	// 3. Integrate velocity and clamp to max speed
	p.vx += ax * Tick
	if p.vx > PaddleMaxSpeed {
		p.vx = PaddleMaxSpeed
	}
	if p.vx < -PaddleMaxSpeed {
		p.vx = -PaddleMaxSpeed
	}

	// 4. Integrate position
	p.x += p.vx * Tick

	// 5. Collision with gameplay area edges – stop and zero velocity
	if p.x < GameAreaLeft+PaddleWidth/2 {
		p.x = GameAreaLeft + PaddleWidth/2
		p.vx = 0
	}
	if p.x > GameAreaRight-PaddleWidth/2 {
		p.x = GameAreaRight - PaddleWidth/2
		p.vx = 0
	}
}

// X returns the center X position of the paddle
func (p *Paddle) X() float64 {
	return p.x
}

// Y returns the Y position of the paddle
func (p *Paddle) Y() float64 {
	return PaddleY
}

// Width returns the width of the paddle
func (p *Paddle) Width() float64 {
	return PaddleWidth
}

// Height returns the height of the paddle
func (p *Paddle) Height() float64 {
	return PaddleHeight
}

// GetBounds returns the paddle's bounding box for collision detection
func (p *Paddle) GetBounds() (left, top, right, bottom float64) {
	left = p.x - PaddleWidth/2
	right = p.x + PaddleWidth/2
	top = PaddleY
	bottom = PaddleY + PaddleHeight
	return
}
