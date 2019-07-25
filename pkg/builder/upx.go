package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUPXResultWrongFormat = errors.New("UPX output is in wrong format")

	upxResultRegexp = regexp.MustCompile(`.*?(\d+).*?(\d+).+?(\d+\.\d+%).*?(\w.*?)\s`)
)

func detectUpxBinary() (bool, string) {
	var outBuf bytes.Buffer

	cmd := exec.Command("upx", []string{
		"--no-color",
		"--no-progress",
		"--version",
	}...)
	cmd.Stdout = &outBuf

	err := cmd.Run()
	if err != nil {
		return false, ""
	}

	reg := regexp.MustCompile(".*upx\\s(.*?)\\s.*")
	res := reg.FindStringSubmatch(outBuf.String())

	upxVersion, _ := strconv.ParseFloat("3.95", 64)
	if upxVersion < 3.95 {
		fmt.Printf("UPX version '%s' is not supported (need 3.95+)", res[1])
		return false, ""
	}

	return true, res[1]
}

type UPXResult struct {
	Algo             string
	CompressedSize   string
	Percent          string
	UncompressedSize string
}

type UPX struct {
	Available bool
	Version   string
}

func NewUPX() *UPX {
	upx := &UPX{}
	upx.Available, upx.Version = detectUpxBinary()
	return upx
}

func (u *UPX) parseUPXResult(output string) (*UPXResult, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 6 {
		return nil, ErrUPXResultWrongFormat
	}
	res := upxResultRegexp.FindAllStringSubmatch(lines[6], -1)
	if len(res) <= 0 && len(res[0]) < 5 {
		return nil, ErrUPXResultWrongFormat
	}

	return &UPXResult{
		Algo:             res[0][4],
		CompressedSize:   res[0][2],
		Percent:          res[0][3],
		UncompressedSize: res[0][1],
	}, nil
}

func (u *UPX) Compress(path string, level int, verbose bool) (*UPXResult, error) {
	basename := filepath.Base(path)
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext) + "-upx" + ext
	target := filepath.Clean(filepath.Dir(path) + "/" + name)

	cmd := exec.Command("upx", []string{
		"--no-color",
		"--no-progress",
		fmt.Sprintf("-%d", level),
		"-o",
		target,
		path,
	}...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if verbose {
		fmt.Println(buf)
	}

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return u.parseUPXResult(buf.String())
}
