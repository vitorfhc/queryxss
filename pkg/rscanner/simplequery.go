package rscanner

import (
	"fmt"
	"io"
)

type SimpleQuery struct {
	Info ScannerInfo
}

func NewSimpleQuery() *SimpleQuery {
	return &SimpleQuery{
		Info: ScannerInfo{
			ID:             "simple-query",
			Description:    "SimpleQuery scanner requests the URL as it is, with no modifications. It looks for query keys and values being reflected.",
			Severity:       "info",
			SuccessMessage: "simple-query detected a reflected query",
		},
	}
}

func (s *SimpleQuery) Scan(client HttpRequester, input *Input) (*Output, error) {
	resp, err := client.Get(input.Url)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}
	bodytext := string(body)

	query, err := GetURLQuery(input.Url)
	if err != nil {
		return nil, fmt.Errorf("error getting query: %v", err)
	}

	if len(query) == 0 {
		return &Output{
			Success: false,
			Info:    &s.Info,
			Url:     input.Url,
		}, nil
	}

	toFind := []string{}
	for key, values := range query {
		toFind = append(toFind, key)
		toFind = append(toFind, values...)
	}

	headers := resp.Header
	found := false
	for key, values := range headers {
		found = found || CaseInsensitiveContains(key, toFind)
		if found {
			break
		}
		for _, value := range values {
			found = found || CaseInsensitiveContains(value, toFind)
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	found = found || CaseInsensitiveContains(bodytext, toFind)

	return &Output{
		Success: found,
		Info:    &s.Info,
		Url:     input.Url,
	}, nil
}

func (s *SimpleQuery) GetInfo() *ScannerInfo {
	return &s.Info
}
