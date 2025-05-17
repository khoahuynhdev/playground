package util

import (
	"os"
)

// FileExists checks if a file or directory exists at the given path then returns boolean and error
func FileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// Schrodinger: file may or may not exist. See err for details.
	// do *NOT* use !os.IsNotExist(err) to test for file existence then
	return false, err
}

func CreateDirectory(path string) error {
	// Check if the directory already exists
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return err
}

func WriteFile(path, content string, mode int, overwrite bool) (bool, error) {
	var fileMode os.FileMode
	if mode == 0 {
		fileMode = os.FileMode(0600)
	} else {
		fileMode = os.FileMode(mode)
	}
	fileExists, err := FileExists(path)
	if err != nil {
		return false, err
	}
	if !fileExists {
		// if not exists, create file with context
		err = os.WriteFile(path, []byte(content), fileMode)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if fileExists && overwrite {
		// if exists and overwrite is true, create file with context
		err = os.WriteFile(path, []byte(content), fileMode)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// file exists and not overwrite with no error
	return false, nil
}
