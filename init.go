package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func initApp() *cli.App {
	app := &cli.App{
		Name:  "CleanDL",
		Usage: "Organize your downloads folder",
		Action: func(cCtx *cli.Context) error {
			createSettings(patternsPath)
			options := []string{"Organize Downloads Folder", "Edit Pattern Settings", "Exit"}
			flags := flagPointers{AgeThreshold: nil, Destination: nil, DeleteFlag: nil}
			option := choice(DefaultOptionsMessage, options)

			switch option {
			case 1:
				clearScreen()
				organizeFolder()
			case 2:
				clearScreen()
				crudPatterns(flags)
			case 3:
				os.Exit(0)

			default:
				println("Invalid choice. Exiting...")
			}
			return nil
		},
		Commands: []*cli.Command{
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
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new pattern",
				Args:    true,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "pattern", Aliases: []string{"p"}, Required: true, Usage: "The pattern in the form of a regex"},
					&cli.IntFlag{Name: "ageThreshold", Aliases: []string{"t"}, Required: true, DefaultText: "14", Usage: "The age threshold in days"},
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return app
}
