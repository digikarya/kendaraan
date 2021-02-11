package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"net/http"
)

type SearchRequest struct {
	Condition []struct{
		Column string `json:"column"  validate:"required,alpha"`
		Value string `json:"value"  validate:"required"`
	} `json:"condition"  validate:"required"`
}

func (payload *SearchRequest) LayoutSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []LayoutResponse{}
	result := db.Where("nama LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("nama asc, nama asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}


func (payload *SearchRequest) JenisKendaraanSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []JenisKendaraanResponse{}
	result := db.Where("nama LIKE ?", "%"+payload.Condition[0].Value+"%").Or("kode LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("nama asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}


func (payload *SearchRequest) CheckListSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []CheckListKendaraanResponse{}
	result := db.Where("jenis_kendaraan LIKE ?", "%"+payload.Condition[0].Value+"%").Or("merek LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	result = result.Order("merek asc").Find(&tmpData)
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}

func (payload *SearchRequest) KategoriKendaraanSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	tmpData := []KategoriKendaraanPayload{}
	sql :=  "SELECT " +
		"	kategori_kendaraan.*," +
		"	jenis_kendaraan.hash_id 'jenis_kendaraan_id', concat(jenis_kendaraan.nama,' - ',jenis_kendaraan.kode) AS 'jenis_kendaraan', " +
		"	check_list_kendaraan.hash_id 'check_list_id',concat(check_list_kendaraan.jenis_kendaraan,' - ',check_list_kendaraan.merek) AS 'check_list', " +
		"	layout_kursi.hash_id 'check_list_id', layout_kursi.nama AS 'layout' " +
		"	FROM kategori_kendaraan" +
		"	JOIN jenis_kendaraan ON kategori_kendaraan.jenis_kendaraan_id=jenis_kendaraan.jenis_Id " +
		"	JOIN check_list_kendaraan ON kategori_kendaraan.check_list_id=check_list_kendaraan.check_list_id " +
		"	JOIN layout_kursi ON kategori_kendaraan.layout_kursi_id=layout_kursi.layout_id " +
		" WHERE kategori_kendaraan.kode LIKE ? OR kategori_kendaraan.nama LIKE  ?"
	result := db.Raw(sql,"%"+payload.Condition[0].Value+"%","%"+payload.Condition[0].Value+"%" ).Scan(&tmpData)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}



func (payload *SearchRequest) TrayekSearch(db *gorm.DB,r *http.Request)  (interface{},error) {
	err := payload.setPayload(r)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	var tmpData = []TrayekPayload{}
	result := db.Find(&tmpData)

	if payload.Condition[0].Column == "asal"{
		if len(payload.Condition) < 2 {
			return nil,errors.New("data tidak sesuai")
		}
		result = db.Where("asal LIKE ? AND tujuan LIKE ? ", "%"+payload.Condition[0].Value+"%","%"+payload.Condition[1].Value+"%").Find(&tmpData)
	}else{
		if payload.Condition[0].Column != "kode" {
			return nil,errors.New("data tidak sesuai")
		}
		result = db.Where("no_trayek LIKE ?", "%"+payload.Condition[0].Value+"%").Find(&tmpData)
	}
	if err := result.Error; err != nil {
		return nil,errors.New("data tidak ditemukan")
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmpData,nil
}

func (payload *SearchRequest) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&payload);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(payload);err != nil {
		return err
	}
	if len(payload.Condition) < 1 {
		return errors.New("invalid payload")
	}
	return nil
}