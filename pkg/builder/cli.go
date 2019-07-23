package builder

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const usage = `Usage:
%s [path/to/main.go]

Flags:
`

var (
	ErrNoEntryFile          = errors.New("no entry file given")
	ErrNoValidEntryFilePath = errors.New("no valid path for entry file given")
	ErrNoValidDistPath      = errors.New("no valid path for dist dir given")
	ErrFlagUpxInvalid       = errors.New("invalid upx value given (>=0 && <10)")
)

var allBuildTargets = Triplets{
	{"windows", "amd64", ""},
	{"linux", "amd64", ""},
	{"linux", "amd64", "musl"},
	{"linux", "arm", "6"},  // (e.g. Pi A, A+, B, B+, Zero)
	{"linux", "arm", "7"},  // (e.g. Pi 2, 3) (32bit)
	{"linux", "arm64", ""}, // (e.g. Pi 3, 4) - GOARM is not available!
}

var nameRegExp = regexp.MustCompile(`\/cmd\/(.*?)\/main.go$`)

type Cli struct {
	flagSet *flag.FlagSet

	gitCommit   string
	userVersion string

	buildTargets Triplets
	cgo          bool
	distDir      string
	entryFile    string
	name         string
	upxLevel     int
	verbose      bool
}

func (cli *Cli) clean(path string) (string, error) {
	tmp, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(tmp), nil
}

func (cli *Cli) entryFilePath(args []string) error {
	if len(args) <= 1 {
		return ErrNoEntryFile
	}

	path, err := cli.clean(args[1])
	if err != nil {
		return ErrNoValidEntryFilePath
	}

	err = cli.exists(path)
	if err != nil {
		return err
	}

	cli.entryFile = path

	return nil
}

func (cli *Cli) exists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoValidEntryFilePath
		}
		return err
	}
	return nil
}

func (cli *Cli) createFlagSet(args []string) {
	cli.flagSet = flag.NewFlagSet("menu", flag.ContinueOnError)
	cli.flagSet.Usage = func() {
		fmt.Printf(usage, args[0])
		cli.flagSet.PrintDefaults()
	}

	cli.flagSet.BoolVar(&cli.cgo, "cgo", false, "enable cgo (e.g. for github.com/mattn/sqlite3)")
	cli.flagSet.StringVar(&cli.distDir, "dist", "./dist", "the output directory for the binaries")
	cli.flagSet.StringVar(&cli.name, "name", "", "the name of the binary (defaults to the entry file name)")
	cli.flagSet.Var(&cli.buildTargets, "target", `defaults to: "windows-amd64,linux-amd64,linux-amd64-musl,linux-arm-6,linux-arm-7,linux-arm64"`)
	cli.flagSet.IntVar(&cli.upxLevel, "upx", 0, `enables binary compression per upx (0=disabled / 1=faster / 9=better`)
	cli.flagSet.BoolVar(&cli.verbose, "verbose", false, "show stdout and stderr from build and upx")

	cli.flagSet.StringVar(&cli.gitCommit, "git", "", "the git commit hash which gets injected into dmain.BuilderVarGitCommit")
	cli.flagSet.StringVar(&cli.userVersion, "version", "", "the version which gets injected into main.BuilderVarVersion")
}

func (cli *Cli) ShowUsage() {
	cli.flagSet.Usage()
}

func (cli *Cli) detectName() string {
	// is /some/path/to/cmd/$name/main.go scheme
	cmdName := nameRegExp.FindStringSubmatch(cli.entryFile)
	if len(cmdName) > 1 {
		return cmdName[1]
	}

	// use the filename from entry file /some/path/to/$name.go
	_, f := filepath.Split(cli.entryFile)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

func ParseCLI(args []string, upx *UPX, now time.Time) (*Builder, *Cli, error) {
	cli := &Cli{
		buildTargets: make(Triplets, 0, len(allBuildTargets)),
	}
	cli.createFlagSet(args)

	err := cli.entryFilePath(args)
	if err != nil {
		return nil, cli, err
	}

	err = cli.flagSet.Parse(args[2:])
	if err != nil {
		return nil, cli, err
	}

	if cli.name == "" {
		cli.name = cli.detectName()
	}

	if cli.upxLevel < 0 || cli.upxLevel > 9 {
		return nil, cli, ErrFlagUpxInvalid
	}

	cli.distDir, err = cli.clean(cli.distDir)
	if err != nil {
		return nil, cli, ErrNoValidDistPath
	}

	if len(cli.buildTargets) == 0 {
		cli.buildTargets = allBuildTargets
	}

	meta := &Meta{
		BuildTime:   now.Format(time.RFC1123),
		GitCommit:   cli.gitCommit,
		GoVersion:   runtime.Version(),
		Name:        cli.name,
		UpxVersion:  upx.Version,
		UserVersion: cli.userVersion,
	}

	builds := []*Build{}
	for _, t := range cli.buildTargets {
		builds = append(builds, &Build{
			Target:   t,
			Cgo:      cli.cgo,
			UpxLevel: cli.upxLevel,
			Verbose:  cli.verbose,
		})
	}

	builder := &Builder{
		builds: builds,
		meta:   meta,
		upx:    upx,

		entryFile: cli.entryFile,
		distDir:   cli.distDir,
	}

	return builder, cli, nil
}
