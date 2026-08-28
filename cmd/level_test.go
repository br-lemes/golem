package cmd

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestLevelValidate(t *testing.T) {
	valid := levelFlags{Group: "skill", Skill: []string{"combat"}}
	tests := []struct {
		name    string
		options levelFlags
		wantErr bool
	}{
		{
			name:    "valid skill group",
			options: valid,
		},
		{
			name:    "valid character group",
			options: levelFlags{Group: "character"},
		},
		{
			name:    "invalid group",
			options: levelFlags{Group: "other"},
			wantErr: true,
		},
		{
			name:    "invalid skill",
			options: levelFlags{Group: "skill", Skill: []string{"other"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := levelValidate(tt.options)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("levelValidate() error = %v, want error = %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsBySkill(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{Name: "Ada", Level: 10, MiningLevel: 20},
		{Name: "Bob", Level: 15, MiningLevel: 25},
	}
	got := levelsBySkill(characters, []string{"combat", "mining"})
	want := map[string]map[string]int{
		"combat": {"Ada": 10, "Bob": 15},
		"mining": {"Ada": 20, "Bob": 25},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("levelsBySkill() = %#v, want %#v", got, want)
	}
}

func TestGroupByCharacter(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{Name: "Ada", Level: 10, MiningLevel: 20},
		{Name: "Bob", Level: 15, MiningLevel: 25},
	}
	previousStdout := console.Stdout
	previousFormat := console.Format
	defer func() {
		console.Stdout = previousStdout
		console.Format = previousFormat
	}()
	var output bytes.Buffer
	console.Stdout = &output
	console.Format = "json"
	err := groupByCharacter(characters, []string{"combat", "mining"})
	if err != nil {
		t.Fatalf("groupByCharacter() error = %v", err)
	}
	got := map[string]map[string]int{}
	err = json.Unmarshal(output.Bytes(), &got)
	if err != nil {
		t.Fatalf("decode groupByCharacter() output: %v", err)
	}
	want := map[string]map[string]int{
		"Ada": {"combat": 10, "mining": 20},
		"Bob": {"combat": 15, "mining": 25},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groupByCharacter() = %#v, want %#v", got, want)
	}
}
