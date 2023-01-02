package main

import (
	"cart/w4"
	"math"
	"strconv"
)

type Point struct {
	Position         Vector
	PreviousPosition Vector
	IsLocked         bool
	TimeOffset       uint64
}

func (p *Point) Draw() {
	*w4.DRAW_COLORS = 0x31
	fr := (frame + p.TimeOffset) % 60
	if fr > 20 {
		*w4.DRAW_COLORS = 0x32
	}
	if fr > 40 {
		*w4.DRAW_COLORS = 0x34
	}
	w4.Oval(int(p.Position.X)-2, int(p.Position.Y)-2, 4, 4)
}

type Stick struct {
	PointA *Point
	PointB *Point
	Length float64
}

func (s *Stick) Draw() {
	*w4.DRAW_COLORS = 0x3
	w4.Line(int(s.PointA.Position.X), int(s.PointA.Position.Y),
		int(s.PointB.Position.X), int(s.PointB.Position.Y))
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
			TimeOffset:       lightIndex * 153,
		}
		lightIndex++
		if i != 0 {
			lastPoint := points[len(points)-1]
			diffX := pos.X - lastPoint.Position.X
			diffY := pos.Y - lastPoint.Position.Y
			len := math.Sqrt(diffX*diffX + diffY*diffY)
			w4.Trace("Length between points is " + strconv.FormatFloat(len, 'f', 3, 64))
			stick := Stick{
				PointA: lastPoint,
				PointB: &point,
				Length: len,
			}
			sticks = append(sticks, &stick)
		}
		points = append(points, &point)
	}
	return points, sticks
}
