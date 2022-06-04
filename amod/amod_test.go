package amod

import (
	"fmt"
	"os"
)

func generateToStdout(str string) {
	_, log, _ := GenerateModel(str)
	err := log.Write(os.Stdout)
	if err != nil {
		fmt.Print(err.Error())
	}
}
