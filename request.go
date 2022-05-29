package main

type sqlQuery struct {
	Query    string `json:"query"`
	Leniency bool   `json:"field_multi_value_leniency"`
}
