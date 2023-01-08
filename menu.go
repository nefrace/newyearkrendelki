package main

import (
	"cart/w4"
	"strconv"
	"unsafe"
)

var readyPlayers = [4]bool{false, false, false, false}
var pvp = false
var playersCount = 0
var maxKills int = 0
var gamepads = []*uint8{
	w4.GAMEPAD1,
	w4.GAMEPAD2,
	w4.GAMEPAD3,
	w4.GAMEPAD4,
}

var options = []string{"Survival (1P)", "PvP (Netplay only)", "Music: ON"}
var selected = 0
var musicEnabled = true
var lastGamepad uint8

func menuStart() {
	w4.DiskR(unsafe.Pointer(&maxKills), 4)
	playersCount = 0
	readyPlayers = [4]bool{false, false, false, false}
	stateUpdate = menuUpdate
	selected = 0
}

func menuUpdate() {
	frame++
	gpPressed := *w4.GAMEPAD1 & (*w4.GAMEPAD1 ^ lastGamepad)
	lastGamepad = *w4.GAMEPAD1
	if gpPressed&w4.BUTTON_DOWN != 0 {
		selected = (selected + 1) % len(options)
	}
	if gpPressed&w4.BUTTON_UP != 0 {
		selected = (selected - 1) % len(options)
	}
	if gpPressed&w4.BUTTON_1 != 0 {
		switch selected {
		case 0:
			readyPlayers = [4]bool{true, false, false, false}
			playersCount = 1
			pvp = false
			gameStart()
			stateUpdate = gameUpdate
		case 1:
			if playersCount > 1 {
				pvp = true
				gameStart()
				stateUpdate = gameUpdate
			}
		case 2:
			musicEnabled = !musicEnabled
			if musicEnabled {
				options[2] = "Music: ON"
			} else {
				options[2] = "Music: OFF"
			}
		}
	}
	for i, gp := range gamepads {
		if *gp&w4.BUTTON_2 != 0 {
			if !readyPlayers[i] {
				readyPlayers[i] = true
				playersCount++
			}
		}
	}
	w4.Text(strconv.Itoa(maxKills)+" kills", 2, 40)
	w4.Text("CHRISTMAS\nLIGHTS\nMASSACRE", 70, 2)
	w4.Text("Press Z to get ready", 0, 150)
	for i, opt := range options {
		*w4.DRAW_COLORS = 0x3
		if i == selected {
			*w4.DRAW_COLORS = 0x4
		}
		w4.Text(opt, 2, 60+i*9)
		*w4.DRAW_COLORS = 0x2
	}
	for i, v := range readyPlayers {
		if v {
			w4.Text("Player "+strconv.Itoa(i+1)+" is ready", 8, 100+i*8)
		}
	}
}
