package main

import (
	"cart/w4"
	"math"
	"strconv"
)

type Player struct {
	Position Vector
	Speed    Vector
	//PointGrabbed *Point
	StickGrabbed *Stick
	StickOffset  float64
	GrabTimeout  uint
	Gamepad      *uint8
	GamepadLast  uint8
}

func (p *Player) Update() {
	lastGamepad := *p.Gamepad & (*p.Gamepad ^ p.GamepadLast)
	p.GamepadLast = *p.Gamepad
	if p.GrabTimeout > 0 {
		p.GrabTimeout--
	}
	isJumping := *p.Gamepad&w4.BUTTON_DOWN == 0
	p.Speed.Y = math.Min(4, p.Speed.Y+gravity)
	if p.StickGrabbed != nil {
		p.Speed.Y = 0
		p.Speed.X = 0
		//		p.Position = p.PointGrabbed.Position
		if *p.Gamepad&w4.BUTTON_LEFT != 0 {
			p.Speed.X -= 1
		}
		if *p.Gamepad&w4.BUTTON_RIGHT != 0 {
			p.Speed.X += 1
		}
		if *p.Gamepad&w4.BUTTON_UP != 0 {
			p.Speed.Y -= 1
		}
		if *p.Gamepad&w4.BUTTON_DOWN != 0 {
			p.Speed.Y += 1
		}
		p.MoveOnRope(p.Speed)
		if p.StickOffset < 0 {
			point := p.StickGrabbed.PointA
			if len(point.Sticks) == 1 {
				p.StickOffset = 0
			} else {
				p.StickGrabbed = point.Sticks[0]
				p.StickOffset += 1
			}
		}
		if p.StickOffset > 1 {
			point := p.StickGrabbed.PointB
			if len(point.Sticks) == 1 {
				p.StickOffset = 1
			} else {
				p.StickGrabbed = point.Sticks[1]
				p.StickOffset -= 1
			}
		}
		p.Position = p.StickGrabbed.GetPosition(p.StickOffset)
		if lastGamepad&w4.BUTTON_2 != 0 {
			p.GrabTimeout = 10
			if isJumping {
				p.GrabTimeout = 10
				p.Speed = p.Speed.MulScalar(2)
				if p.Speed.Y <= 0 {
					p.Speed.Y -= 1 * 2
				}
				impulse := p.Speed.MulScalar(-1)
				p.StickGrabbed.PointA.Position.MoveVec(impulse)
				p.StickGrabbed.PointB.Position.MoveVec(impulse)
			}
			p.StickGrabbed = nil
		}
	} else {
		p.Position.Move(p.Speed.X, p.Speed.Y)
		if *p.Gamepad&w4.BUTTON_DOWN == 0 && p.GrabTimeout == 0 {
			distance := math.MaxFloat64
			var selectedPoint *Point
			for _, point := range points {
				diff := p.Position.Sub(point.Position)
				dSquared := diff.LenSquared()
				if dSquared < distance {
					distance = dSquared
					selectedPoint = point
				}
			}
			stickDistance := math.MaxFloat64
			var selectedStick *Stick
			for _, stick := range selectedPoint.Sticks {
				distance := stick.GetDistance(p.Position)
				if distance < stickDistance {
					stickDistance = distance
					selectedStick = stick
				}
			}
			if stickDistance < 4 {
				w4.Trace(strconv.FormatFloat(stickDistance, 'f', 3, 64))
				p.StickGrabbed = selectedStick
				p.StickOffset = selectedStick.GetOffset(p.Position)
				p.StickGrabbed.PointA.Position.MoveVec(p.Speed)
				p.StickGrabbed.PointB.Position.MoveVec(p.Speed)
			}
		}
	}
	p.Position.X = math.Min(math.Max(0, p.Position.X), 320)
}

func (p *Player) MoveOnRope(motion Vector) {
	newPos := p.Position.Sum(motion)
	offset := p.StickGrabbed.GetOffset(newPos)
	p.StickOffset = offset
}

func (p *Player) Draw() {
	*w4.DRAW_COLORS = 0x34
	w4.Rect(int(p.Position.X)-4-camX, int(p.Position.Y)-4-camY, 8, 8)
}
