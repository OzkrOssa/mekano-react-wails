package main

import (
	"embed"
	"mekano-react-wails/backend/mekano"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "mekano-react-wails",
		Width:  600,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnDomReady: app.onDomReady,
		Bind: []interface{}{
			app,
			&mekano.PaymentStatistics{},
			&mekano.BillingStatistics{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
