package demo_test

import (
	"fmt"
	"testing"
)

// init的作用，一个文件中可以有多个init吗
func init() {
	fmt.Println("init1")
}

func init() {
	fmt.Println("init2")
}

func TestInit(t *testing.T) {}
