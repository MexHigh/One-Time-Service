package main

type GenericResponse struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response"`
}
