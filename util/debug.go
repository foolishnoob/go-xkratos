package util

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

var dumpLock sync.Mutex

func Dump(a ...interface{}) {
	defer dumpLock.Unlock()
	stacks := strings.Split(string(debug.Stack()), "\n\t")
	coroutineRegexp, _ := regexp.Compile("(goroutine \\d+ )")
	matchesCoroutines := coroutineRegexp.FindStringSubmatch(stacks[0])

	pathRegexp, _ := regexp.Compile("(.*\\.go:\\d+ )")
	matchesPaths := pathRegexp.FindStringSubmatch(stacks[3])

	dumpLock.Lock()
	fmt.Printf("%s : %s\n", matchesCoroutines[0], matchesPaths[0])
	spew.Dump(a...)
}
