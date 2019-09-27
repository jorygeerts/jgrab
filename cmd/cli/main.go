package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/andrew-d/go-termutil"
	"io"
	"log"
	"os"
	"strings"
)

type nodesFlag []string

func (f *nodesFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *nodesFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var nf nodesFlag
	flag.Var(&nf, "node", "--node NAME_OF_NODE --node NAME_OF_OTHER_NODE")
	flag.Parse()

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		_, _ = fmt.Printf(" some-program | %s [flags]\n", os.Args[0])
		_, _ = fmt.Printf(" %s [flags] some-json\n", os.Args[0])
		flag.PrintDefaults()
	}

	canRun := false
	if !termutil.Isatty(os.Stdin.Fd()) {
		canRun = true
		processStream(os.Stdin, makeNodeFilter(nf))
	}

	if flag.Arg(0) != "" {
		canRun = true
		processString(flag.Arg(0), makeNodeFilter(nf))
	}

	if !canRun {
		flag.Usage()
		os.Exit(1)
	}

}

type nodeFilter func(map[string]interface{}) map[string]interface{}

func makeNodeFilter(nf nodesFlag) nodeFilter {
	if len(nf) == 0 {
		return func(m map[string]interface{}) map[string]interface{} {
			return m
		}
	}

	nodeList := map[string]bool{}

	for _, k := range nf {
		nodeList[k] = true
	}

	return func(in map[string]interface{}) map[string]interface{} {
		out := make(map[string]interface{})

		for key, _ := range nodeList {
			if value, ok := in[key]; ok {
				out[key] = value
			}
		}

		return out
	}
}

func isPipeMode(fi os.FileInfo) bool {

	// Unix
	if fi.Mode()&os.ModeCharDevice == os.ModeCharDevice {
		return true
	}

	// Windows
	if fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
		return true
	}

	return false
}

func processStream(r io.Reader, nf nodeFilter) {
	dec := json.NewDecoder(r)

	output := log.New(os.Stdout, "", log.LstdFlags)
	for {
		var data map[string]interface{}
		if err := dec.Decode(&data); err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		chuncks := []string{}
		for key, value := range nf(data) {
			chuncks = append(chuncks, key+" -> "+fmt.Sprintf("%s", value))
		}

		output.Println(strings.Join(chuncks, ", "))
	}
}

func processString(s string, nf nodeFilter) {
	output := log.New(os.Stdout, "", 0)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return
	}

	chuncks := []string{}
	for key, value := range nf(data) {
		chuncks = append(chuncks, key+" -> "+fmt.Sprintf("%s", value))
	}

	output.Println(strings.Join(chuncks, ", "))
}
