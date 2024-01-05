package helpers

import (
	"encoding/json"
	"errors"
	"log"

	"net/http"
	"os"

	"server/types"
)

type Envelope map[string]interface{}

type Message struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

var MessageLogs = &Message{
	InfoLog:  infoLog,
	ErrorLog: errorLog,
}

// ReadJSON reads the request body and decodes it into the data interface
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must have only a single json object")
	}

	return nil
}

// WriteJSON writes the data interface to the response writer
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}
	return nil
}

// WriteJSON writes the data interface to the response writer
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload types.JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	WriteJSON(w, statusCode, payload)
}

// `Contains` checks if an element is in an array
func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func InterfaceArrayToStringArray(arr []interface{}) []string {
	var stringArr []string

	for _, v := range arr {
		stringArr = append(stringArr, v.(string))
		// println("InterfaceArrayToStringArray: v=", arr[i])
	}
	return stringArr
}
