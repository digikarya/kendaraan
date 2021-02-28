package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type JadwalPayload struct{
		JadwalID    uint `gorm:"column:jadwal_id; PRIMARY_KEY" json:"-"`
		HashID 		string `json:"id"  validate:""`
		WaktuKeberangkatan 	string `json:"waktu_keberangkatan"  validate:"required"`
		WaktuKedataangan 	string `json:"waktu_kedataangan"  validate:"required"`
		TrayekID 	string `json:"trayek_id"  validate:"required"`
		Trayek 		string `json:"trayek"  validate:""`
}
type JadwalResponse struct{
	JadwalID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	WaktuKeberangkatan 	string `json:"waktu_keberangkatan"  validate:"required"`
	WaktuKedataangan 	string `json:"waktu_kedataangan"         vvalidate:"required"`
	TrayekID 	uint `json:"trayek_id"  validate:"required,numeric"`
	Trayek 		string `json:"trayek"  validate:""`
}

func (JadwalPayload) TableName() string {
	return "jadwal"
}
func (JadwalResponse) TableName() string {
	return "jadwal"
}


func (data *JadwalPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	if err != nil {
		return nil,err
	}
	trx := db.Begin()
	result := trx.Select("waktu_keberangkatan","waktu_kedataangan","kendaraan_id","trayek_id").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.JadwalID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *JadwalPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := JadwalResponse{}
	if err := db.Where("jadwal_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("jadwal_id = ?", id).Save(&tmpUpdate)

	result = db.Select("waktu_keberangkatan","waktu_kedataangan","kendaraan_id","trayek_id").Updates(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return tmpUpdate,nil
}


func (data *JadwalResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	tmp := JadwalPayload{}
	sql :=  `SELECT 
			jadwal.*,
			trayek.trayek_id 'trayekid',trayek.hash_id 'trayek_id',CONCAT(trayek.no_trayek,' | ',trayek.asal,' - ',trayek.tujuan) 'trayek'
			FROM jadwal 
			JOIN trayek ON jadwal.trayek_id=trayek.trayek_id 
			WHERE jadwal_id = ?`
	result := db.Raw(sql+" LIMIT 1", id).Scan(&tmp)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmp,nil
}



func (data *JadwalPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("jadwal_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("jadwal_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *JadwalResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	 result :=  []JadwalPayload{}
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	//trans := db.Limit(limit).Find(&result)
	sql :=  `SELECT 
			jadwal.*,
			trayek.trayek_id 'trayekid',trayek.hash_id 'trayek_id',CONCAT(trayek.no_trayek,' | ',trayek.asal,' - ',trayek.tujuan) 'trayek'
			FROM jadwal 
			JOIN trayek ON jadwal.trayek_id=trayek.trayek_id`
	hashID := string[0]
	param1 := limit
	param2 := limit
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		sql += " WHERE jadwal_id > ?"
		param1 = int(id)
		//trans = trans.Where("jadwal_id > ?",id).Find(&result)
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


func (data *JadwalPayload) defineValue()  (tmp JadwalResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	WaktuKeberangkatan, err := time.Parse("15:04:05", data.WaktuKeberangkatan)
	if err != nil {
		return JadwalResponse{}, errors.New("Invalid payload")
	}
	WaktuKedataangan, err := time.Parse("15:04:05", data.WaktuKedataangan)
	if err != nil {
		return JadwalResponse{}, errors.New("Invalid payload")
	}
	tmp.WaktuKeberangkatan = WaktuKeberangkatan.Format("15:04:05")
	tmp.WaktuKedataangan = WaktuKedataangan.Format("15:04:05")
	tmp.TrayekID,err = helper.DecodeHash(data.TrayekID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.TrayekID,err = helper.DecodeHash(data.TrayekID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}

func (data *JadwalResponse) switchValue(tmp *JadwalResponse) {
	data.WaktuKeberangkatan = tmp.WaktuKeberangkatan
	data.WaktuKeberangkatan = tmp.WaktuKeberangkatan
	data.TrayekID = tmp.TrayekID
}

func (data *JadwalPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *JadwalResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *JadwalPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&JadwalResponse{}).Where("jadwal_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *JadwalPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("jadwal_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
