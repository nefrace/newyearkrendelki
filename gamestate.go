package main

import (
	"cart/w4"
	"math/rand"
	"strconv"
	"unsafe"
)

var gravity = 0.2
var points []*Point = []*Point{}
var sticks []*Stick = []*Stick{}
var actors []Actor = []Actor{}
var enemies uint = 0
var players []*Player
var playersAlive = 0
var kills int = 0
var lightIndex uint64 = 0
var camX int = 0
var camY int = 0
var deathTimer = 180
var music *Music

func gameStart() {
	kills = 0
	deathTimer = 180
	rand.Seed(int64(frame))
	music = &Music{
		KickProb:      [8]uint8{100, 10, 25, 15, 50, 10, 20, 30},
		SnareProb:     [8]uint8{0, 0, 70, 20, 20, 10, 70, 30},
		LeadProb:      [8]uint8{90, 70, 60, 60, 90, 60, 70, 60},
		MainNoteIndex: 1,
	}
	if musicEnabled {
		music.Start()
	}
	points = []*Point{}
	sticks = []*Stick{}
	actors = []Actor{}
	players = []*Player{}

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
	playersAlive = playersCount
	for i, ready := range readyPlayers {
		if !ready {
			continue
		}
		player := &Player{
			Index:        uint8(i),
			Health:       20,
			ShootTimer:   10,
			Position:     Vector{80, 80},
			Speed:        Vector{},
			Gamepad:      gamepads[i],
			StickGrabbed: sticks[rand.Intn(len(sticks)-1)],
		}
		actors = append(actors, player)
		players = append(players, player)
	}

}

func gameUpdate() {
	if (!pvp && playersAlive == 0) || (pvp && playersAlive <= 1) {
		if deathTimer--; deathTimer == 0 {
			if !pvp && kills > maxKills {
				w4.DiskW(unsafe.Pointer(&kills), 4)
			}
			menuStart()
			stateUpdate = menuUpdate
		}
	}
	music.Update()
	frame += 1
	Simulate(points, sticks)
	if frame%120 == 0 && enemies < 10 && !pvp {
		spX := rand.Float64()*2 - 1
		posY := rand.Float64()*200 + 50
		posX := float64(-160)
		if spX < 0 {
			posX = 480
		}

		e := &Enemy{
			Position:   Vector{posX, posY},
			Speed:      Vector{spX, 0},
			ShootTimer: uint8(rand.Uint32() % 120),
		}
		enemies++
		actors = append(actors, e)
	}
	actorsToRemove := []int{}
	for i, actor := range actors {
		if actor.Update() {
			actorsToRemove = append(actorsToRemove, i)
		}
	}
	for i := len(actorsToRemove) - 1; i > 0; i-- {
		actor := actorsToRemove[i]
		actors = append(actors[:actor], actors[actor+1:]...)
	}
	target := players[0]
	isDead := false
	if *w4.NETPLAY&0b100 != 0 {
		playerId := *w4.NETPLAY & 0b11
		target = players[playerId]
		if target.Health == 0 {
			isDead = true
			if plr, ok := target.KilledBy.(*Player); ok {
				target = plr
			}
		}
	}
	camX, camY = int(target.Position.X)-80, int(target.Position.Y)-80
	for _, s := range sticks {
		s.Draw()
	}
	for _, p := range points {
		p.Draw()
	}
	for _, actor := range actors {
		actor.Draw()
	}

	*w4.DRAW_COLORS = 0x23
	w4.Rect(-160-camX, -100-camY, 160, 500)
	w4.Rect(320-camX, -100-camY, 160, 500)
	*w4.DRAW_COLORS = 0x4
	if !pvp {
		w4.Text(strconv.Itoa(kills)+" kills", 1, 1)
	}
	//	w4.Text(strconv.Itoa(playersAlive), 0, 0)
	if !isDead {
		if !pvp || (pvp && playersAlive > 1) {
			w4.Text("health", 57, 151)
			hbar := float32(target.Health) / 20 * 80
			w4.Rect(int(80-hbar), 150, uint(hbar*2), 10)
			*w4.DRAW_COLORS = 0x1
			w4.Text("health", 57, 151)
		} else {
			w4.Text("winner", 57, 151)
		}
	} else {
		w4.Text(" dead ", 57, 151)
	}
}
