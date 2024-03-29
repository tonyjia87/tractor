package utils

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

const netrcFile = `
machine %s
login %s
password %s
`

// helper function to write a netrc file.
func WriteNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}
	out := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	var home string
	u, err := user.Current()
	if err != nil {
		return err
	}
	home = u.HomeDir

	path := filepath.Join(home, ".netrc")
	return ioutil.WriteFile(path, []byte(out), 0600)
}
