package models

type CommonRes struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}
