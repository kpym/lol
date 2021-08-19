// builder package provide an abstraction of latex online compilers.
// Such service shoud provide the Builder interface with single method BuildPDF.
// The Request to send is composed of two parts:
// - Parameters telling how to build the pdf,
// - Files containing the (latex, images...) sources.
package builder

import (
	"fmt"
	"strings"

	"github.com/kpym/lol/log"
)

// Files maps (unix) filename to file data.
// This files are used to build the pdf.
// The files are part of the request.
type Files map[string][]byte

// String provides the Stringer interface for Files.
func (files *Files) String() string {
	w := new(strings.Builder)
	fmt.Fprintln(w, "Files:")
	for fname, fdata := range *files {
		fmt.Fprintf(w, " » %s (%d bytes)\n", fname, len(fdata))
	}

	return w.String()
}

// Parameters contains all context variables.
// In the CLI this variables are set by viper.
type Parameters struct {
	Log       log.Logger
	Service   string
	Compiler  string
	Force     bool
	Biblio    string
	Output    string
	Main      string
	PipedMain bool
	Patterns  []string
}

// String provides the Stringer interface for Parameters.
func (p *Parameters) String() string {
	w := new(strings.Builder)
	fmt.Fprintln(w, "Service:  ", p.Service)
	fmt.Fprintln(w, "Compiler: ", p.Compiler)
	if p.Force {
		fmt.Fprintln(w, "Force:    ", p.Force)
	}
	if p.Biblio != "" {
		fmt.Fprintln(w, "Biblio:   ", p.Biblio)
	}
	if p.Output != "" {
		fmt.Fprintln(w, "Output:   ", p.Output)
	}
	if p.Main != "" {
		fmt.Fprintln(w, "Main:     ", p.Main)
	}
	if p.PipedMain {
		fmt.Fprintln(w, "PipedMain:", p.PipedMain)
	}
	if len(p.Patterns) > 0 {
		fmt.Fprintln(w, "Patterns: ", strings.Join(p.Patterns, ", "))
	}

	return w.String()
}

// Request contains all data necessary to build the pdf.
type Request struct {
	Parameters Parameters
	Files      Files
}

// String provides the Stringer interface for Request.
func (r *Request) String() string {
	return r.Parameters.String() + r.Files.String()
}

// Builder is an interface (service) that can build pdf based on Request.
type Builder interface {
	BuildPDF(Request) ([]byte, error)
}
