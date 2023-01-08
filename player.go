package main

import (
	"cart/w4"
	"math"
	"strconv"
)

type Player struct {
	Index        uint8
	Health       uint8
	KilledBy     Actor
	Position     Vector
	Speed        Vector
	StickGrabbed *Stick
	StickOffset  float64
	GrabTimer    uint8
	ShootTimer   uint8
	ShootDirX    int8
	ShootDirY    int8
	Gamepad      *uint8
	GamepadLast  uint8
}

func (p *Player) GetPosition() Vector {
	return p.Position
}

func (p *Player) Update() bool {
	var dirX int8 = int8(*p.Gamepad&w4.BUTTON_RIGHT>>5) - int8(*p.Gamepad&w4.BUTTON_LEFT>>4)
	var dirY int8 = int8(*p.Gamepad&w4.BUTTON_DOWN>>7) - int8(*p.Gamepad&w4.BUTTON_UP>>6)
	if !(dirX == 0 && dirY == 0) && p.Health > 0 {
		p.ShootDirX = dirX
		p.ShootDirY = dirY
	}
	lastGamepad := *p.Gamepad & (*p.Gamepad ^ p.GamepadLast)
	p.GamepadLast = *p.Gamepad
	if p.GrabTimer > 0 {
		p.GrabTimer--
	}
	if p.ShootTimer > 0 {
		p.ShootTimer--
	} else {
		if *p.Gamepad&w4.BUTTON_1 != 0 && !(p.ShootDirX == 0 && p.ShootDirY == 0) {
			bullet := &Bullet{
				Owner:    p,
				Position: p.Position,
				SpeedX:   p.ShootDirX * 3,
				SpeedY:   p.ShootDirY * 3,
			}
			// w4.Trace(strconv.Itoa(int(p.ShootDirX)))
			// w4.Trace(strconv.Itoa(int(p.ShootDirY)))
			// w4.Trace(p.Position.String())
			w4.Tone(150|50<<16, 10, 100, w4.TONE_MODE1)
			p.ShootTimer = 15
			actors = append(actors, bullet)
		}
	}

	isJumping := *p.Gamepad&w4.BUTTON_DOWN == 0
	p.Speed.Y = math.Min(4, p.Speed.Y+gravity)
	if p.StickGrabbed != nil {
		p.Speed.X = 0
		p.Speed.Y = 0
		if *p.Gamepad&w4.BUTTON_1 == 0 {
			p.Speed.X = float64(dirX)
			p.Speed.Y = float64(dirY)
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
			p.GrabTimer = 10
			if isJumping {
				p.GrabTimer = 10
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
		if p.ShootTimer == 15 {
			impulse := Vector{-float64(dirX), -float64(dirY)}
			if p.StickGrabbed != nil {
				p.StickGrabbed.PointA.Position.MoveVec(impulse)
				p.StickGrabbed.PointB.Position.MoveVec(impulse)
			}

		}
	} else {
		p.Position.Move(p.Speed.X, p.Speed.Y)
		if *p.Gamepad&w4.BUTTON_DOWN == 0 && p.GrabTimer == 0 && p.Health > 0 {
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
				p.StickGrabbed = selectedStick
				p.StickOffset = selectedStick.GetOffset(p.Position)
				p.StickGrabbed.PointA.Position.MoveVec(p.Speed)
				p.StickGrabbed.PointB.Position.MoveVec(p.Speed)
			}
		}
	}
	p.Position.X = math.Min(math.Max(0, p.Position.X), 320)
	if p.Position.Y > 320 && p.Health > 0 {
		p.Health = 0
		p.Death()
	}
	return false
}

func (p *Player) MoveOnRope(motion Vector) {
	newPos := p.Position.Sum(motion)
	offset := p.StickGrabbed.GetOffset(newPos)
	p.StickOffset = offset
}

func (p *Player) Draw() {
	*w4.DRAW_COLORS = 0x34
	drawX := int(p.Position.X) - 4 - camX
	drawY := int(p.Position.Y) - 4 - camY
	w4.Rect(drawX, drawY, 8, 8)

	eyePosX := drawX + 4 + int(p.ShootDirX)
	eyePosY := drawY + 4 - 2 + int(p.ShootDirY)
	w4.Rect(eyePosX-2, eyePosY, 1, 2)
	w4.Rect(eyePosX+1, eyePosY, 1, 2)
	if pvp {
		*w4.DRAW_COLORS = 0x4
		w4.Text("P"+strconv.Itoa(int(p.Index+1)), drawX-3, drawY-10)
	}
}

func (p *Player) TakeHit(from Actor) {
	if p.Health > 0 {
		if p.Health--; p.Health == 0 {
			if bullet, ok := from.(*Bullet); ok {
				p.StickGrabbed = nil
				p.KilledBy = bullet.Owner
				p.Death()
			}
		}
		w4.Tone(250|200<<16, 15, 60, w4.TONE_PULSE1)
	}
}

func (p *Player) Death() {
	p.Speed.Y = -3
	playersAlive--
	music.MuteFor(30)
	w4.Tone(600|200<<16, 15|15<<8, 30, w4.TONE_PULSE2)
}
