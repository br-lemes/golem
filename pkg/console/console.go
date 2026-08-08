package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/goccy/go-yaml"
	"golang.org/x/term"
)

var (
	Debug  bool      = false
	Format string    = "auto"
	Stderr io.Writer = os.Stderr
	Stdin  io.Reader = os.Stdin
	Stdout io.Writer = os.Stdout
	Style  string    = "monokai"
)

func Auto(data any) error {
	osFile, isFile := Stdout.(*os.File)
	isTerminal := isFile && term.IsTerminal(int(osFile.Fd()))

	bytes, ok := data.([]byte)
	if ok {
		var v any
		err := json.Unmarshal(bytes, &v)
		if err != nil {
			return err
		}
		data = v
	}

	if Format == "yaml" {
		return Yaml(data, isTerminal)
	}
	if Format == "auto" {
		if isTerminal {
			return Yaml(data, isTerminal)
		}
	}
	return Json(data, isTerminal)
}

func Json(data any, isTerminal bool) error {
	var jsonData []byte
	var err error
	if isTerminal {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	return output(jsonData, isTerminal)
}

func Yaml(data any, isTerminal bool) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return output(yamlData, isTerminal)
}

func output(data []byte, isTerminal bool) error {
	if isTerminal {
		_ = quick.Highlight(Stdout, string(data), "yaml", "terminal256", Style)
		return nil
	}
	_, err := Stdout.Write(data)
	return err
}

func Confirm(message string) bool {
	Printf("%s (y/N): ", message)
	reader := bufio.NewReader(Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func Input(message string) string {
	var input string
	Printf("%s: ", message)
	_, err := fmt.Scan(&input)
	if err != nil {
		return ""
	}
	return input
}

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(Stdout, format, a...)
}

func Debugf(format string, a ...any) {
	if Debug {
		_, _ = fmt.Fprintf(Stderr, "DEBUG: "+format, a...)
	}
}

func Errorf(format string, a ...any) {
	_, _ = fmt.Fprintf(Stderr, format, a...)
}
