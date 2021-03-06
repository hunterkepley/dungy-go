package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Shadow is an image that gets places below an object
type Shadow struct {
	position Vec2f
	size     Vec2i

	id int // ID of shadow [unique]

	subImageRect image.Rectangle

	isDrawable bool

	image *ebiten.Image
}

func createShadow(subImageRect image.Rectangle, image *ebiten.Image, id int) Shadow {
	return Shadow{
		subImageRect: subImageRect,
		image:        image,
		id:           id,
		size: newVec2i(
			int(subImageRect.Max.X-subImageRect.Min.X),
			int(subImageRect.Max.Y-subImageRect.Min.Y),
		),
		isDrawable: true,
	}
}

func (s *Shadow) update(itemPosition Vec2f, itemSize Vec2i) {
	s.subImageRect = image.Rect(
		s.subImageRect.Min.X,
		s.subImageRect.Min.Y,
		s.subImageRect.Max.X,
		s.subImageRect.Max.Y,
	)
	s.position = newVec2f(
		itemPosition.x+float64(itemSize.x)/2,
		itemPosition.y+float64(itemSize.y)-float64(s.size.y/4),
	)
}

func (s *Shadow) render(screen *ebiten.Image) {
	if s.isDrawable {
		op := &ebiten.DrawImageOptions{}
		// Translate center of image to 0, 0 before rotating
		op.GeoM.Translate(0-float64(s.size.x)/2, 0-float64(s.size.y)/2)
		//op.GeoM.Rotate(s.rotation)
		op.GeoM.Translate(s.position.x, s.position.y)
		op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

		screen.DrawImage(s.image.SubImage(s.subImageRect).(*ebiten.Image), op)
	}
}

func (s *Shadow) getIndex(shadows []Shadow, id int) int {
	for i := 0; i < len(shadows); i++ {
		if shadows[i].id == id {
			return i
		}
	}
	return -1
}

func removeShadow(slice []*Shadow, id int) []*Shadow {
	s := -1
	for i := 0; i < len(slice); i++ {
		if slice[i].id == id {
			s = i
		}
	}
	return append(slice[:s], slice[s+1:]...)
}

func generateUniqueShadowID(g *Game) int {
	g.shadowID++
	return g.shadowID
}
