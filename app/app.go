package app

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kpym/lol/builder"
	"github.com/kpym/lol/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// parameters constants
const (
	// The name of our config file, without the file extension
	// because viper supports many different config file languages.
	defaultConfigFilename = "lol"

	// The environment variable prefix of all environment variables.
	// For example, --server flag is bound to $LOL_SERVER.
	envPrefix = "LOL"

	// If the main file is piped to stdin, this name is used.
	MainNameIfStdin = "main_from_stdin.tex"
)

// The version that is set by goreleaser
var version = "dev"

// Help displays usage message if -h/--help flag is set or in case of falg error.
func Help() {
	var out = os.Stderr
	fmt.Fprintf(out, "lol (version: %s)\n", version)
	fmt.Fprintln(out, "LaTeX online compiler. More info at www.github.com/kpym/lol.")
	fmt.Fprintln(out, "\nAvailable options:")
	pflag.PrintDefaults()

	fmt.Fprintln(out, "\nExamples:")
	fmt.Fprintln(out, "> lol main.tex")
	fmt.Fprintln(out, "> lol  -s ytotech -c xelatex main.tex")
	fmt.Fprintln(out, "> lol main.tex personal.sty images/img*.pdf")
	fmt.Fprintln(out, "> cat main.tex | lol -c lualatex -o out.pdf")
	fmt.Fprintln(out, "")
}

// InitFlag define the CLI flags.
func InitFlags() {
	pflag.StringP("service", "s", "", "Service can be laton or ytotex.")
	pflag.StringP("compiler", "c", "pdflatex", "One of pdflatex,xelatex or lualatex.\nFor ytotex platex, uplatex and context are also available.\n")
	pflag.BoolP("force", "f", false, "Do not use the laton cache. Force compile. Ignored by ytotech.")
	pflag.StringP("biblio", "b", "", "Can be bibtex or biber for ytotex. Not used by laton.")
	pflag.StringP("output", "o", "", "The name of the pdf file. If empty, same as the main tex file.")
	pflag.StringP("main", "m", "", "The main tex file to compile.")
	pflag.BoolP("quiet", "q", false, "Prevent any output.")
	pflag.BoolP("verbose", "v", false, "Print info and errors. No debug info is printed.")
	pflag.Bool("debug", false, "Print everithing (debug info included).")
	pflag.Parse()
}

// stringsIn checks if the first argument is equal to one of the following parameters.
// Used in GetParameters only.
func stringIn(str string, values ...string) bool {
	for _, v := range values {
		if str == v {
			return true
		}
	}
	return false
}

// GetParameters use pflag and viper to set the parameters.
func GetParameters(params *builder.Parameters) error {
	v := viper.New()

	// Bind the current command's flags to viper
	v.BindPFlags(pflag.CommandLine)

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Set environment variables prefix.
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables.
	// I we have --compsed-flags we can use :
	// v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Transfer the parameters values to params struct.
	err := v.Unmarshal(params)
	if err != nil {
		return err
	}

	// set log level
	level := log.ErrorLevel
	if v.GetBool("quiet") {
		level = log.Quiet
	}
	if v.GetBool("verbose") {
		level = log.InfoLevel
	}
	if v.GetBool("debug") {
		level = log.DebugLevel
	}
	// the default writer is os.Stdout (color.Output)
	params.Log = log.New(log.WithLevel(level), log.WithColor())

	// normalise the service name
	params.Service = strings.ToLower(params.Service)
	// chack if the service support the requested options
	if !stringIn(params.Service, "laton", "ytotech", "") {
		return fmt.Errorf("Unknown %s service.", params.Service)
	}
	if stringIn(params.Compiler, "platex", "uplatex", "context") {
		if params.Service == "laton" {
			return fmt.Errorf("Laton do not support %s compiler.", params.Compiler)
		}
		if params.Service == "" {
			params.Service = "ytotech"
		}
	} else if !stringIn(params.Compiler, "pdflatex", "xelatex", "lualatex") {
		return fmt.Errorf("Non supported %s compiler.", params.Compiler)
	}
	if params.Biblio != "" {
		if params.Service == "laton" {
			return fmt.Errorf("Laton do not support %s bibliography.", params.Biblio)
		}
		if params.Service == "" {
			params.Service = "ytotech"
		}
	}
	if params.Service == "" {
		// TODO : choose the fastest ?
		params.Service = "laton"
	}
	// check if the input is piped
	fi, err := os.Stdin.Stat()
	if err == nil {
		params.PipedMain = ((fi.Mode() & os.ModeCharDevice) == 0) && (fi.Mode()&os.ModeNamedPipe != 0)
		params.Log.Debug("Piped input: %v, Stdin mode: %v.", params.PipedMain, fi.Mode())
	}
	// get the patterns
	params.Patterns = append(pflag.Args(), params.Patterns...)
	if len(params.Patterns) == 0 && params.Main == "" && !params.PipedMain {
		return fmt.Errorf("Missing file to compile.")
	}
	if params.Main != "" && params.PipedMain {
		return fmt.Errorf("Main file can't be set when there is piped input.")
	}
	// set the main file (if needed)
	if params.Main == "" {
		if !params.PipedMain {
			params.Main = params.Patterns[0]
		}
	} else {
		params.Patterns = append([]string{params.Main}, params.Patterns...)
	}

	// set the output (if not piped input)
	if params.Output == "" && params.Main != "" {
		params.Output = strings.TrimSuffix(params.Main, ".tex") + ".pdf"
	}

	// set Main if piped input
	if params.PipedMain {
		params.Main = MainNameIfStdin
	}

	return nil
}

// GetFiles read all files based on params.Patterns.
func GetFiles(params builder.Parameters) (builder.Files, error) {
	// temporary variables
	var (
		err      error
		filedata []byte
	)
	// files to be read
	files := make(builder.Files)
	// get the main file
	if params.PipedMain {
		params.Log.Debug("Read the main file from stdin.")
		filedata, err = io.ReadAll(os.Stdin)
	} else {
		params.Log.Debug("Read the main file from %s.", params.Main)
		filedata, err = os.ReadFile(params.Main)
	}
	files[params.Main] = filedata
	if err != nil {
		return nil, fmt.Errorf("Error while reading the main file: %w", err)
	}
	// get all other files (if any) that are readable
	for _, pat := range params.Patterns {
		// check if is folder or pattern
		patInfo, err := os.Stat(pat)
		if err == nil {
			if patInfo.IsDir() {
				pat = path.Join(pat, "*")
			}
		}
		names, _ := filepath.Glob(pat)
		for _, fname := range names {
			// if on Windows, transform to unix name
			uname := filepath.ToSlash(fname)
			// if this file is already present
			if _, ok := files[uname]; ok {
				continue
			}
			// read the file, or skipt it if not readable
			filedata, err = os.ReadFile(fname)
			if err == nil {
				files[uname] = filedata
				params.Log.Debug("File %s (%d bytes) added to the list.", uname, len(filedata))
			} else {
				params.Log.Debug("Probleam reading support file (we skip it): %s.", fname)
			}
		}
	}

	return files, nil
}
