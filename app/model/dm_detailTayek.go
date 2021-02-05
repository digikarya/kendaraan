package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/digikarya/kendaraan/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type DetailTrayekPayload struct{
		DetailTrayekID  uint `gorm:"column:detail_trayek_id; PRIMARY_KEY" json:"-"`
		HashID 			string `json:"id"  validate:""`
		Nama 			string `json:"nama"  validate:"required"`
		Sequence 		int `json:"sequence"  validate:""`
		AgenID 			string `json:"agen_id"  validate:"required"`
		TrayekID 		string `json:"trayek_id"  validate:"required"`
		NamaDaerah 		string `json:"nama_daerah"  validate:""`
		NamaAgen 		string `json:"nama_agen"  validate:""`
}
type DetailTrayekResponse struct{
	DetailTrayekID  uint `gorm:"column:detail_trayek_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  validate:""`
	Nama 			string `json:"nama"  validate:"required"`
	Sequence 		int `json:"sequence"  validate:""`
	AgenID 			uint `json:"agen_id"  validate:"required"`
	TrayekID 		uint `json:"trayek_id"  validate:"required"`
	NamaDaerah 		string `json:"nama_daerah"  validate:""`
	NamaAgen 		string `json:"nama_agen"  validate:""`
}

func (DetailTrayekPayload) TableName() string {
	return "detail_trayek"
}
func (DetailTrayekResponse) TableName() string {
	return "detail_trayek"
}


func (data *DetailTrayekPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	if err != nil {
		return nil, err
	}
	result := trx.Select("nama","sequence","agen_id","trayek_id","nama_daerah","nama_agen").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.DetailTrayekID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *DetailTrayekPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := DetailTrayekResponse{}
	if err := db.Where("detail_trayek_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("detail_trayek_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}

func (data *DetailTrayekResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("detail_trayek_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}

func (data *DetailTrayekPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("detail_trayek_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("detail_trayek_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *DetailTrayekResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []DetailTrayekResponse
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
		trans = trans.Where("detail_trayek_id > ?",id).Find(&result)
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


func (data *DetailTrayekPayload) defineValue()  (tmp DetailTrayekResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	//tmp.Sequence = data.Sequence
	if err = data.checkOtherService(&tmp,data.AgenID); err != nil {
		return tmp,err
	}
	tmp.AgenID,err = helper.DecodeHash(data.AgenID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.TrayekID,err = helper.DecodeHash(data.TrayekID)
	log.Print(tmp.TrayekID)
	log.Print("trayek")

	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}

func (data *DetailTrayekPayload) checkOtherService(tmp *DetailTrayekResponse,hashID ...string)  (err error) {
	checkAgen := helper.GetEndpoint().Kepegawaian.URL+helper.GetEndpoint().Kepegawaian.Agen+"/"+hashID[0]
	code,responseAgen,err := helper.Curl("GET",checkAgen,nil)
	log.Print(responseAgen)
	if err != nil{
		return err
	}
	if code != http.StatusOK {
		return errors.New("Agen tidak ditemukan ")
	}
	var result map[string]interface{}
	json.Unmarshal(responseAgen, &result)
	if val, ok := result["nama"]; ok {
		get := fmt.Sprintf("%v", val)
		tmp.NamaAgen = get
	}
	if val, ok := result["kecamatan"]; ok {
		get := fmt.Sprintf("%v", val)
		tmp.NamaDaerah = get
	}
	return
}

func (data *DetailTrayekResponse) switchValue(tmp *DetailTrayekResponse) {
	// hanya digunakan untuk update
	//data.NoTrayek = tmp.NoTrayek
	//data.Asal = tmp.Asal
	//data.Tujuan = tmp.Tujuan
}

func (data *DetailTrayekPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *DetailTrayekResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *DetailTrayekPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&DetailTrayekResponse{}).Where("detail_trayek_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *DetailTrayekPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("detail_trayek_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
