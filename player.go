package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	lua "github.com/yuin/gopher-lua"
)

// Player is the player in the game
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
	isDead      bool      // Is the player dead?

	canBlinkTimer int        // Timer between blinks
	endBlinkTimer int        // Timer for each blink
	isBlinking    bool       // If blinking
	blinkTrail    BlinkTrail // The animated blue trail

	shadow           Shadow           // The shadow below the player :)
	lightID          int              // The ID of the light that follows the player
	walkSmokeEmitter WalkSmokeEmitter // Emits smoke behind player when he walks

	spritesheet     Spritesheet           // Current spritesheet
	animation       Animation             // Current Animation
	animations      PlayerAnimations      // All animations
	animationSpeeds PlayerAnimationSpeeds // All animation speeds
	isDrawable      bool                  // Is able to be drawn on the screen?

	gun      Gun     // The players gun
	accuracy int     // Player's accuracy with firearms!
	gunRange float64 // Player's range with firearms!

	items []Item // Items held!

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

func createPlayer(position Vec2f, game *Game, lightID int) Player {

	health := 9

	energy := 9
	energyTimer := 120
	sprintEnergyTimer := 20

	canBlinkTimer := 0
	endBlinkTimer := 0

	walkSpeed := 0.9
	runSpeed := 1.3

	shadowRect := image.Rect(0, 231, 14, 237)

	image := iplayerSpritesheet
	// Idle
	idleFrontSpritesheet := createSpritesheet(newVec2i(103, 114), newVec2i(129, 132), 2, image)
	idleBackSpritesheet := createSpritesheet(newVec2i(103, 133), newVec2i(129, 151), 2, image)
	idleLeftSpritesheet := createSpritesheet(newVec2i(103, 95), newVec2i(129, 113), 2, image)
	idleRightSpritesheet := createSpritesheet(newVec2i(103, 76), newVec2i(129, 94), 2, image)
	// Running
	runningFrontSpritesheet := createSpritesheet(newVec2i(103, 57), newVec2i(155, 75), 4, image)
	runningBackSpritesheet := createSpritesheet(newVec2i(103, 38), newVec2i(155, 56), 4, image)
	runningLeftSpritesheet := createSpritesheet(newVec2i(103, 19), newVec2i(155, 37), 4, image)
	runningRightSpritesheet := createSpritesheet(newVec2i(103, 0), newVec2i(155, 18), 4, image)

	return Player{
		position:  position,
		walkSpeed: walkSpeed,
		runSpeed:  runSpeed,

		health:    10,
		maxHealth: health,

		energy:               energy,
		maxEnergy:            energy,
		energyRegenTimer:     energyTimer,
		energyRegenTimerMax:  energyTimer,
		sprintEnergyTimer:    sprintEnergyTimer,
		sprintEnergyTimerMax: sprintEnergyTimer,

		dynamicSize: newVec2i(0, 0),                          // Dynamic size
		staticSize:  runningRightSpritesheet.sprites[0].size, // Static size

		direction:   Down,    // Direction
		movement:    Walking, // Movement
		isMoving:    false,
		isConscious: true,
		isDead:      false,

		canBlinkTimer: canBlinkTimer, // Time between blinks
		endBlinkTimer: endBlinkTimer, // Time for each blink
		isBlinking:    false,
		blinkTrail:    createBlinkTrail(0.5),

		shadow:           createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game)),
		lightID:          lightID,
		walkSmokeEmitter: createWalkSmokeEmitter(2., iwalkSmokeSpritesheet),

		spritesheet: idleFrontSpritesheet,                         // Current animation spritesheet
		animation:   createAnimation(idleFrontSpritesheet, image), // Current animation
		animations: PlayerAnimations{ // All animations
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
		animationSpeeds: PlayerAnimationSpeeds{ // All animation speeds
			1,   // idle
			1.2, // walking
			1.7, // running
			3,   // blink trail
		},
		isDrawable: true,

		// The player's gun is defined here
		gun: Gun{
			position: position,
			image:    iitemsSpritesheet,
			sprite:   createSprite(newVec2i(0, 46), newVec2i(21, 59), newVec2i(21, 13), iitemsSpritesheet),

			fireSpeed:    0,
			firespeedMax: 10,

			baseDamage: 1,

			animation:      createAnimation(createSpritesheet(Vec2i{0, 81}, Vec2i{42, 95}, 2, iitemsSpritesheet), iitemsSpritesheet),
			animationSpeed: 1.5,
		},
		accuracy: 75,
		gunRange: 25,

		items: []Item{},

		image: image, // Entire spritesheet
	}
}

