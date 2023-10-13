package main

import (
	"embed"
	"go_file_sync/src/initial"
	"go_file_sync/src/tcpserver"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	initial := initial.NewInitial(&app.ctx)
	tcpServer := tcpserver.NewTCPServer(&app.ctx)

	// Default Size 1024, 768 ( Width, Height )
	err := wails.Run(&options.App{
		Title:     "go_file_sync",
		Width:     300,
		Height:    300,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Menu:             app.applicationMenu(),
		Bind: []interface{}{
			app,
			initial,
			tcpServer,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
