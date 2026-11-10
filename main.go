package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"interview-copilot/internal"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	configStore := internal.NewConfigStore()
	meetingStore := internal.NewMeetingStore()
	appService := internal.NewAppService(configStore, meetingStore)

	app := application.New(application.Options{
		Name:        "Interview Copilot",
		Description: "Discreet AI Copilot for Live Interviews & Coding",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	appService.SetApp(app)

	win := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:            "Interview Copilot",
		Width:            780,
		Height:           640,
		AlwaysOnTop:      true,
		Frameless:        true,
		BackgroundType:   application.BackgroundTypeTranslucent,
		BackgroundColour: application.NewRGBA(15, 23, 42, 0),
		URL:              "/",
		Windows: application.WindowsWindow{
			HiddenOnTaskbar:                   true,
			DisableFramelessWindowDecorations: true,
			Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
				1: 1, // CoreWebView2PermissionKindMicrophone -> Allow
			},
		},
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHidden,
		},
	})

	appService.SetWindow(win)

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		go func() {
			time.Sleep(300 * time.Millisecond)
			if hwnd, err := win.NativeWindowHandle(); err == nil && hwnd != 0 {
				if configStore.Get().ScreenProtection {
					internal.SetWindowProtection(hwnd, true)
				}
			}
		}()
	})

	// Register global hotkeys
	internal.StartGlobalHotkeys(func(action string) {
		switch action {
		case "assist":
			go appService.RunFeature("assist", "")
		case "say":
			go appService.RunFeature("say", "")
		case "leetcode":
			go appService.RunFeature("leetcode", "")
		case "toggle_hide":
			app.EmitEvent("window:toggle_hide", nil)
		case "quit":
			app.Quit()
		}
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
