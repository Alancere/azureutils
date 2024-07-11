package template_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Alancere/azureutils/template"

	stdTemplate "text/template"
)

func TestParseTmplate(t *testing.T) {
	sdkInfo := map[string]string{
		"rpName":      "myrp",
		"packageName": "arm-myrp",

		"packageTitle": "Azure Compute",

		"packageVersion": "1.0.0",
		"releaseDate":    "2020-01-01",

		"NewClientName": "", // value 参考 ReplaceReadmeNewClientName

		"alancere": "alancere",
	}

	dir := "D:\\Projects\\azureutils\\template\\templates"

	SayHello := func(name string) string {
		return fmt.Sprintf("Hi %s Gopher!!!", name)
	}
	funcMap := stdTemplate.FuncMap{
		"SayHello": SayHello,
	}
	_ = funcMap
	fmt.Println(SayHello("Alancere"))

	err := template.ParseTemplates(dir, "testdata", sdkInfo, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFunc(t *testing.T) {
	// SayHello := func (char string) (string, error) {
	// 	return "Hi" + char, nil
	// }

	tpl, err := stdTemplate.New("hello.tpl").ParseFiles("./templates/hello.tpl") // .Funcs(stdTemplate.FuncMap{"SayHello": SayHello})
	if err != nil {
		t.Fatal(err)
	}

	userName := "xxx"
	// 5. 渲染
	err = tpl.Execute(os.Stdout, userName)
	if err != nil {
		t.Fatal(err)
	}
}
