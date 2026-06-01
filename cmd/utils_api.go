package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

const baseURL = "https://api.artifactsmmo.com"

var httpClient = &http.Client{Timeout: 30 * time.Second}

func apiGet(path string, data map[string]string) ([]byte, error) {
	if len(data) > 0 {
		values := url.Values{}
		for key, val := range data {
			if val != "" {
				values.Add(key, val)
			}
		}
		queryString := values.Encode()
		if queryString != "" {
			path = path + "?" + queryString
		}
	}
	return apiRequest(http.MethodGet, path, nil)
}

func apiPost(path string, data any) ([]byte, error) {
	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return apiRequest(http.MethodPost, path, body)
}

func apiRequest(method, path string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	if debugFlag {
		fmt.Fprintf(os.Stderr, "→ %s %s\n", method, path)
		if len(body) > 0 {
			fmt.Fprintf(os.Stderr, "  Body: %s\n", string(body))
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errMsg := gjson.GetBytes(respBytes, "error.message")
	if errMsg.Exists() {
		return nil, fmt.Errorf("%s", errMsg.String())
	}

	cdResult := gjson.GetBytes(respBytes, "data.cooldown.total_seconds")
	cd := int(cdResult.Int())

	if cd > 0 {
		fmt.Fprintf(os.Stderr, "⏳ Cooldown started: %d seconds\n", cd)
		time.Sleep(time.Duration(cd) * time.Second)
	}

	return respBytes, nil
}

func apiCharacters(name string) (CharacterSchema, error) {
	resp, err := apiGet("/characters/"+name, nil)
	if err != nil {
		return CharacterSchema{}, err
	}
	var data CharacterResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return CharacterSchema{}, err
	}
	return data.Data, nil
}

func apiActionMove(name string, x, y int) (CharacterMovementDataSchema, error) {
	resp, err := apiPost("/my/"+name+"/action/move", map[string]int{
		"x": x,
		"y": y,
	})
	if err != nil {
		return CharacterMovementDataSchema{}, err
	}
	var data CharacterMovementResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return CharacterMovementDataSchema{}, err
	}
	return data.Data, nil
}

func apiActionBankDepositItem(name string, items []SimpleItemSchema) (BankItemTransactionSchema, error) {
	resp, err := apiPost("/my/"+name+"/action/bank/deposit/item", items)
	if err != nil {
		return BankItemTransactionSchema{}, err
	}
	var data BankItemTransactionResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return BankItemTransactionSchema{}, err
	}
	return data.Data, nil
}

func apiActionGathering(name string) (SkillDataSchema, error) {
	resp, err := apiPost("/my/"+name+"/action/gathering", nil)
	if err != nil {
		return SkillDataSchema{}, err
	}
	var data SkillResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return SkillDataSchema{}, err
	}
	return data.Data, nil
}

func apiActionRest(name string) (CharacterRestDataSchema, error) {
	resp, err := apiPost("/my/"+name+"/action/rest", nil)
	if err != nil {
		return CharacterRestDataSchema{}, err
	}
	var data CharacterRestResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return CharacterRestDataSchema{}, err
	}
	return data.Data, nil
}

func apiActionFight(name string, participants []string) (CharacterFightDataSchema, error) {
	resp, err := apiPost("/my/"+name+"/action/fight", FightRequestSchema{
		Participants: &participants,
	})
	if err != nil {
		return CharacterFightDataSchema{}, err
	}
	var data CharacterFightResponseSchema
	if err := json.Unmarshal(resp, &data); err != nil {
		return CharacterFightDataSchema{}, err
	}
	return data.Data, nil
}
