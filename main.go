package main

import "cart/w4"

var frame uint64 = 0
var stateUpdate func()

//go:export start
func start() {
	frame = 0
	w4.PALETTE[0] = 0xfcdeea
	w4.PALETTE[1] = 0x012824
	w4.PALETTE[2] = 0x265935
	w4.PALETTE[3] = 0xff4d6d

	menuStart()
}

//go:export update
func update() {
	frame++
	stateUpdate()
}
