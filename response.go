package essqlclient

import (
	"encoding/json"
	"fmt"
)

type version struct {
	Version struct {
		Number *string `json:"number"`
	} `json:"version"`
}

type sqlColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type precookedResponse struct {
	Columns []sqlColumn `json:"columns"`
	Rows    [][]interface{}
}

type SqlResponse struct {
	Columns []sqlColumn
	Rows    []map[string]interface{}
}

func parseJsonResponse(b []byte) (*SqlResponse, error) {
	var sql SqlResponse
	var raw precookedResponse
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return nil, fmt.Errorf("es: response parse: %s", err.Error())
	}

	sql.Columns = raw.Columns

	for _, rawRow := range raw.Rows {
		row := make(map[string]interface{}, len(sql.Columns))
		for k, v := range rawRow {
			row[sql.Columns[k].Name] = v
		}

		sql.Rows = append(sql.Rows, row)
	}

	return &sql, nil
}
