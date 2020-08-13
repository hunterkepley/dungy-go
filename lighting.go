package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
)

// LightImages is a struct of Rectangles of all the light images
type LightImages struct {
	playerLight          image.Rectangle
	bulletLight          image.Rectangle
	bulletExplosionLight image.Rectangle
}

// this returns a LightImages struct to init the one in the Game struct
func initLightImages() LightImages {
	return LightImages{
		playerLight:          image.Rect(0, 0, 88, 81),
		bulletLight:          image.Rect(88, 0, 139, 50),
		bulletExplosionLight: image.Rect(140, 0, 180, 38),
	}
}

// LightBackground is the "background" [the darkness color] of light
type LightBackground struct {
	position Vec2f
	image    *ebiten.Image // the "background" light image
}

func createLightBackground() LightBackground {
	return LightBackground{
		position: newVec2f(0, 0),
		image:    ilightingBackground,
	}
}

func (b *LightBackground) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(b.image, op)
}

func (b *LightBackground) update() {
}

// Light is a circular image with a position that masks with the LightBackground
type Light struct {
	position Vec2f

	id       int
	rotation float64

	subImage image.Rectangle
	image    *ebiten.Image
}

func createLight(lightRect image.Rectangle, id int, rotation float64) Light {
	return Light{
		subImage: lightRect,
		image:    ilightingSpritesheet,
		id:       id,
		rotation: rotation,
	}

}

func (l *Light) update(center Vec2f, rotation float64) {
	size := newVec2i(l.subImage.Max.X-l.subImage.Min.X, l.subImage.Max.Y-l.subImage.Min.Y)
	l.position = newVec2f(center.x-float64(size.x/2), center.y-float64(size.y/2))
}

// LightHandler controls all lights in the game
type LightHandler struct {
	bg LightBackground
	fg LightBackground

	lights []Light

	lightImages LightImages // All light subimages

	lightID int // Keeps track of light ID's, assigns a unique # to each light

	maskedFgImage *ebiten.Image
}

// this returns a lighthandler to init the one in the Game struct
func initLightHandler() LightHandler {
	maskedFg, _ := ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	return LightHandler{
		bg:            createLightBackground(),
		fg:            createLightBackground(),
		lightImages:   initLightImages(),
		lightID:       0,
		maskedFgImage: maskedFg,
	}
}

/* The lights are a bit complicated so here's an explanation for future reference:
 *
 * Lights have 3 key parts;
 * 1) the masked foreground image [what the lights draw onto],
 * 2) the light background [the default darkness color],
 * 3) the actual background [the screen, everything else drawn],
 *
 * First, we use CompositeModeSourceOver to draw the light images
 * on top of each other. This overlaps like how if you shine two
 * lights in the same spot, you can still see behind both, but
 * the projected light is significantly brighter and harder to see behind.
 *
 * Next, once they are drawn onto the masked foreground image, we change the
 * composite mode to CompositeModeSourceIn which basically takes whatever
 * image we choose and does the opposite to Over. This is what draws the
 * screen under the lights and allows them to make the screen/game brighter.
 *
 * Lastly, we draw the background light image then the masked foreground image
 * onto the screen, this is just the final step and actually displays everything.
 *
 * Overall, not too complicated, dont change anything unless you know what
 * you're doing as it's a little volatile.
 *
 */
func (h *LightHandler) render(screen *ebiten.Image) {
	h.maskedFgImage.Clear()

	// Add lights
	for i := 0; i < len(h.lights); i++ {
		op := &ebiten.DrawImageOptions{}
		op.CompositeMode = ebiten.CompositeModeSourceOver

		lightSize := newVec2i(h.lights[i].subImage.Max.X-h.lights[i].subImage.Min.X, h.lights[i].subImage.Max.Y-h.lights[i].subImage.Min.Y)
		// Rotate light
		op.GeoM.Translate(0-float64(lightSize.x)/2, 0-float64(lightSize.y)/2)
		//op.GeoM.Rotate(h.lights[i].rotation)
		//op.GeoM.Rotate(h.lights[i].rotation)
		op.GeoM.Translate(float64(lightSize.x)/2, float64(lightSize.y)/2)

		// Move light
		op.GeoM.Translate(float64(h.lights[i].position.x), float64(h.lights[i].position.y))
		h.maskedFgImage.DrawImage(h.lights[i].image.SubImage(h.lights[i].subImage).(*ebiten.Image), op)
		// Draw lights
		op = &ebiten.DrawImageOptions{}
		op.CompositeMode = ebiten.CompositeModeSourceIn
		h.maskedFgImage.DrawImage(screen, op)
	}
	// Draw background image with the light image [mask]
	screen.DrawImage(h.bg.image, &ebiten.DrawImageOptions{})
	screen.DrawImage(h.maskedFgImage, &ebiten.DrawImageOptions{})

}

func (h *LightHandler) update() {

}

func (h *LightHandler) addLight(subImage image.Rectangle, rotation float64) int {
	id := h.generateUniqueLightID()
	h.lights = append(h.lights, createLight(subImage, id, rotation))
	return id
}

func removeLight(slice []Light, id int) []Light { // Removes a light given the ID
	l := -1
	for i := 0; i < len(slice); i++ {
		if slice[i].id == id {
			l = i
		}
	}
	fmt.Print("\nRemoving light with ID ", id)
	if l != -1 {
		return append(slice[:l], slice[l+1:]...)
	}
	return slice
}

func (h *LightHandler) generateUniqueLightID() int { // Generates a new ID for a light
	h.lightID++
	return h.lightID
}

func (h *LightHandler) getLightIndex(id int) int { // Get the light index from its ID
	for i := 0; i < len(h.lights); i++ {
		if h.lights[i].id == id {
			return i
		}
	}
	return -1
}
