package main

import (
	"fmt"
	"time"

	"github.com/demaggus83/go-builder/pkg/builder"
)

func main() {
	err := builder.Start(time.Now())
	if err != nil {
		fmt.Printf("Builder returned an error '%s'", err)
	}
}
