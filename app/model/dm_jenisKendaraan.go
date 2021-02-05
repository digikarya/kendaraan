package model

import (
	"errors"
	"github.com/digikarya/kendaraan/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type JenisKendaraanPayload struct{
	JenisID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	Nama 	string `json:"nama"  validate:"required"`
	Kode 	string `json:"kode"  validate:"required"`
	Jenis 	string `json:"jenis"  validate:"required"`
}
type JenisKendaraanResponse struct{
	JenisID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	Nama 	string `json:"nama"  validate:"required"`
	Kode 	string `json:"kode"  validate:"required"`
	Jenis 	string `json:"jenis"  validate:"required"`
}

func (JenisKendaraanPayload) TableName() string {
	return "jenis_kendaraan"
}
func (JenisKendaraanResponse) TableName() string {
	return "jenis_kendaraan"
}


func (data *JenisKendaraanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	result := trx.Select("nama","kode","jenis").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.JenisID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *JenisKendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}

	if err := data.setPayload(r);err != nil {
		return nil, err
	}
	if _,err := data.countData(db,id);err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	tmpUpdate := JenisKendaraanResponse{}
	if err := db.Where("jenis_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("jenis_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *JenisKendaraanResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("jenis_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}


func (data *JenisKendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("jenis_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("jenis_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *JenisKendaraanResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []JenisKendaraanResponse
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	trans := db.Limit(limit).Find(&result)
	hashID := string[0]
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		trans = trans.Where("jenis_id > ?",id).Find(&result)
	}
	exec := trans.Find(&result)
	if exec.Error != nil {
		return result,exec.Error
	}
	return result,nil
}


// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================


func (data *JenisKendaraanPayload) defineValue()  (tmp JenisKendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	tmp.Kode = data.Kode
	tmp.Jenis = data.Jenis
	return tmp,nil
}

func (data *JenisKendaraanResponse) switchValue(tmp *JenisKendaraanResponse) {
	// hanya digunakan untuk update
	data.Nama = tmp.Nama
	data.Kode = tmp.Kode
	data.Jenis = tmp.Jenis
}

func (data *JenisKendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *JenisKendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *JenisKendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&JenisKendaraanResponse{}).Where("jenis_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *JenisKendaraanPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("jenis_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
