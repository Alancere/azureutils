package task_test

import (
	"testing"

	"github.com/Alancere/azureutils/autorest"
	"github.com/Alancere/azureutils/task"
)

func TestRunAutoRestGoTest(t *testing.T) {
	sdkPath := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	autorest_go := "D:/Go/src/github.com/Azure/autorest.go/packages/autorest.go"
	autorest_gotest := "D:/Go/src/github.com/Azure/autorest.go/packages/autorest.gotest"
	tests := []string{
		"--generate-sdk=true",
		// autorest.GOTestOption.Example,
		autorest.GOTestOption.FakeTest,
		// autorest.GOTestOption.MockTest,
		// autorest.GOTestOption.Sample,
		// autorest.GOTestOption.ScenarioTest,
	}
	// 运行 fake_test.go
	err := task.RunAutoRestGoTest(sdkPath, autorest_go, autorest_gotest, tests...)
	if err != nil {
		t.Fatal(err)
	}
}