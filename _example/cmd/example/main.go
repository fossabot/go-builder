package main

import (
	"fmt"
)

var (
	BuilderBuilt   string
	BuilderGo      string
	BuilderGit     string
	BuilderUPX     string
	BuilderVersion string
)

func main() {
	fmt.Println(BuilderBuilt)
	fmt.Println(BuilderGo)
	fmt.Println(BuilderUPX)
	fmt.Println(BuilderGit)
	fmt.Println(BuilderVersion)
}
