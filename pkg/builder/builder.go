package builder

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	ErrNoUPXAvailable = errors.New("upx binary not available")
)

type VersionInfo struct {
	Built     string
	Go        string
	GitCommit string
	UPX       string
	Version   string
}

type Builder struct {
	builds    []*Build
	meta      *Meta
	entryFile string
	distDir   string
	upx       *UPX
}

const ASCII = ` ___      _ _    _
| _ )_  _(_) |__| |___ _ _
| _ \ || | | / _' / -_) '_|
|___/\_,_|_|_\__,_\___|_|
`

func (b *Builder) printBuilderInfo() {
	var output strings.Builder

	output.WriteString(ASCII)
	output.WriteString("===========================\n")

	output.WriteString("Start building '")
	output.WriteString(b.meta.Name)
	output.WriteString("' with:\n")

	output.WriteString("\tBuilt Date:\t")
	output.WriteString(b.meta.BuildTime)
	output.WriteString("\n")

	output.WriteString("\tGo version:\t")
	output.WriteString(runtime.Version())
	output.WriteString("\n")

	if b.meta.GitCommit != "" {
		output.WriteString("\tGit commit:\t")
		output.WriteString(b.meta.GitCommit)
		output.WriteString("\n")
	}

	if b.meta.UserVersion != "" {
		output.WriteString("\tVersion:\t")
		output.WriteString(b.meta.UserVersion)
		output.WriteString("\n")
	}

	if b.upx.Version != "" {
		output.WriteString("\tUPX:\t\t")
		output.WriteString(b.upx.Version)
		output.WriteString("\n")
	}

	fmt.Println(output.String())
}

func (b *Builder) deleteAndCreateDistDir() error {
	err := os.RemoveAll(b.distDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(b.distDir, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) Start() error {
	b.printBuilderInfo()

	fmt.Printf("delete and create dist directory '%s'\n", b.distDir)
	err := b.deleteAndCreateDistDir()
	if err != nil {
		return err
	}

	fmt.Printf("Start the build process...\n\n")

	buildStarted := time.Now()
	for _, build := range b.builds {
		if !b.upx.Available && build.UpxLevel > 0 {
			return ErrNoUPXAvailable
		}

		now := time.Now()
		fmt.Printf("* Building \"%s\"...\n", build.Target.Filename(b.meta.Name))
		distFilePath, err := build.Run(b.meta, b.distDir, b.entryFile)
		if distFilePath == "" && err == nil {
			continue
		}

		dur := time.Since(now).Round(time.Second)
		if err != nil {
			fmt.Printf("Build failed with '%s' after %s\n", err, dur)
			return err
		} else {
			fmt.Printf("\tfinished after %s\n", dur)
		}

		if b.upx.Available && build.UpxLevel > 0 {
			now = time.Now()
			fmt.Printf("* Compressing with UPX...\n")
			result, err := b.upx.Compress(distFilePath, build.UpxLevel, build.Verbose)
			dur := time.Since(now).Round(time.Second)
			if err != nil {
				// UPX is allowed to fail
				fmt.Printf("UPX failed with '%s' after %s\n", err, dur)
				continue
			}

			fmt.Printf("\t%s bytes -> %s bytes (%s)\n",
				result.UncompressedSize, result.CompressedSize, result.Percent)
			fmt.Printf("\tfinished with algorithm '%s' and level '%d' after %s\n",
				result.Algo, build.UpxLevel, dur)
		}

		fmt.Println()
	}

	fmt.Printf("\nAll builds finished after %s\n", time.Since(buildStarted).Round(time.Second))

	return nil
}

func Start(now time.Time, versionInfo *VersionInfo) error {
	builder, cli, err := ParseCLI(os.Args, NewUPX(), now, versionInfo)
	if err != nil {
		cli.ShowUsage()
		return nil
	}

	return builder.Start()
}
