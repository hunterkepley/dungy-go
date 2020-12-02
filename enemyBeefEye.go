package main

import (
	"image"
	"math"
	"math/rand"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	pathfinding "github.com/xarg/gopathfinding"
)

// BeefEyeShockwave is each individual shockwave that can hit the player
type BeefEyeShockwave struct {
	position Vec2f

	velocity  Vec2f
	moveSpeed float64
	size      Vec2i

	animation      Animation
	spritesheet    Spritesheet
	animationSpeed float64
	destroy        bool
	damage         int

	subImageRect image.Rectangle
	image        *ebiten.Image
}

func (b *BeefEyeShockwave) update(beefEye BeefEye, i int) {
	animationSpeed := 0.7

	go func() {
		if isAABBCollision(
			image.Rect(int(b.position.x), int(b.position.y), int(b.position.x)+b.size.x, int(b.position.y)+b.size.y),
			image.Rect(
				int(gameReference.player.position.x),
				int(gameReference.player.position.y),
				int(gameReference.player.position.x)+gameReference.player.staticSize.x,
				int(gameReference.player.position.y)+gameReference.player.staticSize.y,
			),
		) {
			b.destroy = true
			if !gameReference.player.isBlinking {
				gameReference.player.health -= b.damage
			}
		}
	}()

	vec := Vec2f{math.Cos(float64(i)), math.Sin(float64(i))}

	// Move away from beefEye
	b.position.x += vec.x * b.moveSpeed
	b.position.y += vec.y * b.moveSpeed

	b.animation.update(animationSpeed)

	if b.animation.finishedFirstPlay {
		b.destroy = true
	}

}

