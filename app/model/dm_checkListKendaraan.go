package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"
)


type CheckListKendaraanPayloadMany struct{
	CheckListID     uint `gorm:"column:check_list_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  validate:""`
	JenisKendaraan 	string `json:"jenis_kendaraan"  validate:"required"`
	Merek	 		string `json:"merek"  validate:"required"`
	Detail 			[]DetailCheckListKendaraanPayload `gorm:"foreignKey:check_list_id;references:check_list_id" json:"detail"  validate:"" `
}
type CheckListKendaraanResponseMany struct{
	CheckListID     uint `gorm:"column:check_list_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  validate:""`
	JenisKendaraan 	string `json:"jenis_kendaraan"  validate:"required"`
	Merek	 		string `json:"merek"  validate:"required"`
	Detail 			[]DetailCheckListKendaraanResponse `gorm:"foreignKey:check_list_id;references:check_list_id" json:"detail"  validate:"" `
}


type CheckListKendaraanPayload struct{
	CheckListID     uint `gorm:"column:check_list_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  validate:""`
	JenisKendaraan 	string `json:"jenis_kendaraan"  validate:"required"`
	Merek	 		string `json:"merek"  validate:"required"`
	//Detail 			[]DetailCheckListKendaraanPayload `json:"detail"  validate:"required" gorm:"foreignKey:check_list_id;references:check_list_id"`
}
type CheckListKendaraanResponse struct{
	CheckListID    	uint `gorm:"column:check_list_id; PRIMARY_KEY" json:"-"`
	HashID 			string `json:"id"  validate:""`
	JenisKendaraan 	string `json:"jenis_kendaraan"  validate:"required"`
	Merek	 		string `json:"merek"  validate:"required"`
	//Detail 			[]DetailCheckListKendaraanPayload `json:"detail"  validate:"required" gorm:"foreignKey:check_list_id;references:check_list_id"`
}

func (CheckListKendaraanPayloadMany) TableName() string {
	return "check_list_kendaraan"
}
func (CheckListKendaraanResponseMany) TableName() string {
	return "check_list_kendaraan"
}

func (CheckListKendaraanPayload) TableName() string {
	return "check_list_kendaraan"
}
func (CheckListKendaraanResponse) TableName() string {
	return "check_list_kendaraan"
}



func (main *CheckListKendaraanPayloadMany) Create(db *gorm.DB,r *http.Request) (interface{},error){
	//data := CheckListKendaraanPayload{}
	err := main.setPayload(r)
	if err != nil {
		return nil, err
	}
	tmp,tmpDetail,err := main.defineValue()
	trx := db.Begin()
	result := trx.Select("merek","jenis_kendaraan").Save(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := tmp.updateHashId(trx,int(tmp.CheckListID));err != nil{
		trx.Rollback()
		return nil, err
	}
	tmpOut := struct {
		CheckListKendaraanResponse
		Detail 			[]DetailCheckListKendaraanResponse `json:"detail"  validate:"required" `
	}{tmp,nil}
	for _,item := range tmpDetail{
		tmpItem := DetailCheckListKendaraanResponse{}
		tmpItem.CheckListID = tmp.CheckListID
		tmpItem.Nama = item.Nama
		tmpItem.Tipe = item.Tipe
		result := trx.Select("nama","tipe","check_list_id").Create(&tmpItem)
		if result.Error != nil {
			trx.Rollback()
			return nil,result.Error
		}
		if result.RowsAffected < 1 {
			trx.Rollback()
			return nil,errors.New("failed to create data")
		}
		//log.Print(tmp.AgenID)
		if err := item.updateHashId(trx,int(tmpItem.DetailChecklistID));err != nil{
			trx.Rollback()
			return nil, err
		}
		tmpOut.Detail = append(tmpOut.Detail,tmpItem)
	}
	trx.Commit()
	return tmpOut,nil
}

func (data *CheckListKendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := CheckListKendaraanResponse{}
	if err := db.Where("check_list_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Select("nama","tipe","check_list_id").Where("jenis_id = ?", id).Updates(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return tmpUpdate,nil
}


func (data *CheckListKendaraanResponseMany) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Preload(clause.Associations).Where("check_list_id",id).Find(&data)
	if result.Error != nil {
		return nil,result.Error
	}

	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}


func (data *CheckListKendaraanResponseMany) FindByKategori(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	kategori := KategoriKendaraanPayload{}
	result := db.Where("kategori_id",id).Find(&kategori)
	if result.Error != nil {
		return nil,result.Error
	}
	result = db.Preload(clause.Associations).Where("check_list_id",kategori.CheckListID).Find(&data)
	if result.Error != nil {
		return nil,result.Error
	}

	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}



func (data *CheckListKendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	trx := db.Begin()
	result := trx.Where("check_list_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		trx.Rollback()
		return nil,err
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("data tidak ditemukan")
	}
	tmpDetail := DetailCheckListKendaraanResponse{}
	response := trx.Where("check_list_id = ?",id).Delete(&tmpDetail)
	if response.Error != nil {
		trx.Rollback()
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	response = trx.Omit(clause.Associations).Where("check_list_id = ?",id).Delete(&data)
	if response.Error != nil {
		trx.Rollback()
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}

	trx.Commit()
	return data,nil
}


func (data *CheckListKendaraanResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	 result := []CheckListKendaraanPayloadMany{}
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
		trans = trans.Where("check_list_id > ?",id).Find(&result)
	}
	exec := trans.Preload(clause.Associations).Find(&result)
	if exec.Error != nil {
		return result,exec.Error
	}
	return result,nil
}


// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================

func (data *CheckListKendaraanPayloadMany) defineValue()  (tmp CheckListKendaraanResponse,tmpDetail []DetailCheckListKendaraanPayload,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.JenisKendaraan = data.JenisKendaraan
	tmp.Merek = data.Merek
	tmpDetail = data.Detail
	return tmp,tmpDetail,nil
}

func (data *CheckListKendaraanPayloadMany) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *CheckListKendaraanPayload) defineValue()  (tmp CheckListKendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.JenisKendaraan = data.JenisKendaraan
	tmp.Merek = data.Merek
	//tmp.Detail = data.Detail
	return tmp,nil
}

func (data *CheckListKendaraanResponse) switchValue(tmp *CheckListKendaraanResponse) {
	// hanya digunakan untuk update
	data.JenisKendaraan = tmp.JenisKendaraan
	data.Merek = tmp.Merek
	//tmp.Detail = data.Detail
}

func (data *CheckListKendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *CheckListKendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *CheckListKendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&CheckListKendaraanResponse{}).Where("check_list_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *CheckListKendaraanResponse) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	data.HashID = hashID
	response := db.Where("check_list_id",id).Save(&data)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
