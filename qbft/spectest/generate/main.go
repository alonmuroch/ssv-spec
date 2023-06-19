package main

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/bloxapp/ssv-spec/qbft/spectest"
)

//go:generate go run main.go

func main() {
	all := map[string]tests.SpecTest{}
	for _, testF := range spectest.AllTests {
		test := testF()
		post, err := test.GetPostState()
		if err != nil {
			panic(err.Error())
		}
		writeJsonStateComparison(test.TestName(), reflect.TypeOf(test).String(), post)

		n := reflect.TypeOf(test).String() + "_" + test.TestName()
		if all[n] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", n))
		}
		all[n] = test
	}

	byts, err := json.Marshal(all)
	if err != nil {
		panic(err.Error())
	}

	if len(all) != len(spectest.AllTests) {
		panic("did not generate all tests\n")
	}

	fmt.Printf("found %d tests\n", len(all))
	writeJson(byts)
}

func writeJsonStateComparison(name, testType string, post types.Encoder) {
	if post == nil { // If nil, test not supporting post state comparison yet
		fmt.Printf("skipping state comparison json, not supported: %s\n", name)
		return
	}
	fmt.Printf("writing state comparison json: %s\n", name)

	byts, err := json.MarshalIndent(post, "", "		")
	if err != nil {
		panic(err.Error())
	}

	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = filepath.Join(strings.TrimSuffix(basedir, "main.go"), "state_comparison", testType)

	// try to create directory if it doesn't exist
	_ = os.Mkdir(basedir, 0700)
	file := filepath.Join(basedir, fmt.Sprintf("%s.json", name))

	if err := os.WriteFile(file, byts, 0644); err != nil {
		panic(err.Error())
	}
}

func writeJson(data []byte) {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")

	// try to create directory if it doesn't exist
	_ = os.Mkdir(basedir, os.ModeDir)

	file := filepath.Join(basedir, "tests.json")
	fmt.Printf("writing spec tests json to: %s\n", file)
	if err := os.WriteFile(file, data, 0644); err != nil {
		panic(err.Error())
	}
}
