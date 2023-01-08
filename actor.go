package main

type Actor interface {
	Update() bool
	Draw()
	TakeHit(Actor)
	GetPosition() Vector
}
