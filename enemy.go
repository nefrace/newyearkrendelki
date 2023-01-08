package main

import (
	"cart/w4"
	"math"
)

const enemyWidth = 8
const enemyHeight = 8
const enemyFlags = 1 // BLIT_2BPP
var enemy = [16]byte{0x80, 0x80, 0xaa, 0x80, 0xa6, 0xa4, 0xa9, 0x9a, 0xaa, 0xaa, 0xaa, 0x9a, 0xaa, 0xa8, 0x20, 0x20}

type Enemy struct {
	Position   Vector
	Speed      Vector
	Dead       bool
	ShootTimer uint8
}

func (e *Enemy) GetPosition() Vector {
	return e.Position
}

func (e *Enemy) Update() bool {
	if e.Position.Y > 320 {
		return true
	}
	if e.Dead {
		e.Speed.Y = math.Min(e.Speed.Y+0.2, 4)
	}
	e.Position.MoveVec(e.Speed)
	if e.Position.X < 0 && e.Speed.X < 0 {
		e.Position.X = 0
		e.Speed.X *= -1
	}
	if e.Position.X > 320 && e.Speed.X > 0 {
		e.Position.X = 320
		e.Speed.X *= -1
	}
	if e.ShootTimer > 0 {
		e.ShootTimer--
	} else if e.Position.X > 0 && e.Position.X < 320 {
		e.ShootTimer = 180
		distance := math.MaxFloat64
		direction := Vector{}
		var player *Player
		for _, p := range players {
			if p.Health == 0 {
				continue
			}
			diff := p.Position.Sub(e.Position)
			dis := diff.LenSquared()
			if distance > dis {
				distance = dis
				player = p
				direction = diff.Normalized()
			}

		}
		if player != nil {
			sx := int8(direction.X*5) / 3
			sy := int8(direction.Y*5) / 3
			b := &Bullet{
				Position: e.Position,
				Owner:    e,
				SpeedX:   sx,
				SpeedY:   sy,
			}
			w4.Tone(400|300<<16, 3, 5, w4.TONE_PULSE2)
			actors = append(actors, b)
		}
	}
	return false
}

func (e *Enemy) Draw() {
	*w4.DRAW_COLORS = 0x44
	var f uint = w4.BLIT_2BPP
	if e.Speed.X < 0 {
		f |= w4.BLIT_FLIP_X
	}
	*w4.DRAW_COLORS = 0x430
	w4.Blit(&enemy[0], int(e.Position.X)-camX-4, int(e.Position.Y)-camY-4, 8, 8, f)
}

func (e *Enemy) TakeHit(from Actor) {
	if e.Dead {
		return
	}
	if bullet, ok := from.(*Bullet); ok {
		if _, ok := bullet.Owner.(*Player); !ok {
			return
		}
	}
	enemies--
	kills++
	e.Dead = true
	e.Speed.Y = -3
	w4.Tone(300|200<<16, 5, 20, w4.TONE_PULSE2)
}
