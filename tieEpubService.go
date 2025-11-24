package main

import (
	"context"
	"fmt"
	"tieEpub/tieba"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type TieEpubService struct{}

var ServiceApp *application.App

func (tES *TieEpubService) GetConfig(key string) string {
	return viper.GetString(key)
}

func (tES *TieEpubService) SetConfig(key string, value string) bool {
	viper.Set(key, value)
	err := tieba.SaveConfig()
	return err == nil
}

var ErrorChan chan string
var ProcessChan chan int
var SuccessChan chan int
var ExitChan chan int

func (tES *TieEpubService) StartDownload(url string, minimumWord int, onlyLZ bool, filterLink bool, filterImg bool) {

	tieRequest := tieba.TieRequest{
		Url:         url,
		MinimumWord: minimumWord,
		OnlyLZ:      onlyLZ,
		FilterLink:  filterLink,
		FilterImg:   filterImg,
		ErrorChan:   ErrorChan,
		ProcessChan: ProcessChan,
		SuccessChan: SuccessChan,
	}
	go tieba.GetTie(&tieRequest)
}
func (tES *TieEpubService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	go ChanProcess()
	return nil
}
func (tES *TieEpubService) ServiceShutdown() error {
	ExitChan <- 1
	return nil
}

func (tES *TieEpubService) GetTieData() (bool, tieba.TieContent) {
	data := tieba.GetTieData()
	if data == nil {
		ServiceApp.Event.Emit("downloadError", "小说无内容，请检查下载条件")
		return false, tieba.TieContent{}
	}
	return true, *data
}

func (tES *TieEpubService) SaveEPUB(enableCover bool) (bool, error) {
	data := tieba.GetTieData()
	if data == nil {
		return false, nil
	}
	dialog := application.SaveFileDialog()
	dialog.SetOptions(&application.SaveFileDialogOptions{
		CanCreateDirectories: true,
		Title:                "保存小说",
		Filename:             data.Title + ".epub",
		ButtonText:           "保存",
		Filters: []application.FileFilter{
			{
				DisplayName: "EPUB电子小说文件",
				Pattern:     "*.epub",
			},
		},
	})
	if savePath, err := dialog.PromptForSingleSelection(); err == nil && savePath != "" {
		return true, tieba.SaveEpub(savePath, enableCover)
	}
	return false, nil
}

func (tES *TieEpubService) CreateAiCover(keyword string, chapterIndex int) (imageBase64 string, errMsg string) {
	imageBase64, err := tieba.StartGetAiImg(keyword, chapterIndex)
	if err != nil {
		errMsg = err.Error()
	}
	return
}

func (tES *TieEpubService) OnReady() {
	w, _ := ServiceApp.Window.Get("tieEpub")
	w.Show()
}

func (tES *TieEpubService) DeleteChapter(index int) {
	tieba.DeleteChapter(index)
}

func ChanProcess() {
	ProcessChan = make(chan int)
	ErrorChan = make(chan string)
	SuccessChan = make(chan int)
	ExitChan = make(chan int)
	for {
		select {
		case process := <-ProcessChan:
			ServiceApp.Event.Emit("downloadProcess", process)
		case err := <-ErrorChan:
			fmt.Println(err)
			ServiceApp.Event.Emit("downloadError", err)
		case <-SuccessChan:
			ServiceApp.Event.Emit("downloadSuccess", nil)
		case <-ExitChan:
			fmt.Println("Exiting...")
			return
		}
	}
}
