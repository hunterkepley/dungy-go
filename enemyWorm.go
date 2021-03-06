package main

import (
	"image"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	pathfinding "github.com/xarg/gopathfinding"
)

// WormAnimations is the animations for worms
type WormAnimations struct {
	idleFront Animation
}

// WormAnimationSpeeds is the animation speeds for worms
type WormAnimationSpeeds struct {
	idle float64
}

// Worm is a worm type enemy
type Worm struct {
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
	animations      WormAnimations
	animationSpeeds WormAnimationSpeeds

	image *ebiten.Image

	astarNodes  []pathfinding.Node
	path        paths.Path
	canPathfind bool
	pathFinding bool
}

func createWorm(position Vec2f, game *Game) *Worm {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(68, 22), 4, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	return &Worm{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.5,
		weight:    0.3,

		health:       10,
		maxHealth:    10,
		dead:         false,
		attackRadius: 40,

		shadow: &shadow,

		spritesheet: idleFrontSpritesheet,
		animations: WormAnimations{
			idleFront: createAnimation(idleFrontSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: WormAnimationSpeeds{
			idle: 0.9,
		},

		astarNodes:  []pathfinding.Node{},
		pathFinding: false,
		canPathfind: true,

		image: ienemiesSpritesheet,
	}
}

func (w *Worm) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// FLIP DECIDER
	flip := newVec2f(1, 1)
	if w.flipped {
		flip.x = -1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-w.size.x)/2, float64(0-w.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Translate(float64(w.size.x)/2, float64(w.size.y)/2)
	w.subImageRect = image.Rect(
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.y,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(w.position.x), float64(w.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	// Knockback (turning red) render
	w.knockedBack, w.knockedBackTimer = enemyKnockbackRender(op, w.knockedBack, w.knockedBackTimer)

	screen.DrawImage(w.image.SubImage(w.subImageRect).(*ebiten.Image), op) // Draw worm
}

func (w *Worm) update(game *Game) {
	// Start the animation if it's not playing
	if w.animation.state != AnimationPlayingForwards {
		w.animation = w.animations.idleFront
		w.animation.startForwards()
	}
	w.animation.update(w.animationSpeeds.idle)

	// Move worm
	w.position.x += w.velocity.x
	w.position.y += w.velocity.y

	w.size = newVec2i(
		w.spritesheet.sprites[w.animation.currentFrame].size.x,
		w.spritesheet.sprites[w.animation.currentFrame].size.y,
	)
	endPosition := newVec2i(
		int(w.position.x)+w.size.x,
		int(w.position.y)+w.size.y,
	)

	// Update shadow
	w.shadow.update(w.position, w.size)

	w.subImageRect = image.Rect(int(w.position.x), int(w.position.y), endPosition.x, endPosition.y)
	w.center = newVec2f(w.position.x+float64(w.size.x)/2, w.position.y+float64(w.size.y)/2)

	// Attack

	w.attack(game)
}

func (w *Worm) isDead() bool {
	if w.health <= 0 {
		if !w.dead {
			w.dead = true
		}
		return true
	}
	return false
}

func (w *Worm) attack(game *Game) {
	if w.attacking {

	} else if !w.attacking && isCircularCollision(game.player.getBoundsDynamic(), image.Rect(int(w.center.x-w.attackRadius), int(w.center.y-w.attackRadius), int(w.center.x+w.attackRadius), int(w.center.y+w.attackRadius))) {
		w.attacking = true
	}
}

func (w *Worm) getCenter() Vec2f {
	return w.center
}

func (w *Worm) getPosition() Vec2f {
	return w.position
}

func (w *Worm) getCurrentSubImageRect() image.Rectangle {
	return w.subImageRect
}

func (w *Worm) getImage() *ebiten.Image {
	return w.image
}

func (w *Worm) getSize() Vec2i {
	return w.size
}

func (w *Worm) damage(value int) {
	w.health--
}

func (w *Worm) getShadow() Shadow {
	return *w.shadow
}

func (w *Worm) getMoveSpeed() float64 {
	return w.moveSpeed
}

func (w *Worm) getDying() bool {
	return w.dying
}

func (w *Worm) getAttacking() bool {
	return w.attacking
}

func (w *Worm) getPath() *paths.Path {
	return &w.path
}

func (w *Worm) getCanPathfind() bool {
	return w.canPathfind
}

func (w *Worm) getFlipped() bool {
	return w.flipped
}

func (w *Worm) getPathfinding() bool {
	return w.pathFinding
}

func (w *Worm) setPosition(pos Vec2f) {
	w.position = pos
}

func (w *Worm) setFlipped(flipped bool) {
	w.flipped = flipped
}

func (w *Worm) setPath(path paths.Path) {
	w.path = path
}

func (w *Worm) setCanPathfind(canPathfind bool) {
	w.canPathfind = canPathfind
}

func (w *Worm) getWeight() float64 {
	return w.weight
}

func (w *Worm) getKnockedBack() bool {
	return w.knockedBack
}

func (w *Worm) setKnockedBack(knockedBack bool) {
	w.knockedBack = knockedBack
}
