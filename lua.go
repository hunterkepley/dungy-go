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

// SetPlayerHealth sets the player's health to a specific value
func SetPlayerHealth(L *lua.LState) int {
	lv := L.ToInt(1)                 // Get argument
	gameReference.player.health = lv // Set it to player's health
	return 0                         // # of returns/results
}

// PlayerHealth gives the player's health to Lua
func PlayerHealth(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.health) // Get health

	L.Push(lv) // Push result

	return 1 // # of returns/results
}

// SetPlayerEnergy sets the player's energy ...
func SetPlayerEnergy(L *lua.LState) int {
	lv := L.ToInt(1)
	gameReference.player.energy = lv

	return 0
}

// PlayerEnergy gives the player's energy ...
func PlayerEnergy(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.energy)

	L.Push(lv)

	return 1
}

// SetPlayerWalkSpeed sets the player's walk speed ...
func SetPlayerWalkSpeed(L *lua.LState) int {
	lv := L.ToInt(1)
	gameReference.player.walkSpeed = float64(lv)
	return 0
}

// PlayerWalkSpeed gives the player's walk speed ...
func PlayerWalkSpeed(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.walkSpeed)
	L.Push(lv)
	return 1
}

// SetPlayerRunSpeed sets the player's run speed ...
func SetPlayerRunSpeed(L *lua.LState) int {
	lv := L.ToInt(1)
	gameReference.player.runSpeed = float64(lv)
	return 0
}

// PlayerRunSpeed gives the player's run speed ...
func PlayerRunSpeed(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.runSpeed)
	L.Push(lv)
	return 1
}

// SetGunFireSpeed sets the player's gun fire speed ...
func SetGunFireSpeed(L *lua.LState) int {
	lv := L.ToInt(1)
	gameReference.player.gun.fireSpeed = lv
	return 0
}

// GunFireSpeed gives the player's gun fire speed ...
func GunFireSpeed(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.gun.fireSpeed)
	L.Push(lv)
	return 1
}

// SetAccuracy sets the player's accuracy ...
func SetAccuracy(L *lua.LState) int {
	lv := L.ToInt(1)
	gameReference.player.accuracy = lv
	return 0
}

// Accuracy gives the player's accuracy ...
func Accuracy(L *lua.LState) int {
	lv := lua.LNumber(gameReference.player.accuracy)
	L.Push(lv)
	return 1
}

// TODO: Make functions to reset the values to defaults

func initLuaFunctions(L *lua.LState) {
	L.SetGlobal("SetPlayerHealth", L.NewFunction(SetPlayerHealth))
	L.SetGlobal("PlayerHealth", L.NewFunction(PlayerHealth))

	L.SetGlobal("SetPlayerEnergy", L.NewFunction(SetPlayerEnergy))
	L.SetGlobal("PlayerEnergy", L.NewFunction(PlayerEnergy))

	L.SetGlobal("SetPlayerWalkSpeed", L.NewFunction(SetPlayerWalkSpeed))
	L.SetGlobal("PlayerWalkSpeed", L.NewFunction(PlayerWalkSpeed))

	L.SetGlobal("SetPlayerRunSpeed", L.NewFunction(SetPlayerRunSpeed))
	L.SetGlobal("PlayerRunSpeed", L.NewFunction(PlayerRunSpeed))

	L.SetGlobal("SetGunFireSpeed", L.NewFunction(SetGunFireSpeed))
	L.SetGlobal("GunFireSpeed", L.NewFunction(GunFireSpeed))

	L.SetGlobal("SetAccuracy", L.NewFunction(SetAccuracy))
	L.SetGlobal("Accuracy", L.NewFunction(Accuracy))
}
