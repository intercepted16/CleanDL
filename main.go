package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"github.com/urfave/cli/v2"
)

func getFlag[T any](cCtx *cli.Context, flagName string) *T {
	var pattern *T
	if cCtx.IsSet(flagName) {
		patternValue := cCtx.String(flagName)
		pattern = (*T)(unsafe.Pointer(&patternValue))
	}
	return pattern
}

func getUserFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}

func input[T any](prompt string, parseFunc func(string) (T, error)) T {
	var result T
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)
		value, err := parseFunc(input)
		if err == nil {
			result = value
			break
		}
		fmt.Println("Invalid input, please try again.")
	}
	return result
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func writePatternsToFile(patterns regexPatterns) {
	jsonData, err := json.Marshal(regexPatternsJSON{Patterns: patterns})
	if err != nil {
		panic(err)
	}
	file, err := os.Create(patternsPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.Write(jsonData); err != nil {
		panic(err)
	}
}

func getDownloadsFolder() string {
	return filepath.Join(getUserFolder(), "Downloads")
}

type regexInfo struct {
	AgeThreshold int
	Destination  string
	DeleteFlag   bool
}

type regexPatternsJSON struct {
	Patterns regexPatterns
}

type regexPatterns map[string]regexInfo

var patternsPath = filepath.Join(getUserFolder(), "AppData", "Local", "CleanDL", "patterns.json")

func createSettings(path string) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	// if the file doesn't exist, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Serialize the map to JSON
		jsonData, err := json.Marshal(regexPatternsJSON{Patterns: regexPatterns{}})
		if err != nil {
			panic(err) // Consider more graceful error handling
		}

		// Write the JSON data to the file
		if _, err := file.Write(jsonData); err != nil {
			panic(err) // Consider more graceful error handling
		}
	}
}

func getSettings(path string) regexPatterns {
	settingsFile, err := os.Open(path)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", path)
	// defer the closing of our jsonFile so that we can parse it later on
	defer settingsFile.Close()
	byteValue, _ := io.ReadAll(settingsFile)
	// we initialize our custom regex array
	var regexPatternsJSON regexPatternsJSON

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'regexPatternsJSON' which we defined above
	json.Unmarshal(byteValue, &regexPatternsJSON)
	var regexPatterns regexPatterns = regexPatternsJSON.Patterns
	return regexPatterns
}

// Function to move or delete files based on age, type, and regex or simple string
func processFiles(patterns regexPatterns, downloadsFolder string) {
	files, err := os.ReadDir(downloadsFolder)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()

	for _, file := range files {
		filePath := filepath.Join(downloadsFolder, file.Name())
		fileInfo, err := file.Info()
		if err != nil {
			panic(err)
		}
		fileAgeDays := int(currentTime.Sub(fileInfo.ModTime()).Hours() / 24)

		for pattern, info := range patterns {
			matched, err := regexp.MatchString(pattern, file.Name())
			if err != nil {
				panic(err)
			}

			if matched {
				if fileAgeDays > info.AgeThreshold {
					fmt.Printf("File age: %d\n", fileAgeDays)
					fmt.Printf("Age threshold: %d\n", info.AgeThreshold)
					if info.DeleteFlag {
						os.Remove(filePath) // Delete the file
						fmt.Printf("Deleted: %s\n", filePath)
					} else if info.Destination != "" {
						os.Rename(filePath, filepath.Join(info.Destination, file.Name())) // Move the file
						fmt.Printf("Moved: %s to %s\n", filePath, info.Destination)
					}
					break // Exit the loop after processing
				}
			}
		}
	}
}

