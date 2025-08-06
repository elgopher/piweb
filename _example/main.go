package main

import (
	"log"
	"math/rand"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/pimouse"
	"github.com/elgopher/pi/pipad"
	"github.com/elgopher/pi/piscope"
	"github.com/elgopher/piweb"
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
		pi.Circ(rand.Intn(128), rand.Intn(128), rand.Intn(15))
		pi.SetColor(pi.Color(rand.Intn(33)))
		//log.Print("pi.Draw")
	}
	piloop.Target().Subscribe(piloop.EventDraw, func(event piloop.Event, handler pievent.Handler) {
		//log.Println("EventDraw received")
	})

	pikey.Target().SubscribeAll(func(event pikey.Event, handler pievent.Handler) {
		//log.Println(event)
	})

	pipad.ConnectionTarget().SubscribeAll(func(event pipad.EventConnection, handler pievent.Handler) {
		//log.Println(event)
	})

	pipad.ButtonTarget().SubscribeAll(func(event pipad.EventButton, handler pievent.Handler) {
		//log.Println(event)
	})

	pimouse.MoveTarget().SubscribeAll(func(event pimouse.EventMove, handler pievent.Handler) {
		//log.Println(event)
	})

	pimouse.ButtonTarget().SubscribeAll(func(event pimouse.EventButton, handler pievent.Handler) {
		//log.Println(event)
	})

	piscope.Start()

	piweb.Run()
}
