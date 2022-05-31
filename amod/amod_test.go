package amod

import (
	"os"
)

func generateToStdout(str string) {
	_, log, _ := GenerateModel(str)
	log.Write(os.Stdout)
}
