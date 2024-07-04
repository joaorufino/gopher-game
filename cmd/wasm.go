//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/cv-game/pkg/game"
	"go.uber.org/fx"

	_ "net/http/pprof"
)

func main() {
	c := make(chan struct{})
	log.Println("Starting the CV Game in WebAssembly")

	app := fx.New(
		fx.Provide(
			provideConfiguration,
			provideResourceManager,
			provideInputHandler,
			provideEventManager,
			providePhysicsEngine,
			provideGameMap,
			provideCamera,
			provideSettings,
			provideAudioManager,
			provideParticleSystem,
			provideBackgroundImage,
			providePlayer,
			provideItemManager,
			provideAbilitiesManager,
			fx.Annotate(provideScreenWidth, fx.ResultTags(`name:"screenWidth"`)),
			fx.Annotate(provideScreenHeight, fx.ResultTags(`name:"screenHeight"`)),
			game.NewGame,
		),
		fx.Invoke(startGame),
	)

	go func() {
		app.Run()
	}()

	<-c
}

func startGame(lc fx.Lifecycle, gameInstance *game.Game) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ebiten.SetWindowTitle("CV Game")
			go func() {
				if err := ebiten.RunGame(gameInstance); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

// fetchData is a Go function that interacts with the JavaScript fetchData function.
func fetchData(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		log.Println("Path and callback function are required")
		return nil
	}

	path := args[0].String()
	callback := args[1]
	done := make(chan struct{})

	go func() {
		js.Global().Call("fetchData", path, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer close(done)
			if len(args) < 1 {
				log.Println("Data not provided")
				return nil
			}

			data := []byte(args[0].String())
			jsData := js.Global().Get("Uint8Array").New(len(data))
			js.CopyBytesToJS(jsData, data)
			callback.Invoke(jsData)
			return nil
		}))
	}()

	<-done
	return nil
}
