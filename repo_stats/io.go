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

// Get the under input

// Parameters:
//
//   - prompt: The prompt for the user, displayed before input, nothing printed if empty string
//   - promptColor: The Color of the prompt
//
// Returns the user input as a string
func getInput(prompt string, promptColor Color) string {
	return getInputAndRespond(prompt, promptColor, "", None)
}

// Get user input and provide a response
//
// Parameters:
//   - prompt: The prompt for the user, displayed before input, nothing printed if empty string
//   - promptColor: The Color of the prompt
//   - response: The response to respond after user input, nothing printed if empty string
//   - responseType: The Color of the response
//
// Returns the user input as a string
func getInputAndRespond(prompt string, promptColor Color,
	response string, responseType Color) string {
	if prompt != "" {
		fmt.Printf("%s%s: %s", promptColor, prompt, End)
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

// Output a message
//
// Parameters:
//   - message: The message to be displayed
//   - messageColor: The Color of the message
func ouput(message string, messageColor Color) {
	outputWithTitle("", None, message, messageColor)
}

// Output a message with a title
//
// Parameters:
//   - title: The title of the message, displayed above in bold, empty string for no title
//   - titleColor: The Color of the title
//   - message: The message to be displayed, displayed below title
//   - messageColor: The Color of the message
func outputWithTitle(title string, titleColor Color, message string, messageColor Color) {

	if title != "" {
		fmt.Printf("\033[1m%s%s%s\n", titleColor, title, End)
	}

	fmt.Printf("%s%s%s\n", messageColor, message, End)
}
