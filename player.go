package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Player is the player in the fuckin game I hate useless comments
type Player struct {
	position Vec2f
	//center   Vec2f
	walkSpeed float64
	runSpeed  float64

	direction Direction // Up? Down?
	movement  Movement  // Walking? Running?
	isMoving  bool      // Is currently moving due to input?

	spritesheet     Spritesheet           // Current spritesheet
	animation       Animation             // Current Animation
	animations      PlayerAnimations      // All animations
	animationSpeeds PlayerAnimationSpeeds // All animation speeds

	image *ebiten.Image
}

// PlayerAnimations is the animations for the player
type PlayerAnimations struct {
	idleFront    Animation
	idleBack     Animation
	idleLeft     Animation
	idleRight    Animation
	runningFront Animation
	runningBack  Animation
	runningLeft  Animation
	runningRight Animation
}

// PlayerAnimationSpeeds is the animation speeds for the player
type PlayerAnimationSpeeds struct {
	idle    float64
	walking float64
	running float64
}

func createPlayer(position Vec2f) Player {
	walkSpeed := 1.
	runSpeed := 1.5
	image := iplayerSpritesheet
	// Idle
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(60, 26), 5, image)
	idleBackSpritesheet := createSpritesheet(newVec2i(0, 26), newVec2i(60, 52), 5, image)
	idleLeftSpritesheet := createSpritesheet(newVec2i(0, 52), newVec2i(70, 78), 5, image)
	idleRightSpritesheet := createSpritesheet(newVec2i(0, 78), newVec2i(70, 104), 5, image)
	// Running
	runningFrontSpritesheet := createSpritesheet(newVec2i(0, 104), newVec2i(84, 128), 6, image)
	runningBackSpritesheet := createSpritesheet(newVec2i(0, 129), newVec2i(84, 154), 6, image)
	runningLeftSpritesheet := createSpritesheet(newVec2i(0, 154), newVec2i(90, 180), 6, image)
	runningRightSpritesheet := createSpritesheet(newVec2i(0, 180), newVec2i(90, 206), 6, image)
	return Player{
		position,
		walkSpeed,
		runSpeed,

		Down,    // Direction
		Walking, // Movement
		false,

		idleFrontSpritesheet,                         // Current animation spritesheet
		createAnimation(idleFrontSpritesheet, image), // Current animation
		PlayerAnimations{ // All animations
			// Idle
			createAnimation(idleFrontSpritesheet, image),
			createAnimation(idleBackSpritesheet, image),
			createAnimation(idleLeftSpritesheet, image),
			createAnimation(idleRightSpritesheet, image),
			// Running
			createAnimation(runningFrontSpritesheet, image),
			createAnimation(runningBackSpritesheet, image),
			createAnimation(runningLeftSpritesheet, image),
			createAnimation(runningRightSpritesheet, image),
		},
		PlayerAnimationSpeeds{ // All animation speeds
			1,   // idle
			1.6, // walking
			2.3, // running
		},

		image, // Entire spritesheet
	}
}

func (p *Player) update() {
	switch p.movement {
	case (Idle):
		p.animation.play(p.animationSpeeds.idle)
		break
	case (Walking):
		p.animation.play(p.animationSpeeds.walking)
		break
	case (Running):
		p.animation.play(p.animationSpeeds.running)
		break
	}
	p.input()
}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.x, p.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	subImageRect := image.Rect(
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.y,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.y,
	)
	screen.DrawImage(p.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (p *Player) input() {
	// Reset moving
	p.isMoving = false

	// Check if a direction this turn was already decided
	directionDecided := false

	if ebiten.IsKeyPressed(ebiten.KeyA) { // LEFT

		if p.movement == Walking {
			p.position.x -= p.walkSpeed
		} else if p.movement == Running {
			p.position.x -= p.runSpeed
		}

		if !directionDecided {
			p.changeLeft()
		}

		p.isMoving = true
		directionDecided = true
	} else if ebiten.IsKeyPressed(ebiten.KeyD) { // RIGHT

		if p.movement == Walking {
			p.position.x += p.walkSpeed
		} else if p.movement == Running {
			p.position.x += p.runSpeed
		}

		if !directionDecided {
			p.changeRight()
		}

		p.isMoving = true
		directionDecided = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) { // UP

		if p.movement == Walking {
			p.position.y -= p.walkSpeed
		} else if p.movement == Running {
			p.position.y -= p.runSpeed
		}

		if !directionDecided {
			p.changeUp()
		}

		p.isMoving = true
		directionDecided = true
	} else if ebiten.IsKeyPressed(ebiten.KeyS) { // DOWN

		if p.movement == Walking {
			p.position.y += p.walkSpeed
		} else if p.movement == Running {
			p.position.y += p.runSpeed
		}

		if !directionDecided {
			p.changeDown()
		}

		p.isMoving = true
		directionDecided = true
	}

	// Decide movement
	// Not moving
	if !p.isMoving {
		p.movement = Idle
		switch p.direction {
		case (Up):
			p.changeAnimation(p.animations.idleBack)
			break
		case (Down):
			p.changeAnimation(p.animations.idleFront)
			break
		case (Left):
			p.changeAnimation(p.animations.idleLeft)
			break
		case (Right):
			p.changeAnimation(p.animations.idleRight)
			break
		}
	} else {
		// Moving
		if p.movement != Walking {
			// Change to walking before checking if sprinting [to walk again after sprinting]
			p.movement = Walking
		}
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// If holding shift, change to running if not already running!
			if p.movement != Running {
				p.movement = Running
			}
		}
	}
}

func (p *Player) changeUp() {
	p.direction = Up
	if p.movement == Idle {
		p.changeAnimation(p.animations.idleBack)
	} else {
		p.changeAnimation(p.animations.runningBack)
	}
}

func (p *Player) changeDown() {
	p.direction = Down
	if p.movement == Idle {
		p.changeAnimation(p.animations.idleFront)
	} else {
		p.changeAnimation(p.animations.runningFront)
	}
}

func (p *Player) changeLeft() {
	p.direction = Left
	if p.movement == Idle {
		p.changeAnimation(p.animations.idleLeft)
	} else {
		p.changeAnimation(p.animations.runningLeft)
	}
}

func (p *Player) changeRight() {
	p.direction = Right
	if p.movement == Idle {
		p.changeAnimation(p.animations.idleRight)
	} else {
		p.changeAnimation(p.animations.runningRight)
	}
}

// ^ DIRECTION CHANGES ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func (p *Player) changeAnimation(animation Animation) {
	// Only switch animation if not already the current animation
	if p.animation.id != animation.id {
		p.animation = animation
		p.spritesheet = p.animation.spritesheet
	}
}
