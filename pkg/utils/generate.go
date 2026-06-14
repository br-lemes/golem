package utils

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"regexp"
	"strings"
	"text/template"
)

//go:embed generate.tmpl
var generate string

type TemplateParam struct {
	Name        string
	Description string
	Type        string
	GoType      string
}

type TemplateRoute struct {
	Method       string
	Path         string
	GoPathFormat string
	PathArgs     []string
}

type TemplateData struct {
	CommandName string
	Summary     string
	Long        string
	Method      string
	MaxArgs     int
	HasBody     bool
	IsBodyArray bool
	ArgUse      string
	Routes      map[int]TemplateRoute
	Flags       []TemplateParam
}

type flagCollector struct {
	targetCmd        string
	flagsCheck       map[string]TemplateParam
	pathDescriptions map[string]string
	flags            []TemplateParam
}

func BuildTemplateData(targetCmd string) (TemplateData, error) {
	var data TemplateData
	data.Routes = make(map[int]TemplateRoute)

	availableCommands, err := GetCommands()
	if err != nil {
		return data, err
	}

	routes, exists := availableCommands[targetCmd]
	if !exists {
		return data, fmt.Errorf("command '%s' not found", targetCmd)
	}

	data.CommandName = targetCmd

	collector := flagCollector{
		targetCmd:        targetCmd,
		flagsCheck:       make(map[string]TemplateParam),
		pathDescriptions: make(map[string]string),
	}

	for _, routeMap := range routes {
		for method, path := range routeMap {
			routeData, err := GetRoute(path)
			if err != nil {
				return data, err
			}

			details := routeData[method]

			if data.Summary == "" {
				data.Summary = details.Summary
			}
			if data.Method == "" {
				data.Method = method
			}

			if details.IsBodyArray {
				data.IsBodyArray = true
			}

			pathArgs, goPathFormat := extractPathArgs(path)
			argCount := len(pathArgs)

			existingRoute, exists := data.Routes[argCount]
			if exists {
				return data, fmt.Errorf(
					"route conflict: both %s %s and %s %s require %d positional arguments",
					existingRoute.Method, existingRoute.Path, method, path, argCount)
			}

			tRoute := TemplateRoute{
				Method:       method,
				Path:         path,
				GoPathFormat: goPathFormat,
				PathArgs:     pathArgs,
			}
			data.Routes[argCount] = tRoute

			currentPathCount := argCount
			for _, param := range details.Parameters {
				if param.In == "path" {
					collector.pathDescriptions[param.Name] = param.Description
				}
				if param.In == "query" {
					err := collector.collectFlag(param)
					if err != nil {
						return data, err
					}
				}
			}

			if currentPathCount > data.MaxArgs {
				data.MaxArgs = currentPathCount
			}

			if len(details.RequestBody) > 0 {
				data.HasBody = true
				for _, body := range details.RequestBody {
					err := collector.collectFlag(body)
					if err != nil {
						return data, err
					}
				}
			}
		}
	}

	data.Flags = collector.flags
	data.formatArguments(collector.pathDescriptions)

	return data, nil
}

func (fc *flagCollector) collectFlag(field FieldInfo) error {
	goType := mapGoType(field.Type)
	existingFlag, hasFlag := fc.flagsCheck[field.Name]

	if !hasFlag {
		flag := TemplateParam{
			Name:        field.Name,
			Description: field.Description,
			Type:        field.Type,
			GoType:      goType,
		}
		fc.flagsCheck[field.Name] = flag
		fc.flags = append(fc.flags, flag)
		return nil
	}

	if existingFlag.Type == field.Type {
		return nil
	}

	return fmt.Errorf("flag '%s' type conflict: %s vs %s",
		field.Name, existingFlag.Type, field.Type)
}

func (d *TemplateData) formatArguments(descriptions map[string]string) {
	if len(d.Routes) == 0 {
		return
	}

	var syntaxParts []string
	var longArgsLines []string
	seenArgs := make(map[string]bool)
	maxArgLength := 0

	for _, r := range d.Routes {
		if len(r.PathArgs) > 0 {
			var formatParts []string
			for _, arg := range r.PathArgs {
				formatParts = append(formatParts, arg)
				if len(arg) > maxArgLength {
					maxArgLength = len(arg)
				}
			}
			syntaxParts = append(syntaxParts, strings.Join(formatParts, " "))
		}
	}

	for _, r := range d.Routes {
		if len(r.PathArgs) > 0 {
			for _, arg := range r.PathArgs {
				if !seenArgs[arg] {
					seenArgs[arg] = true
					desc := descriptions[arg]
					padding := strings.Repeat(" ", maxArgLength-len(arg))
					longArgsLines = append(
						longArgsLines,
						fmt.Sprintf("  %s%s   %s", arg, padding, desc),
					)
				}
			}
		}
	}

	if len(syntaxParts) > 0 {
		_, acceptsZeroArgs := d.Routes[0]
		if acceptsZeroArgs {
			d.ArgUse = fmt.Sprintf("[%s]", strings.Join(syntaxParts, " | "))
		} else {
			d.ArgUse = fmt.Sprintf("<%s>", strings.Join(syntaxParts, " | "))
		}
	}

	if len(longArgsLines) > 0 {
		longArgs := strings.Join(longArgsLines, "\n")
		d.Long = fmt.Sprintf("%s\n\nArguments:\n%s", d.Summary, longArgs)
	}
}

func RenderTemplate(data TemplateData) ([]byte, error) {
	tmpl, err := template.New("command").Parse(generate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format source code: %w", err)
	}

	return formattedCode, nil
}

func mapGoType(apiType string) string {
	if apiType == "integer" {
		return "Int"
	}
	if apiType == "boolean" {
		return "Bool"
	}
	return "String"
}

func extractPathArgs(path string) ([]string, string) {
	pathParamPattern := regexp.MustCompile(`\{([^}]+)\}`)
	var pathArgs []string
	matches := pathParamPattern.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		pathArgs = append(pathArgs, match[1])
	}
	goPathFormat := pathParamPattern.ReplaceAllString(path, "%s")
	return pathArgs, goPathFormat
}
