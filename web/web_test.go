package web

import (
	"os"
	"testing"

	"github.com/asmaloney/gactar/util/frameworkutil"
	"github.com/urfave/cli/v2"
)

var webTest *Web = nil

func TestMain(m *testing.M) {
	ctx := cli.NewContext(nil, nil, nil)

	frameworks := frameworkutil.CreateFrameworks(ctx, nil)

	webTest, _ = Initialize(ctx, frameworks, nil)

	exitVal := m.Run()

	os.Exit(exitVal)
}
