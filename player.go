package main

import "cart/w4"

type Player struct {
	Position     Vector
	Speed        Vector
	PointGrabbed *Point
	GrabTimeout  uint
	Offset       float64
	Gamepad      *uint8
}

func (p *Player) Update() {
	p.Speed.X = 0
	if p.GrabTimeout > 0 {
		p.GrabTimeout--
	}
	if *p.Gamepad&w4.BUTTON_LEFT != 0 {
		p.Speed.X -= 1
	}
	if *p.Gamepad&w4.BUTTON_RIGHT != 0 {
		p.Speed.X += 1
	}
	isJumping := *p.Gamepad&w4.BUTTON_DOWN == 0
	p.Speed.Y += gravity
	if p.PointGrabbed != nil {
		p.Speed.Y = 0
		p.Position = p.PointGrabbed.Position
		if *p.Gamepad&w4.BUTTON_2 != 0 {
			p.GrabTimeout = 10
			if isJumping {
				p.GrabTimeout = 5
				p.Speed.Y = -4.5
			}
			p.PointGrabbed = nil
		}
	} else {
		p.Position.Move(p.Speed.X, p.Speed.Y)
		if *p.Gamepad&w4.BUTTON_DOWN == 0 && p.GrabTimeout == 0 {
			for _, point := range points {
				diff := p.Position.Sub(point.Position)
				if diff.LenSquared() < 25 {
					// nearPoints = append(nearPoints, p)
					*w4.DRAW_COLORS = 0x44
					w4.Rect(int(point.Position.X), int(point.Position.Y), 3, 3)
					p.PointGrabbed = point
					point.PreviousPosition.Move(-p.Speed.X*2, -p.Speed.Y*2)
					break
				}
			}
		}
	}
}

func (p *Player) Draw() {
	*w4.DRAW_COLORS = 0x34
	w4.Rect(int(p.Position.X)-4, int(p.Position.Y)-4, 8, 8)
}
