package main

import (
	"os"

	"github.com/getlantern/systray"
)

func main() {
	go initApp()
	systray.Run(onReady, onExit)
}

func onReady() {
	//get the icon bytes
	Icon, err := os.ReadFile("assets/icon.ico")
	if err != nil {
		panic(err)
	}
	//set the icon
	systray.SetIcon(Icon)
	systray.SetTitle("CleanDL")
	systray.SetTooltip("Organize your downloads folder")
	mQuit := systray.AddMenuItem("Exit", "Exit the application")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
