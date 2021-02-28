package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type TrayekPayloadMany struct{
	TrayekID    uint `gorm:"column:trayek_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	NoTrayek 	string `json:"no_trayek"  validate:"required"`
	Asal 	string `json:"asal"  validate:"required"`
	Tujuan 	string `json:"tujuan"  validate:"required"`
	Detail 			interface{} `gorm:"foreignKey:trayek_id;references:trayek_id" json:"detail"  validate:"" `
}
type TrayekPayload struct{
		TrayekID    uint `gorm:"column:trayek_id; PRIMARY_KEY" json:"-"`
		HashID 		string `json:"id"  validate:""`
		NoTrayek 	string `json:"no_trayek"  validate:"required"`
		Asal 	string `json:"asal"  validate:"required"`
		Tujuan 	string `json:"tujuan"  validate:"required"`
}
type TrayekResponse struct{
	TrayekID    uint `gorm:"column:trayek_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	NoTrayek 	string `json:"no_trayek"  validate:"required"`
	Asal 	string `json:"asal"  validate:"required"`
	Tujuan 	string `json:"tujuan"  validate:"required"`
}

func (TrayekPayload) TableName() string {
	return "trayek"
}
func (TrayekResponse) TableName() string {
	return "trayek"
}
func (TrayekPayloadMany) TableName() string {
	return "trayek"
}


func (data *TrayekPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	result := trx.Select("no_trayek","asal","tujuan").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.TrayekID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *TrayekPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := TrayekResponse{}
	if err := db.Where("trayek_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Select("no_trayek","asal","tujuan").Where("trayek_id = ?", id).Updates(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return tmpUpdate,nil
}

func (data *TrayekResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("trayek_id",id).Find(&data)
	if result.Error != nil {
		return nil,result.Error
	}

	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	tmp := TrayekPayloadMany{}
	tmp.TrayekID = data.TrayekID
	tmp.HashID = data.HashID
	tmp.NoTrayek = data.NoTrayek
	tmp.Asal = data.Asal
	tmp.Tujuan = data.Tujuan
	tmpDetail := DetailTrayekResponse{}
	tmp.Detail,err = tmpDetail.Find(db,tmp.HashID)
	if err != nil {
		tmp.Detail = []DetailTrayekPayload{}
	}
	return tmp,nil
}

func (data *TrayekPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("trayek_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("trayek_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *TrayekResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	result := []TrayekResponse{}
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
		trans = trans.Where("trayek_id > ?",id).Find(&result)
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


func (data *TrayekPayload) defineValue()  (tmp TrayekResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.NoTrayek = data.NoTrayek
	tmp.Asal = data.Asal
	tmp.Tujuan = data.Tujuan
	return tmp,nil
}

func (data *TrayekResponse) switchValue(tmp *TrayekResponse) {
	// hanya digunakan untuk update
	data.NoTrayek = tmp.NoTrayek
	data.Asal = tmp.Asal
	data.Tujuan = tmp.Tujuan
}

func (data *TrayekPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *TrayekResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *TrayekPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&TrayekResponse{}).Where("trayek_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *TrayekPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("trayek_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
