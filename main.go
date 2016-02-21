package main

import (
	"time"

	"towmer/api"
	"towmer/ecs"
	"towmer/game"
	"towmer/renderer"
)

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

	pool := ecs.NewPool(100)

	game.Spawn(pool)

gameloop:
	for {
		select {
		case key := <-event_queue:
			switch key {
			case api.Key_ESC:
				break gameloop
			// case api.Key_PANUP, api.Key_PANDOWN, api.Key_PANLEFT, api.Key_PANRIGHT, api.Key_ZOOMIN, api.Key_ZOOMOUT:
			// 	renderer.CamAction(key)
			default:
				ecs.Control(pool, key)
			}
		default:
			ecs.Update(pool, 0.01)
			renderer.Render(pool)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
