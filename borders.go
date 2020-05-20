package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// BorderType is a type for the tyletype enum
type BorderType int

const (
	// CornerBorder ... BORDERTYPE ENUM [1]
	CornerBorder BorderType = iota + 1
	// SideWallBorder ... BORDERTYPE ENUM [2]
	SideWallBorder
	// BottomWallBorder ... BORDERTYPE ENUM [3]
	BottomWallBorder
)

func (b BorderType) String() string {
	return [...]string{"Unknown", "CornerBorder", "SideWallBorder", "BottomWallBorder"}[b]
}

// ^ TILETYPE ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Border is the borders on the sides of the screen
type Border struct {
	position Vec2f

	borderType BorderType

	sprite Sprite
	image  *ebiten.Image
}

func createBorder(position Vec2f, borderType BorderType, image *ebiten.Image) Border {
	// Corner border
	sprite := createSprite(newVec2i(17, 51), newVec2i(27, 61), newVec2i(10, 10), itileSpritesheet)
	switch borderType {
	case (SideWallBorder):
		sprite = createSprite(newVec2i(0, 82), newVec2i(10, 262), newVec2i(10, 180), itileSpritesheet)
		break
	case (BottomWallBorder):
		sprite = createSprite(newVec2i(11, 82), newVec2i(28, 262), newVec2i(17, 180), itileSpritesheet)
		break
	}
	return Border{
		position,

		borderType,

		sprite,
		image,
	}
}

func (b *Border) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.position.x), float64(b.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	borderRect := image.Rect(
		b.sprite.startPosition.x,
		b.sprite.startPosition.y,
		b.sprite.endPosition.x,
		b.sprite.endPosition.y,
	)

	screen.DrawImage(b.image.SubImage(borderRect).(*ebiten.Image), op)
}

func (b *Border) update() {

}

func generateBorders(image *ebiten.Image) []Border {
	corner1 := createBorder(newVec2f(0, 0), CornerBorder, image)
	sideWall1 := createBorder(newVec2f(0, float64(corner1.sprite.size.y)), SideWallBorder, image)
	return []Border{corner1, sideWall1}
}
