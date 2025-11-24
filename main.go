package main

import (
	"embed"
	_ "embed"
	"log"
	"tieEpub/tieba"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {

}
func main() {
	configErr := tieba.InitConfig()

	if configErr != nil {
		log.Fatal(configErr)
	}
	app := application.New(application.Options{
		Name:        "百度贴吧小说下载器",
		Description: "百度贴吧小说下载器",
		Services: []application.Service{
			application.NewService(&TieEpubService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:                "tieEpub",
		Title:               "百度贴吧小说下载器",
		Width:               1024,
		Height:              768,
		MinHeight:           440,
		MinWidth:            800,
		Frameless:           true,
		MaximiseButtonState: application.ButtonDisabled,
		BackgroundType:      application.BackgroundTypeTranslucent,
		Hidden:              true,
		BackgroundColour:    application.NewRGBA(27, 38, 54, 0),
		URL:                 "/",
	})
	ServiceApp = app
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
