package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"

	lua "github.com/yuin/gopher-lua"
)

var (
	listOfAllItemNames []string
)

func initListOfAllItemNames() {
	root := "./Assets/Scripts/Items"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".lua") {
			listOfAllItemNames = append(listOfAllItemNames, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

// Item is an item that contains a component that can be picked up by the Player
type Item struct {
	position Vec2f

	size Vec2i

	fileName  string // File name of item
	functions []LuaFunction

	destroy bool // Should we destroy this item?

	sprite Sprite
	image  *ebiten.Image
}

func createItem(position Vec2f, image *ebiten.Image) Item {

	// Decide random item
	rand.Seed(time.Now().UnixNano())
	randNumber := rand.Intn(len(listOfAllItemNames))

	return Item{
		position: position,

		fileName: listOfAllItemNames[randNumber],

		destroy: false,
		image:   image,
	}
}

func (i *Item) init() {
	var parameters []lua.LValue
	var returns []lua.LValue

	// Get state from Lua files
	// Set up Lua State and load files
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(i.fileName); err != nil {
		panic(err)
	}
	// Call function
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("getInformation"), // function name
		NRet:    4,                             // number of returns
		Protect: true,                          // return err or panic
	}, parameters...); err != nil { // ... allows for any # of parameters
		panic(err)
	}

	// Retrieve all returns
	// ALL RETURNS: 0:name, 1:functionName, 2:numReturns, 3:imageBounds
	for i := 0; i < 4; i++ {
		// Get the returned value from the stack and cast it to a lua lstring
		if str, ok := L.Get(-1).(lua.LValue); ok {
			returns = append(returns, str)
		}

		// Pop return value off stack
		L.Pop(1)
	}

	functionNumReturns, err := strconv.Atoi(returns[2].String())
	if err != nil {
		fmt.Println("Error Assets/Scripts/Items/", returns[0], ".lua -- numReturns is not an int/does not exist")
		log.Fatal(err)
	}
	functionName := returns[1].String()
	i.functions = append(
		i.functions,
		createLuaFunction(
			functionName,
			functionNumReturns,
		),
	)

	spriteSizeLua := returns[3].String()
	spriteSizeRaw := [4]int{}

	var temp strings.Builder
	numNeeded := 3
	totalNeeded := 3
	for i := 0; i < len(spriteSizeLua); i++ {
		if spriteSizeLua[i] == ';' {
			num, err := strconv.Atoi(temp.String())

			if err != nil {
				log.Fatal(err)
			}
			spriteSizeRaw[int(math.Abs(float64(totalNeeded-numNeeded)))] = num // Put it in the Raw list!

			temp.Reset() // Reset string builder!

			numNeeded-- // We need one less
		} else if _, err := strconv.Atoi(string(spriteSizeLua[i])); err == nil { // If not a number
			temp.WriteString(string(spriteSizeLua[i])) // Write to the world builder
		}
	}

	spriteMin := newVec2i(
		spriteSizeRaw[0],
		spriteSizeRaw[1],
	)
	spriteMax := newVec2i(
		spriteSizeRaw[2],
		spriteSizeRaw[3],
	)
	i.sprite = createSprite(
		spriteMin,
		spriteMax,
		newVec2i(spriteMax.x-spriteMin.x, spriteMax.y-spriteMin.y),
		i.image,
	)
}

func (i *Item) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(i.position.x, i.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	subImageRect := image.Rect(
		i.sprite.startPosition.x,
		i.sprite.startPosition.y,
		i.sprite.endPosition.x,
		i.sprite.endPosition.y,
	)
	screen.DrawImage(i.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (i *Item) update(player *Player) {

}

func (i *Item) runLuaFunction(function LuaFunction, c chan []lua.LValue) { // Run a lua function!
	// Get state from Lua file
	// Set up Lua State and load file
	L := lua.NewState()
	defer L.Close()
	initLuaFunctions(L)
	if err := L.DoFile(i.fileName); err != nil {
		panic(err)
	}
	// Call function
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(function.name), // function name
		NRet:    function.numReturns,        // number of returns
		Protect: true,                       // return err or panic
	}); err != nil { // this allows any # of parameters without generics :)
		panic(err)
	}

	var returns []lua.LValue

	// Get the returned value from the stack
	// This value is true if the function stops instantly [healthPack] or if it continues forever
	if str, ok := L.Get(-1).(lua.LValue); ok {
		returns = append(returns, str)
	}

	// Pop return value off stack
	L.Pop(1)

	c <- returns
}

func (i *Item) getBounds() image.Rectangle {
	return image.Rect(
		int(i.position.x),
		int(i.position.y),
		int(i.position.x)+i.size.x,
		int(i.position.y)+i.size.y,
	)
}

func updateItems(game *Game) {
	for i := 0; i < len(game.items); i++ {
		if isAABBCollision(game.items[i].getBounds(), game.player.getBoundsStatic()) {

			game.items[i].destroy = true
			if game.items[i].destroy {
				game.player.items = append(game.player.items, game.items[i])
				game.items = removeItem(game.items, i)
			}
			break
		}
	}
}

func renderItems(game *Game, screen *ebiten.Image) {
	for i := 0; i < len(game.items); i++ {
		game.items[i].render(screen)
	}
}

func removeItem(slice []Item, i int) []Item {
	return append(slice[:i], slice[i+1:]...)
}
