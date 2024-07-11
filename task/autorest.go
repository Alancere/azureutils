package task

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alancere/azureutils/common"
	"github.com/Alancere/azureutils/sdk"
)

var autorestLog *slog.Logger

func init() {
	l, err := os.OpenFile("autorest.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	autorestLog = slog.New(slog.NewTextHandler(bufio.NewWriter(l), nil))
}

// test autorest.gotest
// tests: 要生成哪些test类型
func RunAutoRestGoTest(sdkPath string, autorestgo, autorestgotest string, tests ...string) error {
	sdks, err := sdk.GetAllMgmtSDK(sdkPath)
	if err != nil {
		autorestLog.Error("GetAllMgmtSDK", err)
		return err
	}

	for _, s := range sdks {
		autorestLog.Info(fmt.Sprintf("%s/%s", s.Name, s.ArmName))

		outputFolder := s.LocalPath
		autorestMd := filepath.Join(s.LocalPath, "autorest.md")

		// --generate-sdk=false
		if autorestgo == "" {
			autorestgo = "@autorset/go@4.0.0-preview.63"
		}
		if autorestgotest == "" {
			autorestgotest = "@autorest/gotest@4.7.1"
		}
		args := fmt.Sprintf(`--use=%s --use=%s --go --track2 --output-folder=%s --clear-output-folder=false --go.clear-output-folder=false --honor-body-placement=true --remove-unreferenced-types=true `, autorestgo, autorestgotest, outputFolder)
		args += strings.Join(tests, " ")
		args += fmt.Sprintf(" %s", autorestMd)
		output, err := common.AutorestCmd(s.LocalPath, strings.Split(args, " ")...)
		if err != nil {
			autorestLog.Error("AutorestCmd", err, output)
			continue
		}

		// format
		err = common.GoFumpt(s.LocalPath, "-w", ".")
		if err != nil {
			autorestLog.Error("GoFumpt", err)
			continue
		}

		// go mod
		output, err = common.Go(s.LocalPath, "mod", "tidy")
		if err != nil {
			autorestLog.Error("GoMod", err, output)
			continue
		}

		// go vet
		output, err = common.Go(s.LocalPath, "vet", "./...")
		if err != nil {
			autorestLog.Error("GoVet", err, output)
			continue
		}

		// go test
		output, err = common.Go(s.LocalPath, "test", "-v", "./...")
		if err != nil {
			autorestLog.Error("GoTest", err, output)
			continue
		}
	}

	return nil
}
