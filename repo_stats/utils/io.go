package utils

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
	"time"
)

type Color struct {
	code string
}

func (c Color) String() string {
	return c.code
}

// Predefined color constants
var (
	Err         = Color{"\033[38;2;205;41;73m"}
	Success     = Color{"\033[38;2;127;255;212m"}
	Subtle      = Color{"\033[38;2;105;105;105m"}
	Title       = Color{"\u001B[1m\033[38;2;0;255;255m"}
	TitleNoBold = Color{"\033[38;2;0;255;255m"}
	Highlight   = Color{"\033[38;2;153;102;204m"}
	End         = Color{"\033[0m"}
	None        = Color{""}
)

// GetInput
// get input from the user
//
// Parameters:
//
//   - prompt: The prompt for the user, displayed before input, nothing printed if empty string
//   - promptColor: The Color of the prompt
//
// Returns the user input as a string
func GetInput(prompt string, promptColor Color) string {
	return GetInputAndRespond(prompt, promptColor, "", None)
}

// GetInputAndRespond
// get input from the user and respond with message
//
// Parameters:
//   - prompt: The prompt for the user, displayed before input, nothing printed if empty string
//   - promptColor: The Color of the prompt
//   - response: The response to respond after user input, nothing printed if empty string
//   - responseType: The Color of the response
//
// Returns the user input as a string
func GetInputAndRespond(prompt string, promptColor Color,
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
func Output(message string, messageColor Color) {
	OutputWithTitle("", None, message, messageColor)
}

// OutputWithTitle
// output a message with a title
//
// Parameters:
//   - title: The title of the message, displayed above in bold, empty string for no title
//   - titleColor: The Color of the title
//   - message: The message to be displayed, displayed below title
//   - messageColor: The Color of the message
func OutputWithTitle(title string, titleColor Color, message string, messageColor Color) {
	if title != "" {
		fmt.Printf("\033[1m%s%s%s\n", titleColor, title, End)
	}

	fmt.Printf("%s%s%s\n", messageColor, message, End)
}

// OutputFrom
// Print message items in order with corresponding colors
//
// Parameters:
//   - messageItems: List of strings to print on single line of output
//   - messageColors: List of corresponding colors for each message
//
// Returns error if len(messageItems) != len(messageColors)
func OutputFrom(messageItems []string, messageColors []Color) error {
	if len(messageItems) > len(messageColors) {
		return errors.New("each messageItems string must have a Color")
	}
	if len(messageItems) < len(messageColors) {
		return errors.New("each messageColors color must correspond to message")
	}

	for index := range messageItems {
		message := messageItems[index]
		messageColor := messageColors[index]
		fmt.Printf("%s%s%s", messageColor, message, End)
		if index != len(messageItems)-1 {
			fmt.Print(" ")
		}
	}

	// Print newline at the end of message
	fmt.Println()

	return nil
}

// AnimatedLoader
// Create an animated loading status
//
// Parameters:
//   - text: Text to output prior to ellipses
//   - textColor: Color of text
//   - stopCh: Flag to stop the routine, must be false ro start, should be set to true to stop
//
// Should be run as goroutine, stopped by setting stopCh <- true
func AnimatedLoader(text string, textColor Color, stopCh chan bool) {
	count := 0
	ticker := time.NewTicker(250 * time.Millisecond)

	for {
		select {
		case <-stopCh:
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Printf("\r%s%s%s%s%s\r", textColor, text,
				strings.Repeat(".", count), strings.Repeat(" ", 3-count), End)
			count++
			if count > 3 {
				count = 0
			}
		}
	}
}

// ProgressLoader
// A progress loader which fills the term width
//
// Parameters:
//   - text: String to print before the loader
//   - textColor: The color of the prior text
//   - total: The max value of the loader
//
// Returns a loader function which, when called, progresses and prints the loader
func ProgressLoader(text string, textColor Color, total int) func(string) {
	count := 0
	_text := text
	loader := func(updatedText string) {
		if updatedText != "" {
			_text = updatedText
		}

		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width = total
		}

		// Calculate available width for progress bar
		prefixLen := len(_text) + 3 // _text + " [" + "]"
		barWidth := width - prefixLen
		if barWidth < 10 {
			barWidth = 10
		}

		progressPercent := float32(count) / float32(total)
		progressCount := int(float32(barWidth) * progressPercent)
		remainingCount := barWidth - progressCount

		fmt.Printf("\r%s%s%s %s[%s%s]%s",
			textColor, _text, End,
			Subtle,
			strings.Repeat("#", progressCount),
			strings.Repeat("-", remainingCount),
			End)

		count++
		if count > total {
			count = total
		}
	}
	loader("")
	return loader
}

// UpdatableOutputter
// Creates updatable line of text
//
// Returns function for resetting the output, function for updating the output
func UpdatableOutputter() (func(), func(string, Color)) {
	maxOutput := 0
	return func() {
			fmt.Print("\r" + strings.Repeat(" ", maxOutput))
		},
		func(text string, textColor Color) {
			// Clear old text first, then print new text
			fmt.Printf("\r%s\r%s%s%s", strings.Repeat(" ", maxOutput), textColor, text, End)
			maxOutput = max(maxOutput, len(text))
		}
}

// ReadEnv
// Reads in env vars from a file, places them into a map
//
// Parameters:
//   - file: Reader for environment file
//
// Returns map which maps env var to assignment
func ReadEnv(file io.Reader) (map[string]string, error) {
	envs := make(map[string]string)
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.SplitN(line, "=", 2)
		if len(items) != 2 {
			continue
		}
		// Trim spaces and assign a variable to its item
		envs[strings.Trim(items[0], " ")] = strings.Trim(items[1], " ")
	}
	return envs, nil
}
