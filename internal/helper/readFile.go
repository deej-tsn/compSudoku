package helper

import (
	"os"
)

func ReadFileToString(filepath string) string {
	file, err := os.ReadFile(filepath)
	CheckError(err)
	return string(file)
}
