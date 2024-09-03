package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	slErrs "github.com/mrf345/safelock-cli/slErrs"
)

// get the password from pipe or ask the user to enter it
func GetPassword(length int) (password string, err error) {
	pipeInfo, _ := os.Stdin.Stat()

	hasPipe := !strings.HasPrefix(pipeInfo.Mode().String(), "Dcr")

	if !hasPipe {
		fmt.Printf("Enter password (minimum of %d chanters): ", length)
	}

	if password, err = bufio.NewReader(os.Stdin).ReadString('\n'); err != nil && err != io.EOF {
		return
	}

	password = strings.TrimSpace(password)

	if len(password) < length {
		err = &slErrs.ErrInvalidPassword{Len: len(password), Need: length}
	}

	return
}
