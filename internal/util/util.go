package util

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	"os/exec"
)

// ReadConfigFromFile reads config from file and returns a pointer of config struct.
// If the config file is not found or config value not exist in the file,
// it returns a struct that fields has default value(fields using tag `default` indicates its default value).
func ReadConfigFromFile[T any](path string) (*T, error) {
	c := new(T)
	if err := configor.Load(c, path); err != nil {
		return nil, errors.Wrap(err, "read config from file failed")
	}
	return c, nil
}

func Exec(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "execute command failed")
	}
	return stdout, nil
}
