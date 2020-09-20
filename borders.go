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
	rotation float64
	flipped  Vec2b

	borderType BorderType

	sprite Sprite
	image  *ebiten.Image
}

func createBorder(position Vec2f, rotation float64, flipped Vec2b, borderType BorderType, image *ebiten.Image) Border {
	rotationAsRadians := rotation * (Pi / 180)
	// Corner border
	sprite := createSprite(newVec2i(17, 51), newVec2i(34, 61), newVec2i(17, 10), itileSpritesheet)
	switch borderType {
	case (SideWallBorder):
		sprite = createSprite(newVec2i(0, 82), newVec2i(17, 328), newVec2i(17, 246), itileSpritesheet)
		break
	case (BottomWallBorder):
		sprite = createSprite(newVec2i(0, 329), newVec2i(228, 346), newVec2i(228, 17), itileSpritesheet)
		break
	}
	return Border{
		position,
		rotationAsRadians,
		flipped,

		borderType,

		sprite,
		image,
	}
}

func (b *Border) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// FLIP DECIDER
	flip := newVec2f(1, 1)
	if b.flipped.x {
		flip.x = -1
	}
	if b.flipped.y {
		flip.y = 1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-b.sprite.size.x)/2, float64(0-b.sprite.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(float64(b.sprite.size.x)/2, float64(b.sprite.size.y)/2)
	// POSITION
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
	corner1 := createBorder(
		newVec2f(0, 0),
		0,
		newVec2b(false, false),
		CornerBorder,
		image,
	)
	corner2 := createBorder(
		newVec2f(screenWidth-float64(corner1.sprite.size.x), 0),
		0,
		newVec2b(true, false),
		CornerBorder,
		image,
	)
	sideWall1 := createBorder(
		newVec2f(0, float64(corner1.sprite.size.y)),
		0,
		newVec2b(true, false),
		SideWallBorder,
		image,
	)
	sideWall2 := createBorder(
		newVec2f(screenWidth-float64(sideWall1.sprite.size.x), float64(corner1.sprite.size.y)),
		0,
		newVec2b(false, false),
		SideWallBorder,
		image,
	)
	bottomWall1 := createBorder(
		newVec2f(0, screenHeight),
		0,
		newVec2b(false, false),
		BottomWallBorder,
		image,
	)
	bottomWall1.position.y -= float64(bottomWall1.sprite.size.y)
	bottomWall2 := createBorder(
		newVec2f(float64(bottomWall1.sprite.endPosition.x)-1, screenHeight-float64(bottomWall1.sprite.size.y)),
		0,
		newVec2b(true, false),
		BottomWallBorder,
		image,
	)
	return []Border{
		corner1,
		corner2,
		sideWall1,
		sideWall2,
		bottomWall1,
		bottomWall2,
	}
}
