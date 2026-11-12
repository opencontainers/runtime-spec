package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const usage = `Validate is used to check document with specified schema.
You can use validate in following ways:

   1.specify document file as an argument
      validate <schema.json> <document.json>

   2.pass document content through a pipe
      cat <document.json> | validate <schema.json>

   3.input document content manually, ended with ctrl+d(or your self-defined EOF keys)
      validate <schema.json>
      [INPUT DOCUMENT CONTENT HERE]
`

type httpURLLoader struct{}

func (httpURLLoader) Load(url string) (any, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	return jsonschema.UnmarshalJSON(resp.Body)
}

func main() {
	nargs := len(os.Args[1:])
	if nargs == 0 || nargs > 2 {
		fmt.Printf("ERROR: invalid arguments number\n\n%s\n", usage)
		os.Exit(1)
	}

	if os.Args[1] == "help" ||
		os.Args[1] == "--help" ||
		os.Args[1] == "-h" {
		fmt.Printf("%s\n", usage)
		os.Exit(1)
	}

	schemaPath := os.Args[1]
	if !strings.Contains(schemaPath, "://") {
		var err error
		schemaPath, err = formatFileURL(schemaPath)
		if err != nil {
			fmt.Printf("ERROR: invalid schema-file path: %s\n", err)
			os.Exit(1)
		}
	}

	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft4)
	compiler.UseLoader(jsonschema.SchemeURLLoader{
		"file":  jsonschema.FileLoader{},
		"http":  httpURLLoader{},
		"https": httpURLLoader{},
	})

	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var document io.Reader = os.Stdin
	if nargs > 1 {
		documentPath, err := formatFilePath(os.Args[2])
		if err != nil {
			fmt.Printf("ERROR: invalid document-file path: %s\n", err)
			os.Exit(1)
		}
		f, err := os.Open(documentPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()
		document = f
	}

	instance, err := jsonschema.UnmarshalJSON(document)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := schema.Validate(instance); err != nil {
		fmt.Println("The document is not valid. see errors :")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("The document is valid")
}

func formatFilePath(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func formatFileURL(path string) (string, error) {
	fp, err := formatFilePath(path)
	if err != nil {
		return "", err
	}

	path = filepath.ToSlash(fp)
	if runtime.GOOS == "windows" {
		// FileLoader removes one leading slash on Windows.
		// Preserve both leading slashes for UNC paths.
		path = "/" + path
	}
	return (&url.URL{Scheme: "file", Path: path}).String(), nil
}
