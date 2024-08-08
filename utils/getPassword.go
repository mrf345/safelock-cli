package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	myErrs "github.com/mrf345/safelock-cli/errors"
)

func GetPassword() (password string, err error) {
	pipeInfo, _ := os.Stdin.Stat()

	hasPipe := !strings.HasPrefix(pipeInfo.Mode().String(), "Dcr")

	if !hasPipe {
		fmt.Print("Enter password (minimum of 8 chanters): ")
	}

	if password, err = bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
		return
	}

	password = strings.TrimSpace(password)

	if len(password) < 8 {
		err = &myErrs.ErrInvalidPassword{Len: len(password)}
	}

	return
}
