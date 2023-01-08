package main

import "cart/w4"

type Bullet struct {
	Position Vector
	SpeedX   int8
	SpeedY   int8
	Owner    Actor
	Dead     bool
}

func (b *Bullet) GetPosition() Vector {
	return b.Position
}

func (b *Bullet) Update() bool {
	if b.Dead {
		return true
	}
	b.Position.Move(float64(b.SpeedX), float64(b.SpeedY))
	if b.Position.X < 0 || b.Position.X > 320 || b.Position.Y < -160 || b.Position.Y > 320 {
		return true
	}
	for _, actor := range actors {
		if actor == Actor(b.Owner) || actor == b {
			continue
		}
		if _, ok := actor.(*Bullet); ok {
			continue
		}
		diff := b.Position.Sub(actor.GetPosition())
		if diff.LenSquared() < 60 {
			actor.TakeHit(b)
			// b.TakeHit(b)
			b.Dead = true
			return true
		}
	}
	return false
}

func (b *Bullet) Draw() {
	*w4.DRAW_COLORS = 0x41
	w4.Oval(int(b.Position.X)-3-camX, int(b.Position.Y)-3-camY, 6, 6)
}

func (b *Bullet) TakeHit(from Actor) {
	w4.Tone(800, 10<<8, 20, w4.TONE_MODE2)
	b.Dead = true
}
