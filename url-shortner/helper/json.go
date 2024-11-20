package helper

import (
	"encoding/json"
	"net/http"
)

func ReadRequest(r *http.Request, w http.ResponseWriter) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&w)
	if err != nil {
		panic(err)
	}

}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	encode := json.NewEncoder(w)
	err := encode.Encode(data)
	ErrorPanic(err)
}
