package handler

import (

	"github.com/digikarya/kendaraan/helper"
	"github.com/digikarya/kendaraan/app/model"

	"gorm.io/gorm"
	"net/http"
)

func SearchLayout(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.SearchRequest{}
	data,err := serv.LayoutSearch(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func SearchJenisKendaraan(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.SearchRequest{}
	data,err := serv.JenisKendaraanSearch(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}

func SearchKategoriKendaraan(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	serv := model.SearchRequest{}
	data,err := serv.KategoriKendaraanSearch(db,r)
	if err != nil {
		helper.RespondJSONError(w, http.StatusBadRequest, err)
		return
	}
	helper.RespondJSON(w, "Found",http.StatusOK, data)
	return
}



//func SearchAgen(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
//	serv := model.SearchRequest{}
//	data,err := serv.AgenSearch(db,r)
//	if err != nil {
//		helper.RespondJSONError(w, http.StatusBadRequest, err)
//		return
//	}
//	helper.RespondJSON(w, "Found",http.StatusOK, data)
//	return
//}