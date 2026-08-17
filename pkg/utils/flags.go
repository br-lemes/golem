package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterFlags[T any](cmd *cobra.Command) error {
	_, fields, err := flagFields[T]()
	if err != nil {
		return err
	}

	for _, field := range fields {
		name := field.Tag.Get("flag")
		description := field.Tag.Get("desc")
		shorthand := field.Tag.Get("shorthand")
		defaultValue := field.Tag.Get("default")

		switch field.Type.Kind() {
		case reflect.String:
			if shorthand == "" {
				cmd.Flags().String(name, defaultValue, description)
			} else {
				cmd.Flags().StringP(name, shorthand, defaultValue, description)
			}
		case reflect.Bool:
			defaultBool, parseErr := parseBoolDefault(defaultValue, name)
			if parseErr != nil {
				return parseErr
			}
			if shorthand == "" {
				cmd.Flags().Bool(name, defaultBool, description)
			} else {
				cmd.Flags().BoolP(name, shorthand, defaultBool, description)
			}
		case reflect.Int:
			defaultInt, parseErr := parseIntDefault(defaultValue, name)
			if parseErr != nil {
				return parseErr
			}
			if shorthand == "" {
				cmd.Flags().Int(name, defaultInt, description)
			} else {
				cmd.Flags().IntP(name, shorthand, defaultInt, description)
			}
		case reflect.Slice:
			if field.Type != reflect.TypeOf([]string{}) {
				return fmt.Errorf("flag %q has unsupported type %s", name, field.Type)
			}
			values := []string(nil)
			if defaultValue != "" {
				values = strings.Split(defaultValue, ",")
			}
			if shorthand == "" {
				cmd.Flags().StringSlice(name, values, description)
			} else {
				cmd.Flags().StringSliceP(name, shorthand, values, description)
			}
		default:
			return fmt.Errorf("flag %q has unsupported type %s", name, field.Type)
		}
	}
	return nil
}

func ReadFlags[T any](cmd *cobra.Command) (T, error) {
	var result T
	value, fields, err := flagFields(&result)
	if err != nil {
		return result, err
	}
	for _, field := range fields {
		name := field.Tag.Get("flag")
		switch field.Type.Kind() {
		case reflect.String:
			v, err := cmd.Flags().GetString(name)
			if err != nil {
				return result, err
			}
			value.Field(field.Index).SetString(v)
		case reflect.Bool:
			v, err := cmd.Flags().GetBool(name)
			if err != nil {
				return result, err
			}
			value.Field(field.Index).SetBool(v)
		case reflect.Int:
			v, err := cmd.Flags().GetInt(name)
			if err != nil {
				return result, err
			}
			value.Field(field.Index).SetInt(int64(v))
		case reflect.Slice:
			v, err := cmd.Flags().GetStringSlice(name)
			if err != nil {
				return result, err
			}
			value.Field(field.Index).Set(reflect.ValueOf(v))
		}
	}
	return result, nil
}

type flagField struct {
	Type  reflect.Type
	Tag   reflect.StructTag
	Index int
}

func flagFields[T any](target ...*T) (reflect.Value, []flagField, error) {
	typeOf := reflect.TypeOf((*T)(nil)).Elem()
	if typeOf.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("flags type must be a struct")
	}
	var value reflect.Value
	if len(target) > 0 {
		value = reflect.ValueOf(target[0]).Elem()
	}
	fields := []flagField{}
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if field.Tag.Get("flag") == "" {
			continue
		}
		if value.IsValid() && !value.Field(i).CanSet() {
			return reflect.Value{}, nil, fmt.Errorf("flag field %s cannot be set", field.Name)
		}
		fields = append(fields, flagField{field.Type, field.Tag, i})
	}
	return value, fields, nil
}

func parseBoolDefault(value, name string) (bool, error) {
	if value == "" {
		return false, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid default for flag %q: %w", name, err)
	}
	return parsed, nil
}

func parseIntDefault(value, name string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid default for flag %q: %w", name, err)
	}
	return parsed, nil
}
