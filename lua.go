package main

import (
	lua "github.com/yuin/gopher-lua"
)

// LuaFunction is a struct that stores basic info about a function in a Lua file
type LuaFunction struct {
	name       string // Basically, the filename
	numReturns int    // Number of things the function returns
	isFinished bool   // Is the function done running?
}

func createLuaFunction(name string, numReturns int) LuaFunction {
	isFinished := false
	return LuaFunction{name, numReturns, isFinished}
}

// THE ITEM API:

// ItemAPI is the API for Lua-written items in the game
type ItemAPI struct {
	setPlayerHealth    lua.LGFunction // Player health
	playerHealth       lua.LGFunction
	setPlayerEnergy    lua.LGFunction // Player energy
	playerEnergy       lua.LGFunction
	setPlayerWalkSpeed lua.LGFunction // Player walkspeed
	playerWalkSpeed    lua.LGFunction
	setPlayerRunSpeed  lua.LGFunction // Player runspeed
	playerRunSpeed     lua.LGFunction
	setGunFireSpeed    lua.LGFunction // Player gun firespeed
	gunFireSpeed       lua.LGFunction
	setAccuracy        lua.LGFunction // Player accuracy
	accuracy           lua.LGFunction
}

var itemAPI ItemAPI = ItemAPI{
	setPlayerHealth: func(L *lua.LState) int {
		lv := L.ToInt(1)                 // Get argument
		gameReference.player.health = lv // Set it to player's health
		return 0                         // # of returns/results
	},

	playerHealth: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.health) // Get health
		L.Push(lv)                                     // Push result
		return 1                                       // # of returns/results
	},

	setPlayerEnergy: func(L *lua.LState) int {
		lv := L.ToInt(1)
		gameReference.player.energy = lv
		return 0
	},

	playerEnergy: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.energy)
		L.Push(lv)
		return 1
	},

	setPlayerWalkSpeed: func(L *lua.LState) int {
		lv := L.ToInt(1)
		gameReference.player.walkSpeed = float64(lv)
		return 0
	},

	playerWalkSpeed: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.walkSpeed)
		L.Push(lv)
		return 1
	},

	setPlayerRunSpeed: func(L *lua.LState) int {
		lv := L.ToInt(1)
		gameReference.player.runSpeed = float64(lv)
		return 0
	},

	playerRunSpeed: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.runSpeed)
		L.Push(lv)
		return 1
	},

	setGunFireSpeed: func(L *lua.LState) int {
		lv := L.ToInt(1)
		gameReference.player.gun.fireSpeed = lv
		return 0
	},

	gunFireSpeed: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.gun.fireSpeed)
		L.Push(lv)
		return 1
	},

	setAccuracy: func(L *lua.LState) int {
		lv := L.ToInt(1)
		gameReference.player.accuracy = lv
		return 0
	},

	accuracy: func(L *lua.LState) int {
		lv := lua.LNumber(gameReference.player.accuracy)
		L.Push(lv)
		return 1
	},
}

func initLuaFunctions(L *lua.LState) {
	L.SetGlobal("SetPlayerHealth", L.NewFunction(itemAPI.setPlayerHealth))
	L.SetGlobal("PlayerHealth", L.NewFunction(itemAPI.playerHealth))

	L.SetGlobal("SetPlayerEnergy", L.NewFunction(itemAPI.setPlayerEnergy))
	L.SetGlobal("PlayerEnergy", L.NewFunction(itemAPI.playerEnergy))

	L.SetGlobal("SetPlayerWalkSpeed", L.NewFunction(itemAPI.setPlayerWalkSpeed))
	L.SetGlobal("PlayerWalkSpeed", L.NewFunction(itemAPI.playerWalkSpeed))

	L.SetGlobal("SetPlayerRunSpeed", L.NewFunction(itemAPI.setPlayerRunSpeed))
	L.SetGlobal("PlayerRunSpeed", L.NewFunction(itemAPI.playerRunSpeed))

	L.SetGlobal("SetGunFireSpeed", L.NewFunction(itemAPI.setGunFireSpeed))
	L.SetGlobal("GunFireSpeed", L.NewFunction(itemAPI.gunFireSpeed))

	L.SetGlobal("SetAccuracy", L.NewFunction(itemAPI.setAccuracy))
	L.SetGlobal("Accuracy", L.NewFunction(itemAPI.accuracy))
}
