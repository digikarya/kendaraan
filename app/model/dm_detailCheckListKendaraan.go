package model

import (
	"errors"
	"github.com/digikarya/kendaraan/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
)


type DetailCheckListKendaraanPayload struct {
	DetailChecklistID 	uint `gorm:"column:detail_checklist_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	Nama 				string `json:"nama"  validate:"required"`
	Tipe 				string `json:"tipe"  validate:"required"`
	CheckListID 		string `json:"check_list_id"  validate:"" `
}

type DetailCheckListKendaraanResponse struct {
	DetailChecklistID 	uint `gorm:"column:detail_checklist_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	Nama 				string `json:"nama"  validate:"required"`
	Tipe 				string `json:"tipe"  validate:"required"`
	CheckListID 		uint `json:"check_list_id"`
}

func (DetailCheckListKendaraanPayload) TableName() string {
	return "detail_check_list"
}
func (DetailCheckListKendaraanResponse) TableName() string {
	return "detail_check_list"
}


func (data *DetailCheckListKendaraanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	if err != nil {
		trx.Rollback()
		return nil,err
	}
	result := trx.Select("nama","tipe","check_list_id").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.DetailChecklistID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *DetailCheckListKendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := DetailCheckListKendaraanResponse{}
	if err := db.Where("detail_checklist_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Select("nama","tipe").Where("detail_checklist_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return tmpUpdate,nil
}


func (data *DetailCheckListKendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("detail_checklist_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("detail_checklist_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================


func (data *DetailCheckListKendaraanPayload) defineValue()  (tmp DetailCheckListKendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	tmp.Tipe = data.Tipe
	//tmp.CheckListID = data.CheckListID
	if len(data.CheckListID) > 0{
		tmp.CheckListID,err = helper.DecodeHash(data.CheckListID)
		if err != nil {
			return tmp,errors.New("data tidak sesuai")
		}
	}

	return tmp,nil
}

func (data *DetailCheckListKendaraanResponse) switchValue(tmp *DetailCheckListKendaraanResponse) {
	// hanya digunakan untuk update
	data.Nama = tmp.Nama
	data.Tipe = tmp.Tipe
	if data.CheckListID != 0 {
		data.CheckListID = tmp.CheckListID
	}
}

func (data *DetailCheckListKendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *DetailCheckListKendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *DetailCheckListKendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&DetailCheckListKendaraanResponse{}).Where("detail_checklist_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *DetailCheckListKendaraanPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	data.HashID = hashID
	response := db.Select("hash_id").Where("detail_checklist_id",id).Save(&data)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah a id")
	}
	return nil
}

func (data *DetailCheckListKendaraanResponse) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	data.HashID = hashID
	response := db.Select("hash_id").Where("detail_checklist_id",id).Save(&data)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah respone id")
	}
	return nil
}
