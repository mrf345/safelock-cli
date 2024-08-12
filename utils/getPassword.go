package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	myErrs "github.com/mrf345/safelock-cli/errors"
)

// get the password from pipe or ask the user to enter it
func GetPassword(length int) (password string, err error) {
	pipeInfo, _ := os.Stdin.Stat()

	hasPipe := !strings.HasPrefix(pipeInfo.Mode().String(), "Dcr")

	if !hasPipe {
		fmt.Printf("Enter password (minimum of %d chanters): ", length)
	}

	if password, err = bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
		return
	}

	password = strings.TrimSpace(password)

	if len(password) < length {
		err = &myErrs.ErrInvalidPassword{Len: len(password)}
	}

	return
}
