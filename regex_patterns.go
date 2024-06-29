package main

import (
	"errors"
	"fmt"
	"strconv"
)

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

func addPattern(flags flagPointers) {
	var pattern string
	var ageThreshold int
	var destination string
	var deleteFlag bool
	patterns := getSettings(patternsPath)

	if flags.Pattern == nil {
		pattern = input("Enter the pattern (regex): ", func(input string) (string, error) {
			return input, nil // No conversion needed for string
		})
	} else {
		pattern = *flags.Pattern
	}

	if flags.AgeThreshold == nil {
		ageThreshold = input("Enter the age threshold (in days): ", strconv.Atoi)
	} else {
		ageThreshold = *flags.AgeThreshold
	}

	if flags.Destination == nil {
		destination = input("Enter the destination folder: ", func(input string) (string, error) {
			return input, nil // No conversion needed for string
		})
	} else {
		destination = *flags.Destination
	}

	if flags.DeleteFlag == nil {
		deleteFlag = input("Delete the file? (true/false): ", strconv.ParseBool)
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

func deletePattern() {
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
