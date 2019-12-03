package main

import (
	"context"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type verbosity struct {
	silent   bool
	standard bool
	verbose  bool
	debug    bool
	mode     string
}

var (
	_verbosity       verbosity
	action           string
	commandStructure = make(map[string]int)
	flags            = make(map[string]int)
	errLogger        *log.Logger
	outLogger        *log.Logger
)

func init() {
	errLogger = log.New(os.Stderr, "", 0)
	outLogger = log.New(os.Stdout, "", 0)

	setVerbosity()
	for i, c := range os.Args[1:] {
		commandStructure[c] = i
	}
	for key, index := range commandStructure {
		if regexp.MustCompile(`^-.*`).MatchString(key) {
			delete(commandStructure, key)
			flags[key] = index
		} else if regexp.MustCompile(`^--.*`).MatchString(key) {
			delete(commandStructure, key)
			flags[key] = index
		} else {
			commandStructure[key] = index
		}
	}
	for action := range commandStructure {
		switch action {
		case "build":
			if err := build(); err != nil {
				errLogger.Fatal(err)
			}
		case "test":
			if err := test(); err != nil {
				errLogger.Fatal(err)
			}
		case "bench":
			if err := bench(); err != nil {
				errLogger.Fatal(err)
			}
		case "lint":
			if err := lint(); err != nil {
				errLogger.Fatal(err)
			}
		case "fmt":
			if err := fmt(); err != nil {
				errLogger.Fatal(err)
			}
		case "get":
			if err := get(); err != nil {
				errLogger.Fatal(err)
			}
		}
	}
}

func setVerbosity() {
	for _, arg := range os.Args {
		if arg == "-s" || arg == "--silent" {
			_verbosity.silent = true
			_verbosity.mode = "silent"
		}
		if arg == "-d" || arg == "--debug" {
			_verbosity.debug = true
			_verbosity.mode = "debug"
		}
		if arg == "-v" || arg == "--verbose" {
			_verbosity.verbose = true
			_verbosity.mode = "verbose"
		}

	}
}

func execute(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	if !_verbosity.silent {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func build() error {
	return execute(
		ctx,
		viper.GetString("cc"),
		"build",
		func() string {
			if _verbosity.debug {
				return "-x"
			}
			return ""
		}(),
		func() string {
			if _verbosity.verbose {
				return ""
			}
			return ""
		}(),
		viper.GetString("sourceFile"),
		func() string {
			if viper.GetString("outputFile") != "" {
				return "-o"
			}
			return ""
		}(),
		func() string {
			if viper.GetString("outputFile") != "" {
				return viper.GetString("outputFile")
			}
			return ""
		}(),
		func() string {
			if viper.GetString("sourceFile") != "" {
				return viper.GetString("sourceFile")
			}
			return ""
		}(),
	)
}

func test() error {
	return execute(ctx, viper.GetString("cc"), "test")
}

func bench() error {
	return execute(ctx, viper.GetString("cc"), "test", "-bench=.")
}

func lint() error {
	return execute(ctx, viper.GetString("cc"), "lint")
}

func fmt() error {
	return execute(ctx, viper.GetString("cc"), "fmt")
}

func get() error {
	return execute(ctx, viper.GetString("cc"), "get")
}

// Configuration defines the configuration structure for the gomake.yaml file
type Configuration struct {
	SourceFiles     []io.Writer `yaml:"sourceFile"`
	OutputFile      []io.Writer `yaml:"outputFile"`
	CleanOnFailure  bool        `yaml:"cleanOnFailure"`
	Verbosity       verbosity   `yaml:"_verbosity"`
	CC              string      `yaml:"CC"`
	OverrideTargets struct {
		Build   string `yaml:"build"`
		Install string `yaml:"install"`
		Lint    string `yaml:"lint"`
		Fmt     string `yaml:"fmt"`
		Test    string `yaml:"test"`
		Bench   string `yaml:"bench"`
		Get     string `yaml:"get"`
	} `yaml:"override"`
}
