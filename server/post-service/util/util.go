package util

import (
	"encoding/base64"
	"os"
)

// basic slice operations

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Search(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}

func SaveFile(fileBase64, uuidString string) error {
	// decode base64 string
	dec, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		return err
	}

	// create new dir
	err = os.MkdirAll("media/"+uuidString, 0755)
	if err != nil {
		return err
	}

	// create file
	f, err := os.Create("media/" + uuidString + "/" + uuidString)
	if err != nil {
		return err
	}
	defer f.Close()

	// write decoded base64 string in file
	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func LoadFile(uuidString string) (string, error) {
	// check if file exists
	if _, err := os.Stat("media/" + uuidString + "/" + uuidString); err != nil {
		return "", err
	}

	// read file and get content
	content, err := os.ReadFile("media/" + uuidString + "/" + uuidString)
	if err != nil {
		return "", err
	}

	// encode file content to base64 string
	return base64.StdEncoding.EncodeToString(content), nil
}
