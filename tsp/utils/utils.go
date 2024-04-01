package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

func MergeJson(dir string, dst string) error {
	var err error
	dirStat, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !dirStat.IsDir() {
		return fmt.Errorf("%s is not directory", dir)
	}

	filePaths := []string{}
	jsonParsed := gabs.New()

	fmt.Println("merge paths:")
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || strings.Contains(path, "examples") || d.Name() == "openapi.json" {
			return nil
		}

		fmt.Println("\t", path)

		if filepath.Ext(path) == ".json" {
			filePaths = append(filePaths, path)
			//
			newPaseed, err := gabs.ParseJSONFile(path)
			if err != nil {
				return fmt.Errorf("gabs parseJson error: %v", err)
			}
			// if err = jsonParsed.Merge(newPaseed); err != nil {
			// 	return fmt.Errorf("gabs merge json error: %v", err)
			// }

			if err = jsonParsed.MergeFn(newPaseed, func(destination, source interface{}) interface{} {
				return destination
			}); err != nil {
				return fmt.Errorf("gabs merge json error: %v", err)
			}

		}

		return nil
	})
	if err != nil {
		return err
	}

	// fmt.Println("============gabs============")
	// // 过滤重复的数据
	// if err = os.WriteFile("gabs.json", jsonParsed.BytesIndent("", "  "), 0666); err != nil {
	// 	log.Fatal(err)
	// }

	// Specs format
	// spec := Spec{}
	spec := spec.Swagger{}
	myjson := jsoniter.ConfigCompatibleWithStandardLibrary
	err = myjson.Unmarshal(jsonParsed.Bytes(), &spec)
	if err != nil {
		// log.Fatal("spec type unmarshal err:", err)
		return err
	}
	// marshalSpec, err := json.Marshal(&spec)
	marshalSpec, err := json.MarshalIndent(&spec, "", "  ")
	if err != nil {
		// log.Fatal("spec type marshal err:", err)
		return err
	}
	// fmt.Println(string(marshalSpec))
	if err = os.WriteFile(dst, marshalSpec, 0o666); err != nil {
		// log.Fatal(err)
		return err
	}

	return nil
}

// format json to format.json file
func FormatJson(specPath string, dst string) error {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return err
	}

	spec := spec.Swagger{}
	myjson := jsoniter.ConfigCompatibleWithStandardLibrary
	err = myjson.Unmarshal(data, &spec)
	if err != nil {
		return err
	}

	marshalSpec, err := json.MarshalIndent(&spec, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(dst, marshalSpec, 0o666); err != nil {
		return err
	}
	return nil
}

func ComparePath(first, second string, dst string) error {
	spec1, err := loads.JSONSpec(first)
	if err != nil {
		return err
	}

	spec2, err := loads.JSONSpec(second)
	if err != nil {
		return err
	}

	fmt.Println("compare paths:")
	// diff paths
	path1 := make(map[string]bool, len(spec1.Analyzer.AllPaths()))
	for k := range spec1.Spec().Paths.Paths {
		path1[k] = false
	}
	path2 := make(map[string]bool, len(spec2.Analyzer.AllPaths()))
	for k := range spec2.Spec().Paths.Paths {
		path2[k] = false
	}

	for k := range path1 {
		if _, ok := path2[k]; ok {
			path1[k] = true
			path2[k] = true
		}
	}

	firstPaths := make([]string, 0, len(path1))
	for k, v := range path1 {
		if !v {
			firstPaths = append(firstPaths, k)
		}
	}

	secondPaths := make([]string, 0, len(path2))
	for k, v := range path2 {
		if !v {
			secondPaths = append(secondPaths, k)
		}
	}

	// write in paths.json file
	pathFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	pathFile.WriteString("## Compare Path\n\n")

	sort.Strings(firstPaths)
	fmt.Println("\tfirst file paths not container in second:", len(firstPaths))
	pathFile.WriteString(fmt.Sprintf("\n### %s not included in the other:\n\n", first))
	for _, v := range firstPaths {
		fmt.Println("\t", v, lookupOperationId(spec1.Spec().Paths.Paths[v]))
		// fmt.Println("\t", lookupOperationId(spec1.Spec().Paths.Paths[v]), v)
		pathFile.WriteString(fmt.Sprintf("%s %s\n", v, lookupOperationId(spec1.Spec().Paths.Paths[v])))
	}

	sort.Strings(secondPaths)
	fmt.Println("\tsecond file paths not container in first:", len(secondPaths))
	pathFile.WriteString(fmt.Sprintf("\n### %s not included in the other:\n\n", second))
	for _, v := range secondPaths {
		fmt.Println("\t", v, lookupOperationId(spec2.Spec().Paths.Paths[v]))
		// fmt.Println("\t", lookupOperationId(spec2.Spec().Paths.Paths[v]), v)
		pathFile.WriteString(fmt.Sprintf("%s %s\n", v, lookupOperationId(spec2.Spec().Paths.Paths[v])))
	}

	pathFile.Sync()

	return nil
}

func lookupOperationId(pathItem spec.PathItem) string {
	// if pathItem == nil {
	// 	return ""
	// }

	if pathItem.Post != nil {
		return pathItem.Post.ID
	}

	if pathItem.Put != nil {
		return pathItem.Put.ID
	}

	if pathItem.Patch != nil {
		return pathItem.Patch.ID
	}

	if pathItem.Get != nil {
		return pathItem.Get.ID
	}

	if pathItem.Delete != nil {
		return pathItem.Delete.ID
	}

	return ""
}

func CompareFile(c1, c2 string) error {
	f1, err := os.ReadFile(c1)
	if err != nil {
		return err
	}

	f2, err := os.ReadFile(c2)
	if err != nil {
		return err
	}

	fmt.Println(cmp.Diff(f1, f2, nil)) // cmpopts.SortMaps(f1))

	return nil
}
