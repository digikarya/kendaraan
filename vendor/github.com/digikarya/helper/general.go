package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type TimeModel struct {
	Created_at    time.Time
	Updated_at    time.Time
	Deleted_at    gorm.DeletedAt `json:"-" `
}

func (r *TimeModel) GenerateTime(isUpdate bool)  {
	timeGenerate := GenerateTime()
	if isUpdate {
		r.Updated_at = timeGenerate
	}else {
		r.Created_at = timeGenerate
		r.Updated_at = r.Created_at

	}
}
func GenerateTime() time.Time {
	loc,_ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}
func DecodeHash(string2 string)  (uint,error){
	decode, err := base64.URLEncoding.DecodeString(string2)
	if err != nil {
		return 0,err
	}
	id,err := Decrypt(decode)
	if err != nil {
		return 0,errors.New("Invalid ID")
	}
	idint,err := strconv.Atoi(string(id))
	if err != nil {
		return 0,err
	}
	return uint(idint),nil
}

func EncodeHash(id int)  (string,error){
	hash,err := Encrypt(strconv.Itoa(id))
	if err != nil {
		return "",errors.New("Falied to create data")
	}

	return base64.URLEncoding.EncodeToString(hash),nil
}

func ValidateData(data interface{}) error  {
	var validate *validator.Validate
	validate = validator.New()
	err := validate.Struct(data)
	if err != nil {
		return err
	}
	return nil
}
func DecodeJson(r *http.Request,payload interface{})  error {
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func DecodePaginate(r *http.Request) (string,string)  {
	hashid := ""
	url := r.URL.Query()
	query,ok := url["id"]
	if ok{
		hashid = query[0]
	}
	limit := "10"
	return hashid,limit
}

func DecodeHashID(r *http.Request) (string){
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist {
		return ""
	}
	return hashid
}

func DecodeURLParam(r *http.Request) (string,string) {
	limit, ok := r.URL.Query()["limit"]
	if !ok || len(limit[0]) < 1 {
		return "","10"
	}

	hashid, ok := r.URL.Query()["hashid"]
	if !ok || len(hashid[0]) < 1 {
		return "",limit[0]
	}
	return hashid[0],limit[0]
}

