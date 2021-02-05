package handler

import (
	"errors"
	"github.com/digikarya/kendaraan/app/model"
	"github.com/digikarya/kendaraan/helper"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func JadwalCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JadwalPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func JadwalAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JadwalResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func JadwalFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JadwalResponse{}
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

func JadwalUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JadwalPayload{}
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

func JadwalDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JadwalPayload{}
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


