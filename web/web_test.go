package web

import (
	"fmt"
	"os"
	"testing"

	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/framework/ccm_pyactr"
	"github.com/asmaloney/gactar/framework/pyactr"
	"github.com/asmaloney/gactar/framework/vanilla_actr"
)

var webTest *Web = nil

func TestMain(m *testing.M) {
	list := framework.ValidFrameworks[1:]

	frameworks := make(framework.List, len(list))

	for _, f := range list {
		var createErr error
		switch f {
		case "ccm":
			frameworks["ccm"], createErr = ccm_pyactr.New(nil)
		case "pyactr":
			frameworks["pyactr"], createErr = pyactr.New(nil)
		case "vanilla":
			frameworks["vanilla"], createErr = vanilla_actr.New(nil)
		}

		if createErr != nil {
			fmt.Println(createErr.Error())
		}
	}

	webTest, _ = Initialize(nil, frameworks, nil)

	exitVal := m.Run()

	os.Exit(exitVal)
}
