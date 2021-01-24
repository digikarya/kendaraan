package handler

import (
	"errors"
	"github.com/digikarya/kepegawaian/app/model"
	"github.com/digikarya/kepegawaian/helper"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func AgenCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.AgenPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func AgenAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.AgenResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func AgenFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.AgenResponse{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}

	data,err := serv.Find(db,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func AgenUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.AgenPayload{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.Update(db,r,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondJSON(w, "Updated",http.StatusOK, data)
	return
}

func AgenDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.AgenPayload{}
	vars := mux.Vars(r)
	hashid,exist := vars["hashid"]
	if !exist{
		helper.RespondJSONError(w, http.StatusBadRequest, errors.New("Invalid id"))
		return
	}
	data,err := serv.Delete(db,hashid)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}

	helper.RespondJSON(w, "Deleted",http.StatusOK, data)
	return
}


