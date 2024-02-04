package main

import "github.com/getlantern/systray"

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Systray example")
	systray.SetTooltip("Systray example")
}

func onExit() {
}
