package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Define the source directory (Downloads folder)
func getDownloadsFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, "Downloads")
}

// Define the file types/regex patterns, their corresponding age thresholds (in days), destination folders, and deletion flag
var fileTypesAndInfo = map[string]struct {
	ageThreshold int
	destination  string
	deleteFlag   bool
}{
	".pdf":                      {14, filepath.Join(os.Getenv("USERPROFILE"), "OneDrive/Documents"), false},
	".reg":                      {0, `C:\bin\reg`, false},
	".msi":                      {0, `C:\bin\msi`, false},
	`.*(Installer|Setup)\.exe$`: {14, "", true},
	`.*Tool\.exe$`:              {0, `C:\bin\exe`, false},
	// Add more patterns, thresholds, folders, and deletion flags as needed
}

// Function to move or delete files based on age, type, and regex or simple string
func processFiles(downloadsFolder string) {
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

		for pattern, info := range fileTypesAndInfo {
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
				if fileAgeDays > info.ageThreshold {
					fmt.Printf("File age: %d", fileAgeDays)
					fmt.Printf("Age threshold: %d", info.ageThreshold)
					if info.deleteFlag {
						os.Remove(filePath) // Delete the file
						fmt.Printf("Deleted: %s\n", filePath)
					} else if info.destination != "" {
						os.Rename(filePath, filepath.Join(info.destination, file.Name())) // Move the file
						fmt.Printf("Moved: %s to %s\n", filePath, info.destination)
					}
					break // Exit the loop after processing
				}
			}
		}
	}
}

func main() {
	downloadsFolder := getDownloadsFolder()
	print("Processing files in: ", downloadsFolder, "\n")
	processFiles(downloadsFolder)
	print("\nDone!", "\n")
}
