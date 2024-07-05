package main

import (
	"io"
	"log"
	"os"

	"os/exec"
	"syscall"

	"path/filepath"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/windows/registry"
)

func addToStartup(appName, appPath string) error {
	// Open the key for writing
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Run`
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return err
	}
	defer key.Close()

	// Set the value of the registry key
	err = key.SetStringValue(appName, appPath)
	if err != nil {
		return err
	}

	return nil
}

func removeFromStartup(appName string) error {
	// Open the key for writing
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return err
	}
	defer key.Close()

	// Delete the value of the registry key
	err = key.DeleteValue(appName)
	if err != nil {
		return err
	}

	return nil
}

func runDetachedProcess() error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare the command to execute the duplicate process
	cmd := exec.Command(exePath, "schedule", "--no-daemon")

	// On Windows, you can configure some process attributes using cmd.SysProcAttr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // Create in a new process group
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return err
	}

	// No need to wait for the process to finish
	return nil
}

func initApp() *cli.App {
	app := &cli.App{
		Name:  "CleanDL",
		Usage: "Organize your downloads folder",
		Action: func(cCtx *cli.Context) error {
			createSettings(patternsPath)
			options := []string{"Schedule to run daily", "Organize Downloads Folder", "Edit Pattern Settings", "Exit"}
			flags := flagPointers{AgeThreshold: nil, Destination: nil, DeleteFlag: nil}
			option := choice(DefaultOptionsMessage, options)

			switch option {
			case 1:
				clearScreen()
				err := runDetachedProcess()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
				// ScheduleDailyTask()
			case 2:
				clearScreen()
				organizeFolder()
			case 3:
				clearScreen()
				crudPatterns(flags)
			case 4:
				os.Exit(0)
			default:
				println("Invalid choice. Exiting...")
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "startup",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add CleanDL to startup",
						Action: func(cCtx *cli.Context) error {
							currentDir, err := os.Getwd()
							if err != nil {
								log.Fatal("Error getting current directory:", err)
								return err
							}

							exeName := "background_task.exe schedule --no-daemon"

							// Join the current directory with the executable name
							exePath := filepath.Join(currentDir, exeName)

							err = addToStartup("CleanDL", exePath)
							if err != nil {
								log.Fatal(err)
							}
							println("CleanDL added to startup")
							os.Exit(0)
							return nil
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"r"},
						Usage:   "remove CleanDL from startup",
						Action: func(cCtx *cli.Context) error {
							err := removeFromStartup("CleanDL")
							if err != nil {
								log.Fatal(err)
							}
							println("CleanDL removed from startup")
							os.Exit(0)
							return nil
						},
					},
				},
			},
			{
				Name:    "organize",
				Aliases: []string{"o"},
				Usage:   "organize the downloads folder",
				Action: func(cCtx *cli.Context) error {
					organizeFolder()
					return nil
				},
			},
			{
				Name:    "schedule",
				Aliases: []string{"s"},
				Usage:   "schedule the organizer; this runs indefinitely in the background",
				Action: func(cCtx *cli.Context) error {
					if !cCtx.Bool("no-daemon") {
						err := runDetachedProcess()
						if err != nil {
							log.Fatal(err)
						}
						os.Exit(0)
					}
					ScheduleDailyTask()
					return nil
				},
				Flags: []cli.Flag{&cli.BoolFlag{Name: "no-daemon", Aliases: []string{"d"}, Usage: "Don't run the organizer as a daemon", DefaultText: "false"}},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new pattern",
				Args:    true,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "pattern", Aliases: []string{"p"}, Usage: "The pattern in the form of a regex"},
					&cli.IntFlag{Name: "ageThreshold", Aliases: []string{"t"}, DefaultText: "14", Usage: "The age threshold in days"},
					&cli.StringFlag{Name: "destination", Aliases: []string{"m"}, Usage: "The directory to be moved to"},
					&cli.BoolFlag{Name: "deleteFlag", Aliases: []string{"d"}, Usage: "Delete the file"},
				},
				Action: func(cCtx *cli.Context) error {
					// Use a pointer to their `string`, `int` and `bool` to represent their respective types or undefined (nil)
					// These must be used safely by checking if they are nil or not before dereferencing
					var pattern *string = getFlag[string](cCtx, "pattern")
					var ageThreshold *int = getFlag[int](cCtx, "ageThreshold")
					var destination *string = getFlag[string](cCtx, "destination")
					var deleteFlag *bool = getFlag[bool](cCtx, "deleteFlag")
					flags := flagPointers{Pattern: pattern, AgeThreshold: ageThreshold, Destination: destination, DeleteFlag: deleteFlag}
					addPattern(flags)
					return nil
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit a pattern",
				Action: func(cCtx *cli.Context) error {
					editPattern()
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove a pattern",
				Action: func(cCtx *cli.Context) error {
					deletePattern()
					return nil
				},
			},
		},
	}
	defaultHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		defaultHelpPrinter(w, templ, data)
		// as sys tray is blocking, we need to exit the app manually
		os.Exit(0)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return app
}