func (b *BeefEyeShockwave) render(screen *ebiten.Image) {
	if len(b.animation.spritesheet.sprites) == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-b.size.x)/2, float64(0-b.size.y)/2)
	// Rotate here
	op.GeoM.Translate(float64(b.size.x)/2, float64(b.size.y)/2)
	b.subImageRect = image.Rect(
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.y,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(b.position.x), float64(b.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(b.image.SubImage(b.subImageRect).(*ebiten.Image), op) // Draw enemy
}

// BeefEyeShockwaveHandler controls the shockwave that is emitted when the BeefEye attacks
type BeefEyeShockwaveHandler struct {
	shockwaves []BeefEyeShockwave
}

func (b *BeefEyeShockwaveHandler) init(beefEye *BeefEye, numShockwaves int) {

	spritesheet := createSpritesheet(Vec2i{0, 71}, Vec2i{39, 84}, 3, ienemiesSpritesheet)

	for i := 0; i < numShockwaves; i++ {
		b.shockwaves = append(b.shockwaves, BeefEyeShockwave{
			position: Vec2f{beefEye.center.x - float64(spritesheet.size.x/spritesheet.numberOfSprites/2), beefEye.center.y - float64(spritesheet.size.y/spritesheet.numberOfSprites/2)},

			moveSpeed: rand.Float64() + 1.5,

			damage: 1,

			spritesheet: spritesheet,

			animation: createAnimation(spritesheet, ienemiesSpritesheet),

			image: ienemiesSpritesheet,
		})
	}
}

func (b *BeefEyeShockwaveHandler) render(screen *ebiten.Image) {
	for i := 0; i < len(b.shockwaves); i++ {
		b.shockwaves[i].render(screen)
	}
}

func (b *BeefEyeShockwaveHandler) update(beefEye BeefEye) {
	for i := 0; i < len(b.shockwaves); i++ {
		b.shockwaves[i].update(beefEye, i)
		if b.shockwaves[i].destroy {
			b.shockwaves = removeShockwave(b.shockwaves, i)
		}
	}
}

func removeShockwave(slice []BeefEyeShockwave, b int) []BeefEyeShockwave {
	return append(slice[:b], slice[b+1:]...)
}

// BeefEyeAnimations is the animations
type BeefEyeAnimations struct {
	walkSide Animation
	die      Animation
	attack   Animation
}

// BeefEyeAnimationSpeeds is the animation speeds
type BeefEyeAnimationSpeeds struct {
	walk   float64
	die    float64
	attack float64
}

// BeefEye is a beefy eye type enemy
type BeefEye struct {
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

	shadow           *Shadow                 // The shadow below the enemy
	shockwaveHandler BeefEyeShockwaveHandler // The shockwave handler for when the BeefEye attacks

	subImageRect image.Rectangle

	spritesheet     Spritesheet
	animation       Animation
	animations      BeefEyeAnimations
	animationSpeeds BeefEyeAnimationSpeeds

	image *ebiten.Image

	astarNodes  []pathfinding.Node
	path        paths.Path
	canPathfind bool
	pathFinding bool
}

func createBeefEye(position Vec2f, game *Game) *BeefEye {
	walkSideSpritesheet := createSpritesheet(newVec2i(0, 23), newVec2i(234, 47), 9, ienemiesSpritesheet)
	dieSpritesheet := createSpritesheet(newVec2i(0, 48), newVec2i(570, 71), 19, ienemiesSpritesheet)
	attackSpritesheet := createSpritesheet(newVec2i(270, 48), newVec2i(570, 71), 10, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	b := &BeefEye{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.2,
		weight:    0.6,

		health:       25,
		maxHealth:    25,
		idle:         true,
		attackRadius: 30,

		shadow: &shadow,

		spritesheet: dieSpritesheet,
		animations: BeefEyeAnimations{
			walkSide: createAnimation(walkSideSpritesheet, ienemiesSpritesheet),
			die:      createAnimation(dieSpritesheet, ienemiesSpritesheet),
			attack:   createAnimation(attackSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: BeefEyeAnimationSpeeds{
			walk:   1.3,
			die:    2.,
			attack: 5.,
		},

		astarNodes:  []pathfinding.Node{},
		canPathfind: true,

		image: ienemiesSpritesheet,
	}

	return b
}

func (b *BeefEye) render(screen *ebiten.Image) {
	if len(b.animation.spritesheet.sprites) == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	// FLIP DECIDER
	flip := newVec2f(-1, 1)
	if b.flipped {
		flip.x = 1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-b.size.x)/2, float64(0-b.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Translate(float64(b.size.x)/2, float64(b.size.y)/2)
	b.subImageRect = image.Rect(
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.y,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(b.position.x), float64(b.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	// Knockback (turning red) render
	b.knockedBack, b.knockedBackTimer = enemyKnockbackRender(op, b.knockedBack, b.knockedBackTimer)

	b.shockwaveHandler.render(screen)

	screen.DrawImage(b.image.SubImage(b.subImageRect).(*ebiten.Image), op) // Draw enemy
}

func (b *BeefEye) update(game *Game) {
	// Start the animation if it's not playing
	if b.idle {
		b.animation = b.animations.walkSide
		b.animation.startForwards()
		b.idle = false
	} else if b.deathExplosion && !b.dying {
		b.dying = true
		b.animation = b.animations.die
		b.animation.startBackwards()
	}

	// Animation speeds and other specialities related to animations
	switch b.animation.id {
	case b.animations.walkSide.id:
		b.animation.update(b.animationSpeeds.walk)
	case b.animations.die.id: // Special case for explosion animation
		b.animation.update(b.animationSpeeds.die)
		if b.animation.finishedFirstPlay {
			b.deathExplosionFinished = true
			b.deathExplosion = false
		}
	case b.animations.attack.id:
		b.animation.update(b.animationSpeeds.attack)
	}

	if !b.dying && !b.attacking {

		// Move enemy
		b.position.x += b.velocity.x
		b.position.y += b.velocity.y

		b.size = newVec2i(
			b.animation.spritesheet.sprites[b.animation.currentFrame].size.x,
			b.animation.spritesheet.sprites[b.animation.currentFrame].size.y,
		)
		endPosition := newVec2i(
			int(b.position.x)+b.size.x,
			int(b.position.y)+b.size.y,
		)

		// Update shadow
		b.shadow.update(b.position, b.size)

		b.subImageRect = image.Rect(int(b.position.x), int(b.position.y), endPosition.x, endPosition.y)
		b.center = newVec2f(b.position.x+float64(b.size.x)/2, b.position.y+float64(b.size.y)/2)
	}

	// Attack

	b.attack(game)
}

func (b *BeefEye) isDead() bool {
	if b.health <= 0 {

		if !b.deathExplosionFinished {
			b.deathExplosion = true
		} else {
			if !b.dead && !b.deathExplosion {
				b.dead = true
			}
			return true
		}
	}
	return false
}

func (b *BeefEye) attack(game *Game) {
	b.shockwaveHandler.update(*b)
	if !b.dying {
		if b.attacking {
			if b.animation.finishedFirstPlay {
				b.attacking = false
				b.idle = true
				b.shockwaveHandler.init(b, 21)
			}
		} else if !b.attacking && isCircularCollision(game.player.getBoundsDynamic(), image.Rect(int(b.center.x-b.attackRadius), int(b.center.y-b.attackRadius), int(b.center.x+b.attackRadius), int(b.center.y+b.attackRadius))) {

			b.animation = b.animations.attack
			b.animation.startBackwards()

			b.attacking = true

		}
	}
}

func (b *BeefEye) getCenter() Vec2f {
	return b.center
}

func (b *BeefEye) getPosition() Vec2f {
	return b.position
}

func (b *BeefEye) getCurrentSubImageRect() image.Rectangle {
	return b.subImageRect
}

func (b *BeefEye) getImage() *ebiten.Image {
	return b.image
}

func (b *BeefEye) getSize() Vec2i {
	return b.size
}

func (b *BeefEye) damage(value int) {
	b.health--
}

func (b *BeefEye) getShadow() Shadow {
	return *b.shadow
}

func (b *BeefEye) getMoveSpeed() float64 {
	return b.moveSpeed
}

func (b *BeefEye) getDying() bool {
	return b.dying
}

func (b *BeefEye) getAttacking() bool {
	return b.attacking
}

func (b *BeefEye) getPath() *paths.Path {
	return &b.path
}

func (b *BeefEye) getCanPathfind() bool {
	return b.canPathfind
}

func (b *BeefEye) getFlipped() bool {
	return b.flipped
}

func (b *BeefEye) getPathfinding() bool {
	return b.pathFinding
}

func (b *BeefEye) setPosition(pos Vec2f) {
	b.position = pos
}

func (b *BeefEye) setFlipped(flipped bool) {
	b.flipped = flipped
}

func (b *BeefEye) setPath(path paths.Path) {
	b.path = path
}

func (b *BeefEye) setCanPathfind(canPathfind bool) {
	b.canPathfind = canPathfind
}

func (b *BeefEye) getWeight() float64 {
	return b.weight
}

func (b *BeefEye) getKnockedBack() bool {
	return b.knockedBack
}

func (b *BeefEye) setKnockedBack(knockedBack bool) {
	b.knockedBack = knockedBack
}
