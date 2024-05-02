package utils

import (
	"fmt"
	"strings"
)

var Unattended bool = false

func YesNoPrompt(prompt string) bool {
	if Unattended {
		return true
	}

	fmt.Printf("%s (y/[N]): ", prompt)

	var response string
	fmt.Scanln(&response)

	return strings.ToLower(response) == "y"
}
