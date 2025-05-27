package main

import (
	"fmt"
)

type Color struct {
	code string
}

func (c Color) String() string {
	return c.code
}

// Predefined color constants
var (
	Err       = Color{"\033[38;2;205;41;73m"}
	Success   = Color{"\033[38;2;127;255;212m"}
	Subtle    = Color{"\033[38;2;105;105;105m"}
	Title     = Color{"\033[38;2;0;255;255m"}
	Highlight = Color{"\033[38;2;153;102;204m"}
	End       = Color{"\033[0m"}
	None      = Color{""}
)

func getInput(prompt string, promptType Color) string {
	return getInputAndRespond(prompt, promptType, "", None)
}

func getInputAndRespond(prompt string, promptType Color,
	response string, responseType Color) string {
	if prompt != "" {
		fmt.Printf("%s%s: %s", promptType, prompt, End)
	}
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return ""
	}
	if response != "" {
		fmt.Printf("%s%s%s\n", responseType, response, End)
	}
	return input
}

func outputWithTitle(title string, titleColor Color, message string, messageColor Color) {
	if title != "" {
		fmt.Printf("\033[1m%s%s%s\n", titleColor, title, End)
	}

	fmt.Printf("%s%s%s\n", messageColor, message, End)
}

func ouput(message string, messageColor Color) {
	outputWithTitle("", None, message, messageColor)
}
