package utils

import (
	"regexp"
	"runtime"
	"strings"
)

// Thrower returns a possible line, where print/panic was expected.
func Thrower(pkg ...string) (string, int) {
	var goSrcRegexp = regexp.MustCompile(makeSrcRegexp(pkg, `.*\.go`) + `|(libexec/src)`)
	var goTestRegexp = regexp.MustCompile(makeSrcRegexp(pkg, `.*_test\.go`))

	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!goSrcRegexp.MatchString(file) || goTestRegexp.MatchString(file)) {
			return file, line
		}
	}

	return "", 0
}

func makeSrcRegexp(pkg []string, file string) string {
	result := make([]string, 0)
	for _, p := range pkg {
		if p != "" {
			result = append(result, "("+p+"/"+file+")")
		}
	}

	return strings.Join(result, "|")
}
