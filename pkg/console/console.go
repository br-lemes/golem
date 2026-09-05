package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/goccy/go-yaml"
	"golang.org/x/term"
)

var (
	Debug     bool      = false
	Color     bool      = false
	Format    string    = "auto"
	Stderr    io.Writer = os.Stderr
	Stdin     io.Reader = os.Stdin
	Stdout    io.Writer = os.Stdout
	Style     string    = "monokai"
	Exclude   []string
	Only      []string
	ExcludeIf []string
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
	if len(Exclude) > 0 {
		var err error
		data, err = excludePaths(data, Exclude)
		if err != nil {
			return err
		}
	}
	if len(Only) > 0 {
		var err error
		data, err = onlyPaths(data, Only)
		if err != nil {
			return err
		}
	}
	if len(ExcludeIf) > 0 {
		var err error
		data, err = excludeIfPaths(data, ExcludeIf)
		if err != nil {
			return err
		}
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

func excludeIfPaths(data any, expressions []string) (any, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var value any
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	for _, expression := range expressions {
		fields := strings.Fields(expression)
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid exclude-if expression: %s", expression)
		}
		parts := strings.Split(fields[0], ".")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid exclude-if path: %s", fields[0])
		}
		value, err = excludeIfPath(value, parts, fields[1], fields[2])
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func excludeIfPath(value any, parts []string, operator, expected string) (any, error) {
	current, ok := value.(map[string]any)
	if !ok {
		return value, nil
	}
	for key, child := range current {
		matched, err := path.Match(parts[0], key)
		if err != nil {
			return nil, err
		}
		if !matched {
			continue
		}
		matches, err := conditionMatches(child, parts[1:], operator, expected)
		if err != nil {
			return nil, err
		}
		if matches {
			delete(current, key)
		}
	}
	return current, nil
}

func conditionMatches(value any, parts []string, operator, expected string) (bool, error) {
	if len(parts) == 0 {
		return compareValue(value, operator, expected)
	}
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			matched, err := path.Match(parts[0], key)
			if err != nil {
				return false, err
			}
			if matched {
				result, err := conditionMatches(child, parts[1:], operator, expected)
				if result || err != nil {
					return result, err
				}
			}
		}
	case []any:
		if parts[0] == "*" {
			for _, child := range current {
				result, err := conditionMatches(child, parts[1:], operator, expected)
				if result || err != nil {
					return result, err
				}
			}
		}
	}
	return false, nil
}

func compareValue(value any, operator, expected string) (bool, error) {
	actual, ok := value.(float64)
	want, numberErr := strconv.ParseFloat(expected, 64)
	if ok && numberErr == nil {
		switch operator {
		case ">":
			return actual > want, nil
		case ">=":
			return actual >= want, nil
		case "<":
			return actual < want, nil
		case "<=":
			return actual <= want, nil
		case "==":
			return actual == want, nil
		case "!=":
			return actual != want, nil
		}
	}
	if operator == "==" || operator == "!=" {
		matches := fmt.Sprint(value) == expected
		if operator == "!=" {
			return !matches, nil
		}
		return matches, nil
	}
	return false, fmt.Errorf("invalid exclude-if comparison: %s %s", operator, expected)
}

func onlyPaths(data any, patterns []string) (any, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var value any
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	paths := make([][]string, 0, len(patterns))
	for _, pattern := range patterns {
		paths = append(paths, strings.Split(pattern, "."))
	}
	value, _ = onlyPath(value, paths)
	return value, nil
}

func onlyPath(value any, paths [][]string) (any, bool) {
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			keep := false
			var childPaths [][]string
			for _, pathParts := range paths {
				matched, err := path.Match(pathParts[0], key)
				if err != nil || !matched {
					continue
				}
				if len(pathParts) == 1 {
					keep = true
					break
				}
				childPaths = append(childPaths, pathParts[1:])
			}
			if !keep && len(childPaths) > 0 {
				updated, childKeep := onlyPath(child, childPaths)
				if childKeep {
					current[key] = updated
					keep = true
				}
			}
			if !keep {
				delete(current, key)
			}
		}
		return current, true
	case []any:
		keepAny := false
		for index := range current {
			indexString := strconv.Itoa(index)
			keep := false
			var childPaths [][]string
			for _, pathParts := range paths {
				if pathParts[0] != "*" && pathParts[0] != indexString {
					continue
				}
				if len(pathParts) == 1 {
					keep = true
					break
				}
				childPaths = append(childPaths, pathParts[1:])
			}
			if !keep && len(childPaths) > 0 {
				updated, childKeep := onlyPath(current[index], childPaths)
				if childKeep {
					current[index] = updated
					keep = true
				}
			}
			if keep {
				keepAny = true
			} else {
				current[index] = nil
			}
		}
		return current, keepAny
	default:
		return value, false
	}
}

func excludePaths(data any, patterns []string) (any, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var value any
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	for _, pattern := range patterns {
		parts := strings.Split(pattern, ".")
		value, _ = excludePath(value, parts)
	}
	return value, nil
}

func excludePath(value any, parts []string) (any, bool) {
	if len(parts) == 0 {
		return nil, true
	}
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			matched, err := path.Match(parts[0], key)
			if err != nil || !matched {
				continue
			}
			if len(parts) == 1 {
				delete(current, key)
				continue
			}
			updated, remove := excludePath(child, parts[1:])
			if remove {
				delete(current, key)
			} else {
				current[key] = updated
			}
		}
	case []any:
		if parts[0] == "*" {
			if len(parts) == 1 {
				return nil, true
			}
			for index, child := range current {
				updated, _ := excludePath(child, parts[1:])
				current[index] = updated
			}
		} else {
			index, err := strconv.Atoi(parts[0])
			if err != nil || index < 0 || index >= len(current) {
				return value, false
			}
			if len(parts) == 1 {
				return append(current[:index], current[index+1:]...), false
			}
			updated, _ := excludePath(current[index], parts[1:])
			current[index] = updated
		}
	}
	return value, false
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
	if isTerminal || Color {
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
