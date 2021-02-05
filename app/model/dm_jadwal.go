package model

import (
	"errors"
	"github.com/digikarya/kendaraan/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type JadwalPayload struct{
		JadwalID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
		HashID 		string `json:"id"  validate:""`
		WaktuKeberangkatan 	int `json:"waktu_keberangkatan"  validate:"required"`
		WaktuKedataangan 	int `json:"waktu_kedataangan"  validate:"required"`
		TrayekID 	string `json:"trayek_id"  validate:"required"`
}
type JadwalResponse struct{
	JadwalID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	WaktuKeberangkatan 	int `json:"waktu_keberangkatan"  validate:"required"`
	WaktuKedataangan 	int `json:"waktu_kedataangan"  validate:"required"`
	TrayekID 	uint `json:"trayek_id"  validate:"required,numeric"`
}

func (JadwalPayload) TableName() string {
	return "jadwal"
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
	trx := db.Begin()
	tmp,err := data.defineValue()
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
	result := db.Where("jadwal_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
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
	var result []JadwalResponse
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
	tmp.WaktuKeberangkatan = data.WaktuKeberangkatan
	tmp.WaktuKeberangkatan = data.WaktuKeberangkatan
	//tmp.KendaraanID,err = helper.DecodeHash(data.KendaraanID)
	//if err != nil {
	//	return tmp,errors.New("data tidak sesuai")
	//}
	tmp.TrayekID,err = helper.DecodeHash(data.TrayekID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}

func (data *JadwalResponse) switchValue(tmp *JadwalResponse) {
	// hanya digunakan untuk update
	data.WaktuKeberangkatan = tmp.WaktuKeberangkatan
	data.WaktuKeberangkatan = tmp.WaktuKeberangkatan
	//data.KendaraanID = tmp.KendaraanID
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
