package pkg

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword(msg string) string {
	fmt.Print(msg)
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return ""
	}
	fmt.Println()

	return string(bytePassword)
}
