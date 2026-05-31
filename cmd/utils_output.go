package cmd

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/goccy/go-yaml"
	"golang.org/x/term"
)

func output(data any) error {
	osFile, isFile := writer.(*os.File)
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

	if formatFlag == "yaml" {
		return outputYaml(data, isTerminal)
	}
	if formatFlag == "auto" {
		if isTerminal {
			return outputYaml(data, isTerminal)
		}
	}
	return outputJson(data, isTerminal)
}

func outputJson(data any, isTerminal bool) error {
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
	return _output(jsonData, isTerminal)
}

func outputYaml(data any, isTerminal bool) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return _output(yamlData, isTerminal)
}

func _output(data []byte, isTerminal bool) error {
	if isTerminal {
		quick.Highlight(writer, string(data), "yaml", "terminal256", styleFlag)
		return nil
	}
	_, err := writer.Write(data)
	return err
}
