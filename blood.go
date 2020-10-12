package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Blood is a drop of blood :)
type Blood struct {
	position Vec2f
	rotation float64

	subImage image.Rectangle

	sprite Sprite
	image  *ebiten.Image
}

func (b *Blood) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Translate center of image to 0, 0 before rotating
	op.GeoM.Translate(0-float64(b.sprite.size.x)/2, 0-float64(b.sprite.size.y)/2)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.x, b.position.y)

	screen.DrawImage(b.image.SubImage(b.subImage).(*ebiten.Image), op)
}

func (b *Blood) update() {

}

// BloodEmitter sprays out blood onto the ground!
type BloodEmitter struct {
	position Vec2f
	rotation float64

	spawnTimer    float64
	spawnTimerMax float64

	blood             []Blood
	bloodSize         Vec2i
	bloodSizeOriginal Vec2i

	image *ebiten.Image
}

func createBloodEmitter(position Vec2f, spawnTimerMax float64, bloodSizeOriginal Vec2i, image *ebiten.Image) BloodEmitter {
	return BloodEmitter{
		position:          position,
		spawnTimer:        0,
		spawnTimerMax:     spawnTimerMax,
		bloodSize:         bloodSizeOriginal,
		bloodSizeOriginal: bloodSizeOriginal,
		image:             image,
	}
}

func (b *BloodEmitter) render(screen *ebiten.Image) {
	for i := 0; i < len(b.blood); i++ {
		b.blood[i].render(screen)
	}
}

func (b *BloodEmitter) update(position Vec2f) {
	if b.spawnTimer > 0 {
		b.spawnTimer--
	} else { // Spawn blood
		b.spawn(position, b.bloodSize)
		b.spawnTimer = b.spawnTimerMax
	}
	for i := 0; i < len(b.blood); i++ {
		b.blood[i].update()
	}
}

// spawn spawns a drop of blood
func (b *BloodEmitter) spawn(position Vec2f, size Vec2i) {
	size.x = rand.Intn(size.x-1) + 1 // Random size based off of the max
	size.y = rand.Intn(size.y-1) + 1 // ^
	min := newVec2i(
		rand.Intn(ibloodSpritesheet.Bounds().Max.Y+size.x)-size.x,
		rand.Intn(ibloodSpritesheet.Bounds().Max.X+size.y)-size.y,
	)

	subImage := image.Rect(min.x, min.y, min.x+size.x, min.y+size.y)
	rotation := rand.Float64() * Pi
	b.blood = append(
		b.blood,
		Blood{
			position: position,
			rotation: rotation,
			subImage: subImage,
			image:    b.image,
		},
	)
}
