// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xdebug

import (
	"easyshop/common/util/xcolor"
	"easyshop/common/util/xstring"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tidwall/pretty"
)

var (
	isTestingMode     bool
	isDevelopmentMode = os.Getenv("JUPITER_MODE") == "dev"
)

// IsTestingMode 判断是否在测试模式下
var onceTest = sync.Once{}

// IsTestingMode ...
func IsTestingMode() bool {
	onceTest.Do(func() {
		isTestingMode = flag.Lookup("test.v") != nil
	})

	return isTestingMode
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return isDevelopmentMode || isTestingMode
}

// IfPanic ...
func IfPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyJsonPrint ...
func PrettyJsonPrint(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		xcolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(xstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}

// PrettyJsonByte ...
func PrettyJsonByte(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(xstring.Json(obj))), pretty.TerminalStyle))
}

// PrettyKV ...
func PrettyKV(key string, val string) {
	fmt.Printf("%-50s => %s\n", xcolor.Red(key), xcolor.Green(val))
}

// PrettyKV ...
func PrettyKVWithPrefix(prefix string, key string, val string) {
	fmt.Printf(prefix+" %-30s => %s\n", xcolor.Red(key), xcolor.Blue(val))
}

// PrettyMap ...
func PrettyMap(data map[string]interface{}) {
	for key, val := range data {
		fmt.Printf("%-20s : %s\n", xcolor.Red(key), fmt.Sprintf("%+v", val))
	}
}

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // 返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1) // 将\替换成/
}