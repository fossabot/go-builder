// +build !go1.13

package builder

import "strings"

func (b *Build) getSpecialArgs(meta *Meta, distFilePath, entryFile string) []string {
	return []string{}
}

func (b *Build) getTags() []string {
	return []string{"-tags", strings.Join(tags, " ")}
}
