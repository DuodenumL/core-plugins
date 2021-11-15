package utils

import (
	"strings"

	"github.com/docker/go-units"
	"github.com/urfave/cli/v2"
)

// ParseRAMInHuman returns int value in bytes of a human readable string
// e.g. 100KB -> 102400
func ParseRAMInHuman(ram string) (int64, error) {
	if ram == "" {
		return 0, nil
	}
	flag := int64(1)
	if strings.HasPrefix(ram, "-") {
		flag = int64(-1)
		ram = strings.TrimLeft(ram, "-")
	}
	ramInBytes, err := units.RAMInBytes(ram)
	if err != nil {
		return 0, err
	}
	return ramInBytes * flag, nil
}

// SplitEquality transfers a list of
// aaa=bbb, xxx=yyy into
// {aaa:bbb, xxx:yyy}
func SplitEquality(elements []string) map[string]string {
	r := map[string]string{}
	for _, e := range elements {
		p := strings.SplitN(e, "=", 2)
		if len(p) != 2 {
			continue
		}
		r[p[0]] = p[1]
	}
	return r
}

// ExitCoder wraps a cli Action function into
// a function with ExitCoder interface
func ExitCoder(f func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := f(c); err != nil {
			return cli.Exit(err, -1)
		}
		return nil
	}
}
