package demo_test

import (
	"fmt"
	"testing"

	"github.com/iancoleman/strcase"
)

func TestStrCase_ToSnake(t *testing.T) {
	fmt.Println(strcase.ToSnake("StrCase"))

	fmt.Println(strcase.ToSnake("resources"))

	fmt.Println(strcase.ToSnake("managed-identity"))

	fmt.Println(strcase.ToSnake("hello world"))
	fmt.Println(strcase.ToSnake("hello      world"))
}
