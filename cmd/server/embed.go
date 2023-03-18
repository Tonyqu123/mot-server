package main

import (
	_ "embed"
	"fmt"
)

//go:embed env/local.yaml
var s string

func init() {
	fmt.Printf("file local.yaml: %s\n", s)
}
