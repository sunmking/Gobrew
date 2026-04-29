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
	configService := services.NewConfigService()
	app := application.New(application.Options{
		Name:        "Gobrew",
		Description: "A Homebrew GUI Client for macOS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(brewService),
			application.NewService(tapService),
			application.NewService(serviceManager),
			application.NewService(bundleService),
			application.NewService(configService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyRegular,
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	brewService.SetApp(app)
	tapService.SetApp(app)
	serviceManager.SetApp(app)
	bundleService.SetApp(app)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
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

	var tray *application.SystemTray
	if !services.IsDevRuntime() {
		tray = app.SystemTray.New()
		tray.SetLabel("Gobrew")
		tray.SetTooltip("Gobrew")
		menu := app.NewMenu()
		menu.Add("Show Gobrew").OnClick(func(*application.Context) {
			application.InvokeAsync(func() {
				window.UnMinimise()
				window.Restore()
				window.Show()
				window.Focus()
			})
		})
		menu.AddSeparator()
		menu.Add("Quit").OnClick(func(*application.Context) {
			app.Quit()
		})
		tray.SetMenu(menu)
		tray.Run()
	}
	configService.SetMenuBarApplier(func(enabled bool) {
		if tray == nil {
			return
		}
		if enabled {
			tray.Show()
			return
		}
		tray.Hide()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
