package main

import (
	"cart/w4"
	"math/rand"
)

var gravity = 0.2
var points []*Point = []*Point{}
var sticks []*Stick = []*Stick{}
var player *Player
var frame uint64 = 0
var lightIndex uint64 = 0
var camX = 0
var camY = 0

//go:export start
func start() {

	rand.Seed(654654321348654)
	points = []*Point{}
	sticks = []*Stick{}
	w4.PALETTE[0] = 0xfcdeea
	w4.PALETTE[1] = 0x012824
	w4.PALETTE[2] = 0x265935
	w4.PALETTE[3] = 0xff4d6d
	for i := 0; i < 4; i++ {
		p, s := CreateRope(
			Vector{160, rand.Float64()*40 + float64(i*40)},
			Vector{0, rand.Float64()*40 + float64(i*40)},
			10,
		)
		points = append(points, p...)
		sticks = append(sticks, s...)
	}
	for i := 0; i < 3; i++ {
		p, s := CreateRope(
			Vector{float64(i*30 + 30), 0},
			Vector{float64(i*30 + 40), 100},
			10,
		)
		p[len(p)-1].IsLocked = false
		points = append(points, p...)
		sticks = append(sticks, s...)
	}
	player = &Player{
		Position:     Vector{80, 80},
		Speed:        Vector{},
		Gamepad:      w4.GAMEPAD1,
		StickGrabbed: sticks[rand.Intn(len(sticks)-1)],
	}
}

//go:export update
func update() {
	frame += 1
	*w4.DRAW_COLORS = 2
	// w4.Text("Hello from Go!", 10, 10)
	Simulate(points, sticks)
	for _, s := range sticks {
		s.Draw()
	}
	for _, p := range points {
		p.Draw()
	}
	player.Update()
	player.Draw()
}
