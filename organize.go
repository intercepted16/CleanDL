package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

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

// Function to organize the downloads folder
func organizeFolder() {
	patterns := getSettings(patternsPath)
	downloadsFolder := getDownloadsFolder()
	print("Processing files in: ", downloadsFolder, "\n")
	processFiles(patterns, downloadsFolder)
	print("\nDone!", "\n")
}
