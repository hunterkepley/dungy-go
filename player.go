package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Player is the player in the fuckin game I hate useless comments
type Player struct {
	position Vec2f
	//center   Vec2f
	speed       float64
	direction   Direction
	spritesheet Spritesheet      // Current spritesheet
	animation   Animation        // Current Animation
	animations  PlayerAnimations // All animations
	image       *ebiten.Image
}

// PlayerAnimations is the animations for the player
type PlayerAnimations struct {
	idleFront Animation
	idleBack  Animation
	idleLeft  Animation
	idleRight Animation
}

func createPlayer(position Vec2f) Player {
	speed := 3.
	image := i_playerSpritesheet
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(75, 26), 5, image)
	idleBackSpritesheet := createSpritesheet(newVec2i(26, 0), newVec2i(75, 52), 5, image)
	idleLeftSpritesheet := createSpritesheet(newVec2i(52, 0), newVec2i(75, 78), 5, image)
	idleRightSpritesheet := createSpritesheet(newVec2i(78, 0), newVec2i(75, 104), 5, image)
	return Player{
		position,
		speed,
		Down,                 // Direction
		idleFrontSpritesheet, // Current animation spritesheet
		createAnimation(idleFrontSpritesheet, image), // Current animation
		PlayerAnimations{ // All animations
			createAnimation(idleFrontSpritesheet, image),
			createAnimation(idleBackSpritesheet, image),
			createAnimation(idleLeftSpritesheet, image),
			createAnimation(idleRightSpritesheet, image),
		},
		image, // Entire spritesheet
	}
}

func (p *Player) update() {
	p.animation.play(0.4)
	p.input()
}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.x, p.position.y)
	subImageRect := image.Rect(
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.y,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.y,
	)
	screen.DrawImage(p.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (p *Player) input() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.position.y -= p.speed
		if p.direction != Up {
			p.changeUp()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.position.y += p.speed
		if p.direction != Down {
			p.changeDown()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.position.x -= p.speed
		if p.direction != Left {
			p.changeLeft()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.position.x += p.speed
		if p.direction != Right {
			p.changeRight()
		}
	}
}

func (p *Player) changeUp() {
	p.direction = Up
	p.animation = p.animations.idleBack
	p.spritesheet = p.animation.spritesheet
}

func (p *Player) changeDown() {
	p.direction = Down
	p.animation = p.animations.idleFront
	p.spritesheet = p.animation.spritesheet
}

func (p *Player) changeLeft() {
	p.direction = Left
	p.animation = p.animations.idleLeft
	p.spritesheet = p.animation.spritesheet
}

func (p *Player) changeRight() {
	p.direction = Right
	p.animation = p.animations.idleRight
	p.spritesheet = p.animation.spritesheet
}

// ^ DIRECTION CHANGES ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
