package main

import (
	"cart/w4"
	"math/rand"
)

const (
	C int = iota
	Cs
	D
	Ds
	E
	F
	Fs
	G
	Gs
	A
	As
	B
)

var frequencies = []uint{262, 277, 294, 311, 330, 349, 370, 392, 415, 440, 466, 494}

type Music struct {
	KickProb      [8]uint8
	KickLine      [8]bool
	SnareProb     [8]uint8
	SnareLine     [8]bool
	LeadProb      [8]uint8
	LeadLine      [8]uint
	MainNoteIndex int
	Frame         uint64
	CurrentPart   uint8
	CurrentLine   uint8
	Active        bool
	MutedFrames   uint
}

func (m *Music) Start() {
	m.GenLines()
	m.Frame = 0
	m.CurrentLine = 0
	m.CurrentPart = 0
	m.Active = true
}
func (m *Music) Pause() {
	m.Active = false
}
func (m *Music) Resume() {
	m.Active = true
	m.MutedFrames = 0
}

func (m *Music) MuteFor(frames uint) {
	m.MutedFrames = frames
}

var pentaScale = []int{0, 3, 5, 7, 10, 12}

func (m *Music) PentaNote(n int) uint {
	pentaIndex := n % 6
	octave := 1 + n/6
	// w4.Trace("Octave: " + strconv.Itoa(octave))
	noteIndex := (m.MainNoteIndex + pentaScale[pentaIndex]) % 12
	return frequencies[noteIndex] * uint(octave)
}

func (m *Music) Update() {
	if !m.Active {
		return
	}
	if m.MutedFrames > 0 {
		m.MutedFrames--
	}
	if m.Frame%15 == 0 {
		if m.MutedFrames == 0 {
			if m.KickLine[m.CurrentLine] {
				w4.Tone(200|1<<16, 5, 60, w4.TONE_TRIANGLE)
			}
			if m.SnareLine[m.CurrentLine] {
				w4.Tone(200|50<<16, 10<<8, 10, w4.TONE_NOISE)
			}
			if m.LeadLine[m.CurrentLine] != 0 {
				w4.Tone(m.LeadLine[m.CurrentLine], 10|15<<8, 5, w4.TONE_MODE2|w4.TONE_PULSE2)
			}
		}
		m.CurrentLine++
		if m.CurrentLine == 8 {
			m.CurrentLine = 0
			m.CurrentPart++
			if m.CurrentPart == 4 {
				m.CurrentPart = 0
				m.GenLines()
			}
		}
	}
	m.Frame++
}

func (m *Music) GenLines() {
	m.KickLine = m.GenDrumLine(m.KickProb)
	m.SnareLine = m.GenDrumLine(m.SnareProb)
	m.LeadLine = m.GenLeadLine(m.LeadProb)
	m.MainNoteIndex = rand.Intn(5)
}

func (m *Music) GenDrumLine(prob [8]uint8) [8]bool {
	line := [8]bool{}
	for i, p := range prob {
		line[i] = rand.Intn(101) < int(p)
	}
	return line
}

func (m *Music) GenLeadLine(prob [8]uint8) [8]uint {
	line := [8]uint{}
	pattern := [8]int{}
	offset := rand.Intn(6)
	switch rand.Intn(4) {
	case 0:
		pattern = [8]int{0, 1, 2, 3, 0, 1, 2, 3}
	case 1:
		pattern = [8]int{3, 2, 1, 0, 3, 2, 1, 0}
	case 2:
		pattern = [8]int{7, 5, 6, 4, 5, 3, 4, 2}
	case 3:
		pattern = [8]int{0, 2, 1, 3, 2, 4, 3, 5}
	case 4:
		pattern = [8]int{0, 2, 4, 1, 3, 5, 2, 4}
	case 5:
		pattern = [8]int{4, 2, 5, 3, 1, 4, 2, 0}
	}
	for i, p := range prob {
		if rand.Intn(101) < int(p) {
			line[i] = m.PentaNote(pattern[i]+offset) / 2
			continue
		}
		line[i] = 0
	}
	return line
}
