package main

import (
	"fmt"
	"time"

	"github.com/demaggus83/go-builder/pkg/builder"
)

var (
	BuilderBuilt   string
	BuilderGo      string
	BuilderGit     string
	BuilderUPX     string
	BuilderVersion string
)

func main() {
	err := builder.Start(time.Now(), &builder.VersionInfo{
		Built:     BuilderBuilt,
		Go:        BuilderGo,
		GitCommit: BuilderGit,
		UPX:       BuilderUPX,
		Version:   BuilderVersion,
	})
	if err != nil {
		fmt.Printf("Builder returned an error '%s'", err)
	}
}
