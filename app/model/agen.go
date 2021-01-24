package model

import (
	"errors"
	"github.com/digikarya/kepegawaian/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type AgenPayload struct{
		AgenID        	uint `gorm:"column:agen_id; PRIMARY_KEY" json:"-"`
		HashID 			string `json:"-"  validate:""`
		Nama 			string `json:"nama"  validate:"required"`
		Alamat 			string `json:"alamat"  validate:"required"`
		Tipe 			string `json:"tipe"  validate:"required,alpha"`
		NoTlpn 			string `json:"telpon"  validate:"required,number"`
		NoWa			string `json:"whatsapp"  validate:"required,number"`
		DaerahID		string `json:"daerah_id"  validate:"required,base64url"`
}
type AgenResponse struct{
		AgenID        	uint `gorm:"column:agen_id; PRIMARY_KEY" json:"-"`
		HashID 			string `json:"id"  validate:"required,base64url"`
		Nama 			string `json:"nama"  validate:"required"`
		Alamat 			string `json:"alamat"  validate:"required"`
		Tipe 			string `json:"tipe"  validate:"required,alpha"`
		NoTlpn 			string `json:"telpon"  validate:"required,number"`
		NoWa			string `json:"whatsapp"  validate:"required,number"`
		DaerahID		uint `json:"daerah_id"  validate:"required"`
		Kecamatan		string `json:"kecamatan"  validate:""`
		Kabupaten		string `json:"kabupaten"  validate:""`
		Provinsi		string `json:"provinsi"  validate:""`
	}

func (AgenPayload) TableName() string {
	return "agen"
}
func (AgenResponse) TableName() string {
	return "agen"
}


func (data *AgenPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	result := trx.Select("nama","alamat","tipe","no_tlpn","no_wa","daerah_id").Create(&tmp)
	if result.Error != nil {
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.AgenID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *AgenPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := AgenResponse{}
	if err := db.Where("agen_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("agen_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *AgenResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	sql :=  "SELECT " +
		"	 agen.hash_id,agen.nama,agen.alamat,agen.no_tlpn,agen.tipe," +
		"    daerah.daerah_id, daerah.kabupaten, daerah.kecamatan, daerah.provinsi" +
		"	 FROM agen JOIN daerah ON agen.daerah_id=daerah.daerah_id WHERE agen_id = ?"
	result := db.Raw(sql+" LIMIT 1", id).Scan(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}



func (data *AgenPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("agen_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("agen_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *AgenResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []AgenResponse
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	//trans := db.Limit(limit).Find(&result)
	sql :=  "SELECT " +
		"	 agen.hash_id,agen.nama,agen.alamat,agen.no_tlpn,agen.tipe," +
		"    daerah.daerah_id, daerah.kabupaten, daerah.kecamatan, daerah.provinsi" +
		"	 FROM agen JOIN daerah ON agen.daerah_id=daerah.daerah_id"
	hashID := string[0]
	param1 := limit
	param2 := limit
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		sql += " WHERE agen_id > ?"
		param1 = int(id)
		//trans = trans.Where("agen_id > ?",id).Find(&result)
	}
	exec := db.Raw(sql+" LIMIT ?", param1,param2).Scan(&result)
	if exec.Error != nil {
		return result,exec.Error
	}
	return result,nil
}


// General Function =================================================================================
// ==================================================================================================
// ==================================================================================================


func (data *AgenPayload) defineValue()  (tmp AgenResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	tmp.Alamat = data.Alamat
	tmp.Tipe = data.Tipe
	tmp.NoTlpn = data.NoTlpn
	tmp.NoWa = data.NoWa
	tmp.DaerahID,err = helper.DecodeHash(data.DaerahID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}

func (data *AgenResponse) switchValue(tmp *AgenResponse) {
	// hanya digunakan untuk update
	data.Nama = tmp.Nama
	data.Alamat = tmp.Alamat
	data.Tipe = tmp.Tipe
	data.NoTlpn = tmp.NoTlpn
	data.NoWa = tmp.NoWa
	data.DaerahID = tmp.DaerahID
}
func (data *AgenPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *AgenResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *AgenPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&AgenResponse{}).Where("agen_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *AgenPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("agen_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
