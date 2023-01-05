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
var camX int = 0
var camY int = 0

//go:export start
func start() {

	rand.Seed(654348654)
	points = []*Point{}
	sticks = []*Stick{}
	w4.PALETTE[0] = 0xfcdeea
	w4.PALETTE[1] = 0x012824
	w4.PALETTE[2] = 0x265935
	w4.PALETTE[3] = 0xff4d6d

	for i := 0.0; i < 8; i++ {
		var y1, y2 float64
		if int(i)%2 == 0 {
			y1 = rand.Float64()*40 - 20
			y2 = rand.Float64() * 40
		} else {
			y1 = rand.Float64() * 40
			y2 = rand.Float64()*40 - 20
		}

		p, s := CreateRope(
			Vector{0, y1 + i*30},
			Vector{320, y2 + i*30},
			14,
		)
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
	player.Update()
	camX, camY = int(player.Position.X)-80, int(player.Position.Y)-80
	for _, s := range sticks {
		s.Draw()
	}
	for _, p := range points {
		p.Draw()
	}
	player.Draw()

	*w4.DRAW_COLORS = 0x23
	w4.Rect(-160-camX, -100-camY, 160, 500)
	w4.Rect(320-camX, -100-camY, 160, 500)
}
