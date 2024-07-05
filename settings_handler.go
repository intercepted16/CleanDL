package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				fmt.Println("Error closing patterns file")
			}
		}(file)

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
	defer func(settingsFile *os.File) {
		err = settingsFile.Close()
		if err != nil {
			println("Error closing settings file")
		}
	}(settingsFile)
	byteValue, _ := io.ReadAll(settingsFile)
	// we initialize our custom regex array
	var regexPatternsJSON regexPatternsJSON

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'regexPatternsJSON' which we defined above
	err = json.Unmarshal(byteValue, &regexPatternsJSON)
	if err != nil {
		return nil
	}
	var regexPatterns = regexPatternsJSON.Patterns
	return regexPatterns
}
