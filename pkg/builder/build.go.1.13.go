// +build go1.13

package builder

import "strings"

func (b *Build) getSpecialArgs(meta *Meta, distFilePath, entryFile string) []string {
	return []string{"-trimpath"}
}

// -tags tag,list : a comma-separated list of build tags to consider satisfied during the build...
// todo: may need to modify when go 1.14 is released https://github.com/golang/go/issues/26492
func (b *Build) getTags() []string {
	return []string{"-tags", strings.Join(tags, ",")} // todo: go 1.13 prefers ,
}
