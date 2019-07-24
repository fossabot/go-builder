package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Meta struct {
	BuildTime   string
	GitCommit   string
	GoVersion   string
	Name        string
	UpxVersion  string
	UserVersion string
}

type Build struct {
	Cgo      bool
	Target   *Triplet
	UpxLevel int
	Verbose  bool
}

var tags = []string{
	"osusergo", "netgo", "static_build",
}

func (b *Build) getBuilderVars(meta *Meta) string {
	var builderVars strings.Builder
	builderVars.WriteString(fmt.Sprintf(`-X "main.BuilderBuilt=%s" `, meta.BuildTime))
	builderVars.WriteString(fmt.Sprintf(`-X "main.BuilderGo=%s" `, runtime.Version()))
	if meta.UpxVersion != "" {
		builderVars.WriteString(fmt.Sprintf(`-X "main.BuilderUPX=%s" `, meta.UpxVersion))
	}
	if meta.GitCommit != "" {
		builderVars.WriteString(fmt.Sprintf(`-X "main.BuilderGit=%s" `, meta.GitCommit))
	}
	if meta.UserVersion != "" {
		builderVars.WriteString(fmt.Sprintf(`-X "main.BuilderVersion=%s" `, meta.UserVersion))
	}
	return builderVars.String()
}

// todo: may need to modify when go 1.14 is released https://github.com/golang/go/issues/26492
func (b *Build) getCgoFlags() string {
	cgoFlags := ""
	if b.Cgo {
		cgoFlags = "-extldflags '-fno-PIC -static' -linkmode=external"
	}
	return cgoFlags
}

func (b *Build) getArgs(meta *Meta, distFilePath, entryFile string) []string {
	args := []string{
		"build",
		"-a",
	}
	args = append(args, b.getSpecialArgs(meta, distFilePath, entryFile)...)
	args = append(args, b.getTags()...)
	// ldflags
	args = append(args, "-ldflags", fmt.Sprintf(`-s -w %s%s`, b.getBuilderVars(meta), b.getCgoFlags()))

	// output / input
	return append(args, "-o", distFilePath, entryFile)
}

// Returns the distFilePath
func (b *Build) Run(meta *Meta, distPath, entryFile string) (string, error) {
	if !b.Cgo && b.Target.Special == "musl" {
		fmt.Println("skipping build...")
		return "", nil
	}

	distFilePath := filepath.Join(distPath, b.Target.Filename(meta.Name))

	cmd := exec.Command("go", b.getArgs(meta, distFilePath, entryFile)...)

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOOS=%s", b.Target.GOOS),
		fmt.Sprintf("GOARCH=%s", b.Target.GOARCH),
	)

	if b.Target.GOARCH == "arm" && b.Target.Special != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARM=%s", b.Target.Special))
	}
	
	cgo := 0
	if b.Cgo {
		cgo = 1
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("CGO_ENABLED=%d", cgo),
			fmt.Sprintf("CC=%s", b.Target.Compiler()),
		)
	}

	cmd.Stderr = os.Stderr

	if b.Verbose {
		cmd.Stdout = os.Stdout
	}

	return distFilePath, cmd.Run()
}
