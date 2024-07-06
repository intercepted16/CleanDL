package main

import (
	"errors"
	"fmt"
	"strconv"
)

var DefaultOptionsMessage = "Enter your choice: "

func choice(message string, options []string) int {
	for i := 0; i < len(options); i++ {
		fmt.Printf("%d. %s\n", i+1, options[i])
	}
	choice := input(message, func(input string) (int, error) {
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			return 0, errors.New("invalid choice")
		}
		return choice, nil
	})
	return choice
}
