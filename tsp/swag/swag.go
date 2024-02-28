package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Alancere/azureutils/tsp/utils"
)

// go build -o swag.exe .
func main() {
	if len(os.Args) < 3 {
		log.Fatal("please input success args: arg[1]=originDir, arg[2]=compileFile")
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	subDir := workDir
	if len(os.Args) == 4 {
		subDir = filepath.Join(subDir, os.Args[3])
		_, err = os.Stat(subDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if err = os.Mkdir(subDir, 0666); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	fmt.Printf("Args:\n\t%s\n\t%s\n", os.Args[1], os.Args[2])

	fmt.Println("merge path files to merge.json...")
	mergePath := filepath.Join(subDir, "merge.json")
	formatPath := filepath.Join(subDir, "format.json")
	if err := utils.MergeJson(os.Args[1], mergePath); err != nil {
		log.Fatal(err)
	}

	fmt.Println("format compile file to format.json...")
	if err := utils.FormatJson(os.Args[2], formatPath); err != nil {
		log.Fatal(err)
	}
	// if err := swag.CompareFile("merge.json", "format.json"); err != nil {
	// 	log.Fatal(err)
	// }
	if err := utils.ComparePath(mergePath, formatPath, filepath.Join(subDir, "paths.md")); err != nil {
		log.Fatal(err)
	}

	fmt.Println("merge.json path:", mergePath)
	fmt.Println("format.json path:", formatPath)
}