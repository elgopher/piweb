package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/piweb"
	"log"
)

func main() {
	pi.SetScreenSize(128, 128)
	pi.SetColor(7)
	pi.RectFill(33, 30, 50, 50)
	pi.SetColor(3)
	picofont.Print("PIWEB", 1, 1)

	pi.Init = func() {
		log.Print("Game initialized")
	}
	piloop.Target().Subscribe(piloop.EventInit, func(event piloop.Event, handler pievent.Handler) {
		//log.Println("EventInit received")
	})

	pi.Update = func() {
		//log.Print("pi.Update")
	}
	piloop.Target().Subscribe(piloop.EventUpdate, func(event piloop.Event, handler pievent.Handler) {
		//log.Println("EventUpdate received")
	})

	pi.Draw = func() {
		//log.Print("pi.Draw")
	}
	piloop.Target().Subscribe(piloop.EventDraw, func(event piloop.Event, handler pievent.Handler) {
		//log.Println("EventDraw received")
	})

	piweb.Run()
}
