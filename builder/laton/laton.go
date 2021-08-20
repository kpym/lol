// laton package provides a builder.Builder interface to latexonline.cc service.
// The request send to latexonline.cc is composed by
// a single tar.gz file that contains all sources
// and by url parameters:
// - target : containing the name of the main file
// - compiler : containing the compiler (pdflatex|xelatex|lualatex)
// - force : skip the cached version if present
// In case of success the response body contains the pdf.
// In case of error the response body contains the log.
package laton

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/kpym/lol/builder"
)

// *laton is a Builder.
type laton struct{}

// NewBuilder provides a new Builder interface to latexonline.cc.
func NewBuilder() builder.Builder {
	return new(laton)
}

// filesToTar stock all source files in a single .tar.gz.
func filesToTar(files builder.Files) ([]byte, error) {
	var err error
	// .tar.gz buffer
	var tarbuf bytes.Buffer
	gzw := gzip.NewWriter(&tarbuf)
	tw := tar.NewWriter(gzw)

	for name, data := range files {
		// prepare header
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(data)),
		}
		// write header
		err = tw.WriteHeader(hdr)
		if err != nil {
			return nil, err
		}
		// write file
		_, err = tw.Write(data)
		if err != nil {
			return nil, err
		}
	}
	// close tar
	err = tw.Close()
	if err != nil {
		return nil, err
	}
	// close gzip
	err = gzw.Close()
	if err != nil {
		return nil, err
	}

	return tarbuf.Bytes(), nil
}

// newTarRequest prepare the http.Request to be send to latexonline.cc.
// The values from params are encoded as url values and the tardata is send as request body.
func newTarRequest(params builder.Parameters, tardata []byte) (*http.Request, error) {
	// write the tar.gz as request body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "laton.tar.gz")
	if err != nil {
		return nil, err
	}
	part.Write(tardata)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// encode the params values in the url
	urlParams := url.Values{}
	urlParams.Add("target", params.Main)
	if params.Force {
		urlParams.Add("force", "true")
	}
	urlParams.Add("command", params.Compiler)

	// return the request
	httpReq, err := http.NewRequest("POST", "https://texlive2020.latexonline.cc/data?"+urlParams.Encode(), body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", writer.FormDataContentType())
	return httpReq, nil
}

// BuildPDF send the request to latexonline.cc and returns the resulting pdf.
func (y *laton) BuildPDF(req builder.Request) ([]byte, error) {
	var err error
	// prepare the tar file to submit
	tardata, err := filesToTar(req.Files)
	if err != nil {
		return nil, err
	}
	// create a request
	httpReq, err := newTarRequest(req.Parameters, tardata)
	if err != nil {
		return nil, err
	}
	// send compile request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
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
		return nil, fmt.Errorf("Laton compilation error (status code %d):\n%s\n", resp.StatusCode, respBody)
	}

	// respBody contains the resulting pdf
	return respBody, nil
}
