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

	health    int
	maxHealth int

	energy               int // Energy!
	maxEnergy            int // How much energy he can have max
	energyRegenTimer     int // Regeneration timer
	energyRegenTimerMax  int // Regeneration timer tick rate
	sprintEnergyTimer    int // Keeps track of how fast sprint is depleting energy
	sprintEnergyTimerMax int // Sprint energy depletion tick rate

	dynamicSize Vec2i // This is the player's dynamic size
	staticSize  Vec2i // This value is the player's largest size for wall collisions

	direction   Direction // Up? Down?
	movement    Movement  // Walking? Running?
	isMoving    bool      // Is currently moving due to input?
	isConscious bool      // Is player conscious? [Can use input?]

	canBlinkTimer int        // Timer between blinks
	endBlinkTimer int        // Timer for each blink
	isBlinking    bool       // If blinking
	blinkTrail    BlinkTrail // The animated blue trail

	shadow Shadow // The shadow below the player :)

	spritesheet     Spritesheet           // Current spritesheet
	animation       Animation             // Current Animation
	animations      PlayerAnimations      // All animations
	animationSpeeds PlayerAnimationSpeeds // All animation speeds
	isDrawable      bool                  // Is able to be drawn on the screen?

	gun Gun // The players gun

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
	idle       float64
	walking    float64
	running    float64
	blinkTrail float64
}

func createPlayer(position Vec2f) Player {

	health := 9

	energy := 9
	energyTimer := 120
	sprintEnergyTimer := 20

	canBlinkTimer := 0
	endBlinkTimer := 0

	walkSpeed := 0.8
	runSpeed := 1.6

	shadowRect := image.Rect(0, 231, 14, 237)

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

		health,
		health,

		energy,
		energy,
		energyTimer,
		energyTimer,
		sprintEnergyTimer,
		sprintEnergyTimer,

		newVec2i(0, 0),                          // Dynamic size
		runningRightSpritesheet.sprites[0].size, // Static size

		Down,    // Direction
		Walking, // Movement
		false,
		true,

		canBlinkTimer, // Time between blinks
		endBlinkTimer, // Time for each blink
		false,
		createBlinkTrail(0.5),

		createShadow(shadowRect, iplayerSpritesheet),

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
			1.3, // walking
			2.5, // running
			3,   // blink trail
		},
		true,

		Gun{
			position:     position,
			image:        iitemsSpritesheet,
			sprite:       createSprite(newVec2i(0, 46), newVec2i(21, 59), newVec2i(21, 13), iitemsSpritesheet),
			fireSpeed:    0,
			firespeedMax: 50,
		},

		image, // Entire spritesheet
	}
}

func (p *Player) update(cursor Cursor) {
	switch p.movement {
	case (Idle):
		p.animation.update(p.animationSpeeds.idle)
	case (Walking):
		p.animation.update(p.animationSpeeds.walking)
	case (Running):
		p.animation.update(p.animationSpeeds.running)
	}
	p.input(cursor)
	go p.updateLevels()

	// Blink update
	p.updateBlinkTrail()
	p.blinkTrail.update()

	// Gun update
	p.gun.update(
		newVec2f(p.position.x+float64(p.dynamicSize.x)/2, p.position.y+float64(p.dynamicSize.y)/2),
		newVec2i(cursor.center.x, cursor.center.y),
	)

	// Set size
	p.dynamicSize = p.animation.spritesheet.sprites[0].size
	p.wallCollisions()

	// Check if drawable
	if p.isDrawable == p.isBlinking {
		p.isDrawable = !p.isDrawable
		p.isConscious = !p.isConscious // No shoot or move while blink
	}

	p.shadow.isDrawable = p.isDrawable
	p.shadow.update(p.position, p.dynamicSize)
}

func (p *Player) render(screen *ebiten.Image) {

	if p.isDrawable {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.position.x, p.position.y)
		op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

		currentFrame := p.spritesheet.sprites[p.animation.currentFrame]

		subImageRect := image.Rect(
			currentFrame.startPosition.x,
			currentFrame.startPosition.y,
			currentFrame.endPosition.x,
			currentFrame.endPosition.y,
		)

		screen.DrawImage(p.image.SubImage(subImageRect).(*ebiten.Image), op) // Draw player
	}
	p.renderBlinkTrail(screen) // Draw blink trail

}

func (p *Player) renderBlinkTrail(screen *ebiten.Image) {
	p.blinkTrail.render(screen)
}

func (p *Player) updateBlinkTrail() {
	if p.movement != Running && !p.isBlinking {
		go p.energyRegeneration()
	}
	for i, s := range p.blinkTrail.sections {
		if s.delete {
			p.blinkTrail.sections = append(p.blinkTrail.sections[:i], p.blinkTrail.sections[i+1:]...)
		}
	}
}

