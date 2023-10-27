package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	tfhoon "github.com/wim-web/tfoon"
)

func TreesPresenter() func(tfhoon.ModuleTreeList) (string, error) {
	return func(trees tfhoon.ModuleTreeList) (string, error) {
		b, err := json.Marshal(trees)
		return string(b), err
	}
}

func M2ePresenter() func(tfhoon.Module2EntryPoint) (string, error) {
	return func(m2e tfhoon.Module2EntryPoint) (string, error) {
		b, err := json.Marshal(m2e)
		return string(b), err
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	// --separator -s flag
	separator := flag.String("separate", ",", "default is comma")
	m2e := flag.Bool("m2e", false, "module to entrypoint mode")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		logger.Error("e.g.: tfoon path1,path2")
		os.Exit(1)
	}

	paths := strings.Split(args[0], *separator)
	output, err := getOutput(paths, *m2e, logger)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(output)
}

func getOutput(paths []string, m2eMode bool, logger *slog.Logger) (string, error) {
	trees, err := tfhoon.FromPaths(paths)
	if err != nil {
		return "", err
	}

	var output string

	if m2eMode {
		m2e := trees.ToModule2EntryPoint()
		presenter := M2ePresenter()
		output, err = presenter(m2e)
	} else {
		presenter := TreesPresenter()
		output, err = presenter(trees)
	}

	if err != nil {
		return "", err
	}

	return output, nil
}