func (p *Player) update(g *Game) {
	switch p.movement {
	case (Idle):
		p.animation.update(p.animationSpeeds.idle)
	case (Walking):
		p.animation.update(p.animationSpeeds.walking)
	case (Running):
		p.animation.update(p.animationSpeeds.running)
	}
	p.input(g)
	go p.updateLevels(g)

	// Blink update
	p.updateBlinkTrail()
	p.blinkTrail.update()

	// Item/Lua update
	p.updateItems()

	// Gun update
	p.gun.update(
		newVec2f(p.position.x+float64(p.dynamicSize.x)/2, p.position.y+float64(p.dynamicSize.y)/2),
		newVec2i(g.cursor.center.x, g.cursor.center.y),
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
	g.lightHandler.lights[g.lightHandler.getLightIndex(p.lightID)].update(
		newVec2f(
			p.position.x+float64(p.staticSize.x/2),
			p.position.y+float64(p.staticSize.y/2),
		),
		0,
	)

	spawnWalkSmoke := false
	if p.movement == Running && p.isMoving {
		spawnWalkSmoke = true
	}

	p.walkSmokeEmitter.update(p.position, p.dynamicSize, p.direction, spawnWalkSmoke)

	if p.isDead {
		p.isDrawable = false
		p.shadow.isDrawable = false
	}

}

func (p *Player) render(screen *ebiten.Image) {

	p.renderWalkSmoke(screen)

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

func (p *Player) renderWalkSmoke(screen *ebiten.Image) {
	p.walkSmokeEmitter.render(screen)
}

func (p *Player) input(g *Game) {
	// Reset moving
	p.isMoving = false

	// Check if a direction this turn was already decided
	directionDecided := false

	// Blink
	p.blinkHandler()

	// Mouse button input
	p.mouseButtonInput(g)

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
func (p *Player) mouseButtonInput(g *Game) {
	if p.isConscious {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.gun.fire(p.position, g)
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
func (p *Player) updateLevels(g *Game) {
	if p.health < 0 {
		p.health = 0
	}
	// DEATH
	if p.health == 0 {
		g.state = 0
		p.die(g)
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

// Updates items and Lua functions, also deals with returns
func (p *Player) runLuaFunctions() {
	for i := 0; i < len(p.items); i++ {
		for j := 0; j < len(p.items[i].functions); j++ {
			// Check if the function is done running
			if p.items[i].functions[j].isFinished {
				continue
			}

			luaReturnChan := make(chan []lua.LValue)
			go p.items[i].runLuaFunction(p.items[i].functions[j], luaReturnChan) // Get returns

			returns := <-luaReturnChan

			p.items[i].functions[j].isFinished = bool(returns[0].(lua.LBool)) // Is it finished?
		}
	}
}

// Basically just a wrapper function for runLuaFunctions for now :)
func (p *Player) updateItems() {
	p.runLuaFunctions()
}

// Collision bounds
func (p *Player) getBoundsDynamic() image.Rectangle {
	return image.Rect(int(p.position.x), int(p.position.y), int(p.position.x)+p.dynamicSize.x, int(p.position.y)+p.dynamicSize.y)
}

func (p *Player) getBoundsStatic() image.Rectangle {
	return image.Rect(int(p.position.x), int(p.position.y), int(p.position.x)+p.staticSize.x, int(p.position.y)+p.staticSize.y)
}

// ----------------

// die
func (p *Player) die(g *Game) {
	gibAmount := 10 // Gib setting 1
	gibSize := 6    // Gib setting 1

	switch g.settings.Graphics.Gibs { // Gib setting 2
	case 2:
		gibAmount = 20
		gibSize = 6
	case 0:
		gibAmount = 0
	}

	gibHandler := createGibHandler() // Gibs

	gibHandler.explode(
		gibAmount,
		gibSize,
		newVec2f(p.position.x+float64(p.dynamicSize.x/2), p.position.y+float64(p.dynamicSize.y/2)), // Center
		p.animation.spritesheet.sprites[p.animation.currentFrame].getBounds(),
		p.image,
	)

	g.gibHandlers = append(g.gibHandlers, gibHandler)

	// TODO: Delete shadow
	//removeShadow(g.shadows, p.shadow.id)

	p.isDead = true
}
