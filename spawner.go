package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// ImageBit is a piece of an image
type ImageBit struct {
	startPosition Vec2f // Where it flies from
	endPosition   Vec2f // Where it flies to
	position      Vec2f

	rotation  float64
	scale     Vec2f
	imageRect image.Rectangle

	render bool

	image *ebiten.Image
}

func createImageBits(position Vec2f, subImage image.Rectangle, _image *ebiten.Image) []ImageBit {
	var bits []ImageBit

	size := 3

	bitPosition := position

	for i := subImage.Min.X; i < subImage.Max.X; i += size {
		for j := subImage.Min.Y; j < subImage.Max.Y; j += size {

			max := Vec2i{i + size, j + size}

			if max.x > subImage.Max.X {
				max.x = subImage.Max.X
			}
			if max.y > subImage.Max.Y {
				max.y = subImage.Max.Y
			}

			startPosition := Vec2f{float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))}

			// Randomly corrupt
			corruptionChance := rand.Intn(100)
			im := _image
			imageRect := image.Rect(i, j, max.x, max.y)
			if corruptionChance > 85 { // 15% chance; 85+
				im = icorruptionSpritesheet
				x := rand.Intn(ienemiesSpritesheet.Bounds().Dx() - size)
				y := rand.Intn(ienemiesSpritesheet.Bounds().Dy() - size)
				x2 := x + size
				y2 := y + size
				imageRect = image.Rect(x, y, x2, y2)
			}

			bits = append(bits, ImageBit{
				startPosition: startPosition,
				position:      startPosition,
				endPosition:   bitPosition,

				rotation:  0.,
				scale:     Vec2f{1, 1},
				imageRect: imageRect,

				render: false,

				image: im,
			})
			bitPosition.y += float64(size)
		}
		bitPosition.y = position.y
		bitPosition.x += float64(size)
	}
	return bits
}

// EnemySpawner spawns enemies in on an interval and deals with the spawn visuals
type EnemySpawner struct {
	position Vec2f

	bits        []ImageBit
	bitTimer    int
	bitTimerMax int
	currentBit  int

	speed    float64
	speedMax float64

	enemyType EnemyType

	subImage image.Rectangle

	image *ebiten.Image
}

func createEnemySpawner(position Vec2f, enemyType EnemyType, speed float64, subImage image.Rectangle) EnemySpawner {
	image := ienemiesSpritesheet
	return EnemySpawner{
		position: position,

		bits:        createImageBits(position, subImage, image),
		bitTimer:    3,
		bitTimerMax: 3,

		speed:    speed,
		speedMax: speed,

		enemyType: enemyType,

		subImage: subImage,

		image: image,
	}
}

func (e *EnemySpawner) update(g *Game) {
	if e.speed > 0 {
		e.speed--
	} else {
		if e.bitTimer > 0 {
			e.bitTimer--

			for i := 0; i < len(e.bits); i++ {

				moveSpeed := rand.Float64() + 2.

				if e.currentBit < len(e.bits) {
					// Calculate movement using an imaginary vector :)
					dx := e.bits[i].endPosition.x - e.bits[i].position.x
					dy := e.bits[i].endPosition.y - e.bits[i].position.y

					ln := math.Sqrt(dx*dx + dy*dy)

					dx /= ln
					dy /= ln

					// Move towards portal
					e.bits[i].position.x += dx * moveSpeed
					e.bits[i].position.y += dy * moveSpeed

					e.bits[i].rotation = dx * dy // Rotation
				}
			}
		} else {
			if e.currentBit == len(e.bits) {
				// Finished bitting
				e.speed = e.speedMax
				e.currentBit = 0
				// Spawn enemy
				switch e.enemyType {
				case EBeefEye:
					g.enemies = append(g.enemies, createBeefEye(e.position, g))
				case EWorm:
					g.enemies = append(g.enemies, createWorm(e.position, g))
				}

				e.position = g.currentMap.randomPosition()
				e.bits = createImageBits(e.position, e.subImage, e.image)

			} else {
				e.bits[e.currentBit].render = true

				e.currentBit++ // Iterate bit
			}
			e.bitTimer = e.bitTimerMax
		}
	}
}

func (e *EnemySpawner) render(screen *ebiten.Image) {
	for _, b := range e.bits {
		if b.render {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Rotate(b.rotation)
			op.GeoM.Scale(b.scale.x, b.scale.y)
			op.GeoM.Translate(b.position.x, b.position.y)
			op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

			screen.DrawImage(b.image.SubImage(b.imageRect).(*ebiten.Image), op)
		}
	}
}
