package main

import (
	"time"

	"towmer/api"
	"towmer/ecs"
	"towmer/game"
	"towmer/renderer"
)

var w, h, x, y int

func main() {

	err := renderer.Init()
	if err != nil {
		panic(err)
	}
	defer renderer.Close()

	// event loop
	event_queue := make(chan api.Key)
	go func() {
		for {
			event_queue <- renderer.PollEvent()
		}
	}()

	w, h = renderer.Size()

	HudY := 4

	ecs.Bounds(float64(w), float64(h-HudY))

	pool := ecs.NewPool(100)

	game.Spawn(pool, float64(w), float64(h-HudY))

gameloop:
	for {
		select {
		case key := <-event_queue:
			switch key {
			case api.Key_ESC:
				break gameloop
			default:
				ecs.Control(pool, key)
			}
		default:
			ecs.Update(pool, 0.01)
			renderer.Render(pool, float64(w), float64(h-HudY))
			time.Sleep(10 * time.Millisecond)
		}
	}
}
