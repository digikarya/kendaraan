package handler

import (
	"errors"
	"github.com/digikarya/kendaraan/app/model"
	"github.com/digikarya/helper"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func CheckListKendaraanCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.CheckListKendaraanPayloadMany{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func CheckListKendaraanAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.CheckListKendaraanResponse{}
	hashID,limit := helper.DecodeURLParam(r)
	data,err := serv.All(db,hashID,limit)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}

func CheckListKendaraanFind(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.CheckListKendaraanResponseMany{}
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

func CheckListKendaraanUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.CheckListKendaraanPayload{}
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

func CheckListKendaraanDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.CheckListKendaraanPayload{}
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


func DetailCheckListKendaraanUpdate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.DetailCheckListKendaraanPayload{}
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
func DetailCheckListKendaraanDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.DetailCheckListKendaraanPayload{}
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


func DetailCheckListKendaraanCreate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.DetailCheckListKendaraanPayload{}
	data,err := serv.Create(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Success",http.StatusOK, data)
	return
}



