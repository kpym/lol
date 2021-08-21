package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kpym/lol/app"
	"github.com/kpym/lol/builder"
	"github.com/kpym/lol/builder/laton"
	"github.com/kpym/lol/builder/ytotech"
	"github.com/kpym/lol/log"
	"github.com/spf13/pflag"
)

// Error checking
func check(logger log.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	var params builder.Parameters

	// init the flags
	pflag.Usage = app.Help
	pflag.CommandLine.SortFlags = false
	app.InitFlags()

	// get parameters from flags, envs and config file
	err = app.GetParameters(&params)
	check(params.Log, err)

	// get the files content based on params.Patterns
	files, err := app.GetFiles(params)
	check(params.Log, err)

	// build the pdf
	var compiler builder.Builder
	switch params.Service {
	case "ytotech":
		compiler = ytotech.NewBuilder()
	case "laton":
		compiler = laton.NewBuilder()
	default:
		check(params.Log, fmt.Errorf("Unknown service %s", params.Service))
	}
	req := builder.Request{Parameters: params, Files: files}
	params.Log.Info("Send request with the following parameters:\n%s", req.String())
	sendtime := time.Now()
	pdf, err := compiler.BuildPDF(req)
	params.Log.Info("Answer received in %1.1f seconds.", time.Since(sendtime).Seconds())
	check(params.Log, err)

	// write the pdf
	if params.Output != "" {
		params.Log.Info("Write %s.", params.Output)
		err = os.WriteFile(params.Output, pdf, 0644)
		check(params.Log, err)
	} else {
		params.Log.Info("Write to stdout.")
		_, err = os.Stdout.Write(pdf)
		check(params.Log, err)
	}
}
