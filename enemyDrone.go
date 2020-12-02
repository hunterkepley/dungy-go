package main

import (
	"image"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	pathfinding "github.com/xarg/gopathfinding"
)

// DroneAnimations is the animations for Drones
type DroneAnimations struct {
	idleFront Animation
}

// DroneAnimationSpeeds is the animation speeds for Drones
type DroneAnimationSpeeds struct {
	idle float64
}

// Drone is a Drone type enemy
type Drone struct {
	position  Vec2f
	center    Vec2f
	size      Vec2i
	velocity  Vec2f
	moveSpeed float64
	weight    float64

	health                 int
	maxHealth              int
	dead                   bool
	deathExplosion         bool // When the death explosion is playing
	deathExplosionFinished bool // When the death explosion is finished
	dying                  bool
	remove                 bool    // Do we remove this enemy?
	flipped                bool    // Is the enemy flipped?
	idle                   bool    // Is the enemy idling?
	attacking              bool    // Is the enemy attacking?
	attackRadius           float64 // When the player is in this radius, the enemy will attack!
	knockedBack            bool    // Is the enemy being knocked back?
	knockedBackTimer       float64

	shadow *Shadow // The shadow below the enemy

	subImageRect image.Rectangle

	spritesheet     Spritesheet
	animation       Animation
	animations      DroneAnimations
	animationSpeeds DroneAnimationSpeeds

	image *ebiten.Image

	astarNodes  []pathfinding.Node
	path        paths.Path
	canPathfind bool
	pathFinding bool
}

func createDrone(position Vec2f, game *Game) *Drone {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(68, 22), 4, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	return &Drone{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.4,
		weight:    0.3,

		health:       15,
		maxHealth:    15,
		dead:         false,
		attackRadius: 60,

		shadow: &shadow,

		spritesheet: idleFrontSpritesheet,
		animations: DroneAnimations{
			idleFront: createAnimation(idleFrontSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: DroneAnimationSpeeds{
			idle: 0.9,
		},

		astarNodes:  []pathfinding.Node{},
		pathFinding: false,
		canPathfind: true,

		image: ienemiesSpritesheet,
	}
}

func (d *Drone) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// FLIP DECIDER
	flip := newVec2f(1, 1)
	if d.flipped {
		flip.x = -1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-d.size.x)/2, float64(0-d.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Translate(float64(d.size.x)/2, float64(d.size.y)/2)
	d.subImageRect = image.Rect(
		d.spritesheet.sprites[d.animation.currentFrame].startPosition.x,
		d.spritesheet.sprites[d.animation.currentFrame].startPosition.y,
		d.spritesheet.sprites[d.animation.currentFrame].endPosition.x,
		d.spritesheet.sprites[d.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(d.position.x), float64(d.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	// Knockback (turning red) render
	d.knockedBack, d.knockedBackTimer = enemyKnockbackRender(op, d.knockedBack, d.knockedBackTimer)

	screen.DrawImage(d.image.SubImage(d.subImageRect).(*ebiten.Image), op) // Draw Drone
}

func (d *Drone) update(game *Game) {
	// Start the animation if it's not playing
	if d.animation.state != AnimationPlayingForwards {
		d.animation = d.animations.idleFront
		d.animation.startForwards()
	}
	d.animation.update(d.animationSpeeds.idle)

	// Move Drone
	d.position.x += d.velocity.x
	d.position.y += d.velocity.y

	d.size = newVec2i(
		d.spritesheet.sprites[d.animation.currentFrame].size.x,
		d.spritesheet.sprites[d.animation.currentFrame].size.y,
	)
	endPosition := newVec2i(
		int(d.position.x)+d.size.x,
		int(d.position.y)+d.size.y,
	)

	// Update shadow
	d.shadow.update(d.position, d.size)

	d.subImageRect = image.Rect(int(d.position.x), int(d.position.y), endPosition.x, endPosition.y)
	d.center = newVec2f(d.position.x+float64(d.size.x)/2, d.position.y+float64(d.size.y)/2)

	// Attack

	d.attack(game)
}

func (d *Drone) isDead() bool {
	if d.health <= 0 {
		if !d.dead {
			d.dead = true
		}
		return true
	}
	return false
}

func (d *Drone) attack(game *Game) {

}

func (d *Drone) getCenter() Vec2f {
	return d.center
}

func (d *Drone) getPosition() Vec2f {
	return d.position
}

func (d *Drone) getCurrentSubImageRect() image.Rectangle {
	return d.subImageRect
}

func (d *Drone) getImage() *ebiten.Image {
	return d.image
}

func (d *Drone) getSize() Vec2i {
	return d.size
}

func (d *Drone) damage(value int) {
	d.health--
}

func (d *Drone) getShadow() Shadow {
	return *d.shadow
}

func (d *Drone) getMoveSpeed() float64 {
	return d.moveSpeed
}

func (d *Drone) getDying() bool {
	return d.dying
}

func (d *Drone) getAttacking() bool {
	return d.attacking
}

func (d *Drone) getPath() *paths.Path {
	return &d.path
}

func (d *Drone) getCanPathfind() bool {
	return d.canPathfind
}

func (d *Drone) getFlipped() bool {
	return d.flipped
}

func (d *Drone) getPathfinding() bool {
	return d.pathFinding
}

func (d *Drone) setPosition(pos Vec2f) {
	d.position = pos
}

func (d *Drone) setFlipped(flipped bool) {
	d.flipped = flipped
}

func (d *Drone) setPath(path paths.Path) {
	d.path = path
}

func (d *Drone) setCanPathfind(canPathfind bool) {
	d.canPathfind = canPathfind
}

func (d *Drone) getWeight() float64 {
	return d.weight
}

func (d *Drone) getKnockedBack() bool {
	return d.knockedBack
}

func (d *Drone) setKnockedBack(knockedBack bool) {
	d.knockedBack = knockedBack
}
