package handler

import (
	"errors"
	"github.com/digikarya/kendaraan/app/model"
	"github.com/digikarya/helper"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func KendaraanCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KendaraanPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func KendaraanAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KendaraanResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func KendaraanFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KendaraanResponse{}
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

func KendaraanUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KendaraanPayload{}
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

func KendaraanDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.KendaraanPayload{}
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


