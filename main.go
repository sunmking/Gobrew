package main

import (
	"embed"
	"log"

	"changeme/services"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("brew-output")
	application.RegisterEvent[string]("brew-complete")
}

func main() {
	brewService := services.NewBrewService(nil)
	tapService := services.NewTapService(nil)
	serviceManager := services.NewServiceManager(nil)
	bundleService := services.NewBundleService(nil)

	app := application.New(application.Options{
		Name:        "Gobrew",
		Description: "A Homebrew GUI Client for macOS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(brewService),
			application.NewService(tapService),
			application.NewService(serviceManager),
			application.NewService(bundleService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	brewService.SetApp(app)
	tapService.SetApp(app)
	serviceManager.SetApp(app)
	bundleService.SetApp(app)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Gobrew",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(246, 246, 246),
		URL:              "/",
		Width:            1100,
		Height:           700,
		MinWidth:         800,
		MinHeight:        500,
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