func (p *Player) input(cursor Cursor) {
	// Reset moving
	p.isMoving = false

	// Check if a direction this turn was already decided
	directionDecided := false

	// Blink
	p.blinkHandler()

	// TEMPORARY
	if ebiten.IsKeyPressed(ebiten.KeyY) {
		p.energy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		p.health++
	}
	// TEMPORARY

	// Mouse button input
	p.mouseButtonInput(cursor)

	if !p.isBlinking {
		// Deplete energy!
		if p.movement == Running {
			if p.sprintEnergyTimer <= 0 {
				p.sprintEnergyTimer = p.sprintEnergyTimerMax
				p.energy--
			} else {
				p.sprintEnergyTimer--
			}
		}
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
			} else { // Then it is two keys at once!
				if p.direction == Left {
					p.direction = UpLeft
				} else {
					p.direction = UpRight
				}
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
			} else {
				if p.direction == Left {
					p.direction = DownLeft
				} else { // Then it is two keys at once!
					p.direction = DownRight
				}
			}

			p.isMoving = true
			directionDecided = true
		}

		// Decide movement
		// Not moving
		if !p.isMoving {
			p.movement = Idle
			switch {
			case p.direction == Up:
				p.changeAnimation(p.animations.idleBack)
			case p.direction == Down:
				p.changeAnimation(p.animations.idleFront)
			case p.direction == Left || p.direction == UpLeft || p.direction == DownLeft:
				p.changeAnimation(p.animations.idleLeft)
			case p.direction == Right || p.direction == UpRight || p.direction == DownRight:
				p.changeAnimation(p.animations.idleRight)
			}
		} else {
			// Moving
			if p.movement != Walking {
				// Change to walking before checking if sprinting [to walk again after sprinting]
				p.movement = Walking
			}
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				// If holding shift, change to running if not already running!
				// Also check if you have energy!
				if p.movement != Running && p.energy > 0 {
					p.movement = Running
				}
			}
		}
	}
}

// Mouse input!
func (p *Player) mouseButtonInput(cursor Cursor) {
	if p.isConscious {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.gun.fire(p.position, cursor.center)
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			// Maybe a charge shot or special ability?
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

// BLINK
func (p *Player) blinkHandler() {
	betweenBlinkTime := 20
	blinkTime := 15

	blinkEnergyDepleter := 3
	// FIX THIS MONSTROSITY AT SOME POINT JESUS CHRIST
	if p.canBlinkTimer >= betweenBlinkTime &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		!p.isBlinking && p.energy >= blinkEnergyDepleter {
		p.energy -= blinkEnergyDepleter // Get rid of some energy and blink!
		p.isBlinking = true
		p.canBlinkTimer = 0
	} else {
		p.canBlinkTimer++
	}

	// If actually blinking
	if p.endBlinkTimer <= blinkTime && p.isBlinking {
		p.blink()
	} else {
		p.isBlinking = false
		p.animation.startForwards()
		p.endBlinkTimer = 0
	}
}

func (p *Player) blink() {
	blinkSpeed := p.runSpeed * 2

	p.blinkTrail.spawnUpdate(p.position, p.direction)

	p.endBlinkTimer++
	p.animation.pause()
	switch p.direction {
	case (Right):
		p.position.x += blinkSpeed
	case (Left):
		p.position.x -= blinkSpeed
	case (Up):
		p.position.y -= blinkSpeed
	case (Down):
		p.position.y += blinkSpeed
	case (UpRight):
		p.position.x += blinkSpeed
		p.position.y -= blinkSpeed
	case (UpLeft):
		p.position.x -= blinkSpeed
		p.position.y -= blinkSpeed
	case (DownRight):
		p.position.x += blinkSpeed
		p.position.y += blinkSpeed
	case (DownLeft):
		p.position.x -= blinkSpeed
		p.position.y += blinkSpeed
	}

}

func (p *Player) changeAnimation(animation Animation) {
	// Only switch animation if not already the current animation
	if p.animation.id != animation.id {
		p.animation = animation
		p.spritesheet = p.animation.spritesheet
	}
}

func (p *Player) wallCollisions() {
	// Left/Right wall width: 17
	// Bottom wall height: 17
	if p.position.x <= 17 {
		p.position.x = 17
	} else if p.position.x+float64(p.staticSize.x) >= screenWidth-17 {
		p.position.x = screenWidth - 17 - float64(p.staticSize.x)
	}
	if p.position.y <= 8 {
		p.position.y = 8
	} else if p.position.y+float64(p.staticSize.y) >= screenHeight-17 {
		p.position.y = screenHeight - 17 - float64(p.staticSize.y)
	}
}

// Updates health, energy, maybe etc
func (p *Player) updateLevels() {
	if p.health < 0 {
		p.health = 0
	}
	if p.health > p.maxHealth {
		p.health = p.maxHealth
	}
	if p.energy < 0 {
		p.energy = 0
	}
	if p.energy > p.maxEnergy {
		p.energy = p.maxEnergy
	}
}

// Deals with energy regeneration
func (p *Player) energyRegeneration() {
	if p.energyRegenTimer <= 0 {
		p.energy++ // Add an energy
		p.energyRegenTimer = p.energyRegenTimerMax
	} else {
		p.energyRegenTimer--
	}
}
