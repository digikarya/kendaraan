package helper

import (
	"encoding/json"
	"net/http"
)

type defaultResponse struct {
	Status struct{
		Code int
		Desc string
	}
	Data interface{}

}
func RespondJSON(w http.ResponseWriter,msg string, status int, payload interface{}) {
	res := defaultResponse{}
	res.Status.Code = status
	res.Status.Desc = msg
	res.Data = payload
	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_,_ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_,_ =w.Write(response)
}

func RespondJSONError(w http.ResponseWriter, status int,err error) {
	res := defaultResponse{}
	res.Status.Code = status
	res.Status.Desc = err.Error()
	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_,_ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_,_ = w.Write(response)
}


func CorsHelper(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,DELETE,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token,Authorization")
	w.Header().Set("Access-Control-Max-Age", "7200")

}