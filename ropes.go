package main

import (
	"cart/w4"
	"math"
)

type Point struct {
	Position         Vector
	PreviousPosition Vector
	IsLocked         bool
	TimeOffset       uint16
	Sticks           []*Stick
}

func (p *Point) AddStick(stick *Stick) *Point {
	p.Sticks = append(p.Sticks, stick)
	return p
}

func (p *Point) GetMotion() Vector {
	return p.Position.Sub(p.PreviousPosition)
}

func (p *Point) Draw() {
	*w4.DRAW_COLORS = 0x31
	fr := (frame + uint64(p.TimeOffset)) % 60
	if fr > 20 {
		*w4.DRAW_COLORS = 0x34
	}
	if fr > 40 {
		*w4.DRAW_COLORS = 0x32
	}
	w4.Oval(int(p.Position.X)-2-camX, int(p.Position.Y)-2-camY, 4, 4)

}

type Stick struct {
	PointA *Point
	PointB *Point
	Length float64
}

func (s *Stick) GetVector() Vector {
	return s.PointB.Position.Sub(s.PointA.Position)
}

func (s *Stick) GetDistance(point Vector) float64 {
	ab := s.GetVector()
	ap := point.Sub(s.PointA.Position)
	bp := point.Sub(s.PointB.Position)

	abbp := (ab.X*bp.X + ab.Y*bp.Y)
	abap := (ab.X*ap.X + ab.Y*ap.Y)

	if abbp > 0 {
		return bp.Len()
	} else if abap < 0 {
		return ap.Len()
	}
	mod := ab.Len()
	return math.Abs(ab.X*ap.Y-ab.Y*ap.X) / mod
}

func (s *Stick) GetPosition(offset float64) Vector {
	diff := s.GetVector()
	return s.PointA.Position.Sum(diff.MulScalar(offset))
}

func (s *Stick) GetOffset(p Vector) float64 {
	ab := s.GetVector()
	ap := p.Sub(s.PointA.Position)
	projection := ap.ProjectTo(ab)
	dot := ab.Dot(ap)
	result := projection.Len() / ab.Len()
	if dot < 0 {
		result *= -1
	}
	return result
}

func (s *Stick) Draw() {
	*w4.DRAW_COLORS = 0x3
	w4.Line(int(s.PointA.Position.X)-camX, int(s.PointA.Position.Y)-camY,
		int(s.PointB.Position.X)-camX, int(s.PointB.Position.Y)-camY)
}

func Simulate(points []*Point, sticks []*Stick) {
	for _, p := range points {
		if !p.IsLocked {
			positionBeforeUpdate := p.Position
			diff := p.Position.Sub(p.PreviousPosition)
			p.Position = p.Position.Sum(diff)
			p.Position.Y += gravity
			p.PreviousPosition = positionBeforeUpdate
		}
	}

	cycles := 20
	for i := 0; i < cycles; i++ {
		for _, s := range sticks {
			centerX := (s.PointA.Position.X + s.PointB.Position.X) / 2
			centerY := (s.PointA.Position.Y + s.PointB.Position.Y) / 2
			diff := s.PointA.Position.Sub(s.PointB.Position)
			direction := diff.Normalized()
			if !s.PointA.IsLocked {
				s.PointA.Position.X = centerX + direction.X*s.Length/2
				s.PointA.Position.Y = centerY + direction.Y*s.Length/2
			}
			if !s.PointB.IsLocked {
				s.PointB.Position.X = centerX - direction.X*s.Length/2
				s.PointB.Position.Y = centerY - direction.Y*s.Length/2
			}
		}
	}
}

func CreateRope(start Vector, end Vector, divisions int) ([]*Point, []*Stick) {
	var points []*Point = []*Point{}
	var sticks []*Stick = []*Stick{}
	for i := 0; i <= divisions; i++ {
		k := float64(i) / float64(divisions)
		diffX := end.X - start.X
		diffY := end.Y - start.Y
		posX := start.X + diffX*k
		posY := start.Y + diffY*k
		pos := Vector{posX, posY}
		// w4.Trace("Point created at " + pos.String())
		point := Point{
			Position:         pos,
			PreviousPosition: pos,
			IsLocked:         (i == 0 || i == divisions),
			TimeOffset:       uint16(lightIndex * 147),
		}
		lightIndex++
		if i != 0 {
			lastPoint := points[len(points)-1]
			diffX := pos.X - lastPoint.Position.X
			diffY := pos.Y - lastPoint.Position.Y
			len := math.Sqrt(diffX*diffX + diffY*diffY)
			stick := Stick{
				PointA: lastPoint,
				PointB: &point,
				Length: len,
			}
			stick.PointA.AddStick(&stick)
			stick.PointB.AddStick(&stick)
			sticks = append(sticks, &stick)
		}
		points = append(points, &point)
	}
	return points, sticks
}
