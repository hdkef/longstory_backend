package helper

import "os"

func GetAbsPath() (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}
