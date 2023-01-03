package main

import (
	"math"
	"strconv"
)

type Vector struct {
	X float64
	Y float64
}

func Vec(x float64, y float64) Vector {
	return Vector{
		x,
		y,
	}
}

func (v *Vector) String() string {
	return "(" + strconv.FormatFloat(v.X, 'f', 3, 64) + "; " + strconv.FormatFloat(v.Y, 'f', 3, 64) + ")"
}

func (v *Vector) LenSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

func (v *Vector) Sum(v2 Vector) Vector {
	return Vector{
		v.X + v2.X,
		v.Y + v2.Y,
	}
}

func (v *Vector) Sub(v2 Vector) Vector {
	return Vector{
		v.X - v2.X,
		v.Y - v2.Y,
	}
}

func (v *Vector) MulScalar(c float64) Vector {
	return Vector{
		v.X * c,
		v.Y * c,
	}
}

func (v *Vector) DivScalar(c float64) Vector {
	return Vector{
		v.X * c,
		v.Y * c,
	}
}

func (v *Vector) Normalized() Vector {
	len := math.Sqrt(v.X*v.X + v.Y*v.Y)
	x := v.X / len
	y := v.Y / len
	return Vector{x, y}
}

func (v *Vector) Dot(v2 Vector) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v *Vector) ProjectTo(v2 Vector) Vector {
	abdot := v.Dot(v2)
	bbdot := v2.Dot(v2)
	return v2.MulScalar(abdot / bbdot)
}

func (v *Vector) Move(x float64, y float64) {
	v.X += x
	v.Y += y
}
func (v *Vector) MoveVec(v2 Vector) {
	v.X += v2.X
	v.Y += v2.Y
}
