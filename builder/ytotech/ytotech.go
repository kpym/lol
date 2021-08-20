// ytotech package provides a builder.Builder interface to latex.ytotech.com service.
// The request send to latex.ytotech.com is a single json that looks like this
// ```json
// {
//     "compiler": "pdflatex",
//     "resources": [
//         {
//             "main": true,
//             "file": "...base64 encoded file..."
//         },
//         {
//             "path": "logo.png",
//             "file": "...base64 encoded file..."
//         }
//     ]
// }
// ```
package ytotech

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kpym/lol/builder"
)

// *ytotech is a Builder.
type ytotech struct{}

// NewBuilder provides a new Builder interface to latex.ytotech.com.
func NewBuilder() builder.Builder {
	return new(ytotech)
}

// reqToJson encode (part of) the Request as json that is send to latex.ytotech.com.
func reqToJson(req builder.Request) string {
	json := new(strings.Builder)
	json.WriteString(`{`)
	fmt.Fprintf(json, `"compiler":"%s",`, req.Parameters.Compiler)
	if req.Parameters.Biblio != "" {
		fmt.Fprintf(json, `"options":{"bibliography":{"command":"%s"}},`, req.Parameters.Biblio)
	}
	json.WriteString(`"resources": [`)
	addComma := false
	for fname, fdata := range req.Files {
		if addComma {
			json.WriteString(`,`)
		} else {
			addComma = true
		}
		if fname == req.Parameters.Main {
			json.WriteString(`{"main": true,`)
		} else {
			fmt.Fprintf(json, `{"path": "%s",`, fname)
		}
		json.WriteString(`"file": "`)
		json.WriteString(base64.StdEncoding.EncodeToString(fdata))
		json.WriteString(`"}`)
	}
	json.WriteString(`]}`)

	return json.String()
}

// compilationError corresponds to the json returned in case of error.
type compilationError struct {
	Error string
	Logs  string
}

// BuildPDF send the request to latex.ytotech.com and returns the resulting pdf.
func (y *ytotech) BuildPDF(req builder.Request) ([]byte, error) {
	// prepare the json to submit
	body := strings.NewReader(reqToJson(req))
	// send comile request
	resp, err := http.Post("https://latex.ytotech.com/builds/sync", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read pdf or error from response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Problem reading response: %w\n", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// respBody contains a json encoded compilationError
		var comperr compilationError
		err = json.Unmarshal(respBody, &comperr)
		if err != nil {
			return nil, fmt.Errorf("LaTeX-on-HTTP compilation error (status code %d). The answer is not a valid json:\n%s\n", resp.StatusCode, respBody)
		}
		return nil, fmt.Errorf("LaTeX-on-HTTP compilation error (status code %d):\n%s\n", resp.StatusCode, comperr.Logs)
	}

	// respBody contains the resulting pdf
	return respBody, nil
}
