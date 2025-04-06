package version

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/pflag"
)

type value int

// Get return the value
func (v *value) Get() interface{} {
	return *v
}

// Set implement pflag.value Set interface
func (v *value) Set(s string) error {
	if s == strRawVersion {
		*v = raw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = boolTrue
	} else {
		*v = boolFalse
	}
	return err
}

// String returns the string representation of the value
func (v *value) String() string {
	if *v == raw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", *v == boolTrue)
}

// Type is the type of the flag as required by the pflag.value interface
func (v *value) Type() string {
	return "version"
}

const (
	boolFalse     value  = 0
	boolTrue      value  = 1
	raw           value  = 2
	flagName      string = "version"
	flagShortHand string = "V"
	strRawVersion string = "raw"
)

var (
	v value = boolFalse
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-main+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit sha1 from git, output of $(git rev-parse HEAD)
	GitCommit = "$Format:%H$"
)

func Show() {
	fmt.Printf(`GitVersion: %s
BuildDate: %s
GitCommit: %s
GoVersion: %s
GOOS/Arch: %s/%s
`,
		GitVersion, BuildDate, GitCommit,
		runtime.Version(), runtime.GOOS, runtime.GOARCH,
	)
}

// AddFlags registers this package's flags on arbitrary FlagSets, such that they
// point to the same value as the global flags.
func AddFlags(flag *pflag.FlagSet) {
	flag.VarP(&v, flagName, flagShortHand, "Print version information and quit.")
	// "--version" will be treated as "--version=true"
	flag.Lookup(flagName).NoOptDefVal = "true"
}

func PrintAndExitIfRequested(appName string) {
	if v == raw {
		Show()
		os.Exit(0)
	} else if v == boolTrue {
		fmt.Printf("%s %s\n", appName, GitVersion)
		os.Exit(0)
	}
}
