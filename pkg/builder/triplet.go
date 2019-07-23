package builder

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrNoValidTriplet = errors.New("no valid triplet")
)

type Triplets []*Triplet

func (t *Triplets) String() string {
	l := len(*t)

	var b strings.Builder

	const sep = ","
	for i, v := range *t {
		b.WriteString(v.String())
		if i < l-1 {
			b.WriteString(sep)
		}
	}

	return b.String()
}

func (t *Triplets) Set(value string) error {
	values := strings.Split(value, ",")

	for _, v := range values {
		trip, err := NewTriplet(v)
		if err != nil {
			return err
		}

		*t = append(*t, trip)
	}

	return nil
}

// All available GOOS and GOARCH are listed in go src/go/build/syslist.go
type Triplet struct {
	GOOS    string
	GOARCH  string
	Special string // GOARM 6|7 or musl
}

func NewTriplet(triplet string) (*Triplet, error) {
	tmp := strings.Split(triplet, "-")
	if len(tmp) < 3 {
		return nil, ErrNoValidTriplet
	}

	special := ""
	if len(tmp) >= 3 {
		special = tmp[2]
	}

	return &Triplet{
		GOOS:    tmp[0],
		GOARCH:  tmp[1],
		Special: special,
	}, nil
}

func (t *Triplet) Compiler() string {
	if t.GOOS == "windows" {
		return "x86_64-w64-mingw32-gcc"
	}

	if t.GOOS == "linux" {
		if t.Special == "musl" {
			return "musl-gcc"
		}

		if t.GOARCH == "arm" {
			cc := "arm-linux-gnueabihf-gcc-8"
			i, _ := strconv.Atoi(t.Special)
			if i <= 6 {
				cc = "arm-linux-gnueabi-gcc-8"
			}
			return cc
		}

		if t.GOARCH == "arm64" {
			return "aarch64-linux-gnu-gcc-8"
		}

		return "x86_64-linux-gnu-gcc"
	}
	return ""
}

func (t *Triplet) Filename(name string) string {
	var f strings.Builder

	f.WriteString(name)
	f.WriteString("-")

	f.WriteString(t.GOOS)
	f.WriteString("-")

	if t.GOOS == "linux" && t.Special == "musl" {
		f.WriteString("musl")
		f.WriteString("-")
	}

	f.WriteString(t.GOARCH)
	if t.GOARCH == "arm" {
		// arm6, arm7...
		f.WriteString(t.Special)
	}

	if t.GOOS == "windows" {
		f.WriteString(".exe")
	}

	return f.String()
}

func (t *Triplet) String() string {
	var b strings.Builder

	b.WriteString(t.GOOS)
	b.WriteString("-")
	b.WriteString(t.GOARCH)

	if t.Special != "" {
		b.WriteString("-")
		b.WriteString(t.Special)
	}

	return b.String()
}
