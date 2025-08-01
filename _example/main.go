package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/piweb"
	"log"
)

func main() {
	pi.SetScreenSize(128, 128)
	pi.SetColor(7)
	pi.RectFill(33, 30, 50, 50)
	pi.SetColor(1)
	picofont.Print("piweb2", 0, 0)

	surface := pi.NewSurface[int](3, 3)
	surface.SetAll(1, 2, 3, 4, 5, 6, 7, 8, 9)
	log.Print(surface.Get(0, 1))

	piweb.Run()
}
