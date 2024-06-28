package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func getUserFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}

// Define the source directory (Downloads folder)
func getDownloadsFolder() string {
	return filepath.Join(getUserFolder(), "Downloads")
}

// Define the file types/regex patterns, their corresponding age thresholds (in days), destination folders, and deletion flag
var defaultPatterns = regexPatternsJSON{
	Patterns: regexPatterns{
		".pdf": {
			AgeThreshold: 14,
			Destination:  filepath.Join(os.Getenv("USERPROFILE"), "OneDrive/Documents"),
			DeleteFlag:   false,
		},
		".reg": {
			AgeThreshold: 0,
			Destination:  `C:\bin\reg`,
			DeleteFlag:   false,
		},
		".msi": {
			AgeThreshold: 0,
			Destination:  `C:\bin\msi`,
			DeleteFlag:   false,
		},
		`.*(Installer|Setup)\.exe$`: {
			AgeThreshold: 14,
			Destination:  "",
			DeleteFlag:   true,
		},
		`.*Tool\.exe$`: {
			AgeThreshold: 0,
			Destination:  `C:\bin\exe`,
			DeleteFlag:   false,
		},
		// Add more patterns, thresholds, folders, and deletion flags as needed
	},
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

var fileTypesAndInfoPath = filepath.Join(getUserFolder(), "AppData", "Local", "CleanDL", "patterns.json")

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
		jsonData, err := json.Marshal(defaultPatterns)
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
	// if we os.Open returns an error then handle it
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
	// jsonFile's content into 'users' which we defined above
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
			matched := false
			if strings.HasSuffix(pattern, "$") {
				// Regex pattern match
				matched, err = regexp.MatchString(pattern, file.Name())
				if err != nil {
					panic(err)
				}
			} else {
				// Simple string match
				matched = strings.HasSuffix(file.Name(), pattern)
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
		editSettings()
	case 3:
		os.Exit(0)

	default:
		println("Invalid choice. Exiting...")
	}
}

func organizeFolder() {
	patterns := getSettings(patternsPath)
	downloadsFolder := getDownloadsFolder()
	print("Processing files in: ", downloadsFolder, "\n")
	processFiles(patterns, downloadsFolder)
	print("\nDone!", "\n")
}

func editSettings() {
	options := []string{"Add Pattern", "Edit Pattern", "Delete Pattern", "Exit"}
	println("Choose an option:\n")
	for i := 0; i < len(options); i++ {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		addFileType()
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

func addFileType() {
	patterns := getSettings(patternsPath)
	const (
		AgeThreshold = "AgeThreshold"
		Destination  = "Destination"
		DeleteFlag   = "DeleteFlag"
	)

	var messages = map[string]string{
		AgeThreshold: "Enter the age threshold (in days): ",
		Destination:  "Enter the destination folder: ",
		DeleteFlag:   "Delete the file? (true/false): ",
	}
	var pattern string
	var ageThreshold int
	var destination string
	var deleteFlag bool
	for key, message := range messages {
		fmt.Print(message)
		var value string
		fmt.Scanln(&value)
		switch key {
		case AgeThreshold:
			thisAgeThreshold, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			ageThreshold = thisAgeThreshold
		case Destination:
			destination = value
		case DeleteFlag:
			thisDeleteFlag, err := strconv.ParseBool(value)
			if err != nil {
				panic(err)
			}
			deleteFlag = thisDeleteFlag
		}

	}
	patterns[pattern] = regexInfo{AgeThreshold: ageThreshold, Destination: destination, DeleteFlag: deleteFlag}
	writePatternsToFile(patterns)
}

func editFileType() {
	//TODO: Implement this
	println("WIP")
}

func deleteFileType() {
	patterns := getSettings(patternsPath)
	println("Choose a Pattern to delete:")
	i := 1
	for key := range patterns {
		fmt.Printf("%d. %s\n", i, key)
		i++
	}
	var choice int
	fmt.Scanln(&choice)
	i = 1
	for key := range patterns {
		if i == choice {
			delete(patterns, key)
			break
		}
		i++
	}
	writePatternsToFile(patterns)
}