func main() {
	app := &cli.App{
		Name:  "CleanDL",
		Usage: "Organize your downloads folder",
		Action: func(cCtx *cli.Context) error {
			createSettings(patternsPath)
			options := []string{"Organize Downloads Folder", "Edit Pattern Settings", "Exit"}
			println("Choose an option:\n")
			for i := 0; i < len(options); i++ {
				fmt.Printf("%d. %s\n", i+1, options[i])
			}
			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				clearScreen()
				organizeFolder()
			case 2:
				clearScreen()
				editSettings(cCtx.Args())
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
					var pattern *string = getFlag[string](cCtx, "pattern")
					var ageThreshold *int = getFlag[int](cCtx, "ageThreshold")
					var destination *string = getFlag[string](cCtx, "destination")
					var deleteFlag *bool = getFlag[bool](cCtx, "deleteFlag")
					// Safely use the pointers by checking if they are not nil before dereferencing
					if pattern != nil {
						println("Pattern:", *pattern)
					} else {
						println("Pattern not provided")
					}

					if ageThreshold != nil {
						println("Age threshold:", *ageThreshold)
					} else {
						println("Age threshold not provided")
					}

					if destination != nil {
						println("Destination:", *destination)
					} else {
						println("Destination not provided")
					}

					if deleteFlag != nil {
						println("Delete flag:", *deleteFlag)
					} else {
						println("Delete flag not provided")
					}
					addFileType(cCtx.Args())
					return nil
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit a pattern",
				Action: func(cCtx *cli.Context) error {
					editFileType()
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove a pattern",
				Action: func(cCtx *cli.Context) error {
					deleteFileType()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func organizeFolder() {
	patterns := getSettings(patternsPath)
	downloadsFolder := getDownloadsFolder()
	print("Processing files in: ", downloadsFolder, "\n")
	processFiles(patterns, downloadsFolder)
	print("\nDone!", "\n")
}

func editSettings(args cli.Args) {
	options := []string{"Add Pattern", "Edit Pattern", "Delete Pattern", "Exit"}
	println("Choose an option:\n")
	for i := 0; i < len(options); i++ {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		addFileType(args)
	case 2:
		editFileType()
	case 3:
		deleteFileType()
	case 4:
		clearScreen()
		main()
	default:
		println("Invalid choice. Exiting...")
	}
}

func addFileType(args cli.Args) {
	println("args", args.Present())
	patterns := getSettings(patternsPath)

	pattern := input("Enter the pattern (regex): ", func(input string) (string, error) {
		return input, nil // No conversion needed for string
	})

	ageThreshold := input("Enter the age threshold (in days): ", strconv.Atoi)

	destination := input("Enter the destination folder: ", func(input string) (string, error) {
		return input, nil // No conversion needed for string
	})

	deleteFlag := input("Delete the file? (true/false): ", strconv.ParseBool)
	delete(patterns, pattern)
	patterns[pattern] = regexInfo{AgeThreshold: ageThreshold, Destination: destination, DeleteFlag: deleteFlag}
	writePatternsToFile(patterns)
}

func editFileType() {
	patterns := getSettings(patternsPath)
	println("Choose a Pattern to edit:")
	keys := make([]string, 0, len(patterns))
	i := 1
	for key := range patterns {
		fmt.Printf("%d. %s\n", i, key)
		keys = append(keys, key)
		i++
	}
	var choice int
	fmt.Scanln(&choice)
	oldPattern := keys[choice-1]
	options := []string{"Pattern", "Age Threshold", "Destination", "Delete Flag"}
	println("Choose an option to edit:")
	for i := 0; i < len(options); i++ {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}

	choice = input("Enter your choice: ", func(input string) (int, error) {
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 3 {
			return 0, errors.New("invalid choice")
		}
		return choice, nil
	})

	var ageThreshold int
	var destination string
	var deleteFlag bool
	var pattern string

	switch choice {
	case 1:
		newPattern := input("Enter the new pattern (regex or simple string): ", func(input string) (string, error) {
			return input, nil // No conversion needed for string
		})
		pattern = newPattern
	case 2:
		newAgeThreshold := input("Enter the new age threshold (in days): ", strconv.Atoi)
		ageThreshold = newAgeThreshold
	case 3:
		newDestination := input("Enter the new destination folder: ", func(input string) (string, error) {
			return input, nil // No conversion needed for string
		})
		destination = newDestination
	case 4:
		newDeleteFlag := input("Delete the file? (true/false): ", strconv.ParseBool)
		deleteFlag = newDeleteFlag
	default:
		println("Invalid choice. Exiting...")
	}
	delete(patterns, oldPattern)
	patterns[pattern] = regexInfo{AgeThreshold: ageThreshold, Destination: destination, DeleteFlag: deleteFlag}
	writePatternsToFile(patterns)
}

func deleteFileType() {
	patterns := getSettings(patternsPath)
	println("Choose a Pattern to delete:")
	keys := make([]string, 0, len(patterns))
	i := 1
	for key := range patterns {
		fmt.Printf("%d. %s\n", i, key)
		keys = append(keys, key)
		i++
	}

	choice := input("Enter your choice: ", func(input string) (int, error) {
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(keys) {
			return 0, errors.New("invalid choice")
		}
		return choice, nil
	})

	delete(patterns, keys[choice-1])
	writePatternsToFile(patterns)
}
