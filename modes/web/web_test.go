package web

import (
	"os"
	"testing"

	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/frameworkutil"
)

var webTest *Web = nil

func TestMain(m *testing.M) {
	settings := &cli.Settings{}

	frameworks := frameworkutil.CreateFrameworks(settings, nil)

	settings.Frameworks = frameworks

	webTest, _ = Initialize(settings, 8181, nil)

	exitVal := m.Run()

	os.Exit(exitVal)
}
