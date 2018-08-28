package uuid

import (
	"os/exec"
	"strings"
)

func New() (string, error) {
	if o, e := exec.Command("uuidgen").Output(); e != nil {
		return "", nil
	} else {
		return strings.Trim(string(o), "\n"), nil
	}
}
