package handler

import (
	"errors"
	"github.com/digikarya/kendaraan/app/model"
	"github.com/digikarya/kendaraan/helper"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func JenisKendaraanCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JenisKendaraanPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func JenisKendaraanAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JenisKendaraanResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func JenisKendaraanFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JenisKendaraanResponse{}
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

func JenisKendaraanUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JenisKendaraanPayload{}
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

func JenisKendaraanDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.JenisKendaraanPayload{}
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


