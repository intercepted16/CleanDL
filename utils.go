package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/fatih/color"

	"github.com/urfave/cli/v2"
)

func getDownloadsFolder() string {
	return filepath.Join(getUserFolder(), "Downloads")
}

func getUserFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}

func getFlag[T any](cCtx *cli.Context, flagName string) *T {
	if !cCtx.IsSet(flagName) {
		return nil
	}

	var result T
	resultType := reflect.TypeOf(result)
	var value interface{}

	switch resultType.Kind() {
	case reflect.String:
		value = cCtx.String(flagName)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = cCtx.Int(flagName)
	case reflect.Bool:
		value = cCtx.Bool(flagName)
	default:
		// Handle unsupported types
		return nil
	}

	// This is safe because we're checking the type above

	// Convert the value to type T and return a pointer to it
	result = value.(T)
	return &result
}

func input[T any](prompt string, parseFunc func(string) (T, error)) T {
	var result T
	for {
		print(prompt)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			if err == io.EOF {
				// Exit on EOF as it is likely the user is trying to quit
				color.Red("Exiting...")
				os.Exit(0)
			}
			// Non critical error, just print it and continue
			println("Error reading input, please try again.")
			continue
		}
		value, err := parseFunc(input)
		if err == nil {
			result = value
			break
		}
		color.Red("Invalid input: %s", err.Error())
	}
	return result
}

func clearScreen() {
	print("\033[H\033[2J")
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	if _, err := file.Write(jsonData); err != nil {
		panic(err)
	}
}
