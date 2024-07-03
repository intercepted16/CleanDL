package main

import (
	"errors"
	"os"
	"regexp"
	"strconv"
)

func getPattern() string {
	// get pattern input with error handling
	pattern := input("Enter the pattern (regex): ", func(input string) (string, error) {
		var error error
		if input != "" {
			if !isValidRegex(input) {
				error = errors.New("invalid regex pattern")
			}
		} else {
			error = errors.New("pattern cannot be empty")
		}
		if error != nil {
			return input, error
		}
		return input, nil
	})

	return pattern
}

func getAgeThreshold() int {
	// get age threshold input with error handling
	ageThreshold := input("Enter the age threshold (in days): ", strconv.Atoi)
	return ageThreshold
}

func getDestination() string {
	// get destination input with error handling
	destination := input("Enter the destination folder: ", func(input string) (string, error) {
		var error error
		if input != "" {
			// verify that the directory exists
			if _, err := os.Stat(input); os.IsNotExist(err) {
				error = errors.New("directory does not exist")
			}
		} else {
			error = errors.New("destination cannot be empty")
		}
		if error != nil {
			return input, error
		}
		return input, nil
	})
	return destination
}

func getDeleteFlag() bool {
	// get delete flag input with error handling
	deleteFlag := input("Delete the file? (true/false): ", strconv.ParseBool)
	return deleteFlag
}

func crudPatterns(flags flagPointers) {
	options := []string{"Add Pattern", "Edit Pattern", "Delete Pattern", "Exit"}
	choice := choice("Enter an option to edit: ", options)
	switch choice {
	case 1:
		addPattern(flags)
	case 2:
		editPattern()
	case 3:
		deletePattern()
	case 4:
		clearScreen()
		main()
	default:
		println("Invalid choice. Exiting...")
	}
}

func isValidRegex(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

func addPattern(flags flagPointers) {
	var pattern string
	var ageThreshold int
	var destination string
	var deleteFlag bool
	patterns := getSettings(patternsPath)

	if flags.Pattern == nil {
		pattern = getPattern()
	} else {
		pattern = *flags.Pattern
	}

	if flags.AgeThreshold == nil {
		ageThreshold = getAgeThreshold()
	} else {
		ageThreshold = *flags.AgeThreshold
	}

	if flags.Destination == nil {
		destination = getDestination()
	} else {
		destination = *flags.Destination
	}

	if flags.DeleteFlag == nil {
		deleteFlag = getDeleteFlag()
	} else {
		deleteFlag = *flags.DeleteFlag
	}
	delete(patterns, pattern)
	patterns[pattern] = regexInfo{AgeThreshold: ageThreshold, Destination: destination, DeleteFlag: deleteFlag}
	writePatternsToFile(patterns)
}

func editPattern() {
	patterns := getSettings(patternsPath)
	keys := make([]string, 0, len(patterns))
	for key := range patterns {
		keys = append(keys, key)
	}
	patternToEdit := choice("Choose a pattern to edit: ", keys)
	oldPattern := keys[patternToEdit-1]
	options := []string{"Pattern", "Age Threshold", "Destination", "Delete Flag"}
	optionToEdit := choice("Choose an option to edit: ", options)

	var ageThreshold int
	var destination string
	var deleteFlag bool
	var pattern string

	switch optionToEdit {
	case 1:
		newPattern := getPattern()
		pattern = newPattern
	case 2:
		newAgeThreshold := getAgeThreshold()
		ageThreshold = newAgeThreshold
	case 3:
		newDestination := getDestination()
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

func deletePattern() {
	patterns := getSettings(patternsPath)
	keys := make([]string, 0, len(patterns))
	for key := range patterns {
		keys = append(keys, key)
	}
	patternToDelete := choice("Choose a pattern to delete: ", keys)

	delete(patterns, keys[patternToDelete-1])
	writePatternsToFile(patterns)
}
