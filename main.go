package main

import (
	"fmt"
	"os"
	"syscall"

	"strings"

	"bytes"
	"os/exec"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
)

func doesServiceExist(serviceName string) (bool, error) {
	// Execute the `sc query state=all` command
	cmd := exec.Command("sc", "query", "state=all")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to execute sc query: %w", err)
	}

	// Read and parse the command's output
	output := out.String()
	if strings.Contains(output, serviceName) {
		return true, nil // Service exists
	}

	return false, nil // Service does not exist
}

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

func isAdmin() bool {
	elevated := windows.GetCurrentProcessToken().IsElevated()
	return elevated
}

func runAsAdmin() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	return err
}
