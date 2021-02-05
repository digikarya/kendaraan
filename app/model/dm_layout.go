package model

import (
	"errors"
	"github.com/digikarya/kendaraan/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type LayoutPayload struct{
		LayoutID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
		HashID 		string `json:"-"  validate:""`
		Nama 		string `json:"nama"  validate:"required"`
		BrsKiri 	int `json:"baris_kiri"  validate:"required,numeric,min=0"`
		KlmKiri 	int `json:"kolom_kiri"  validate:"required,numeric,min=0"`
		BrsKanan 	int `json:"baris_kanan"  validate:"required,numeric,min=0"`
		KlmKanan 	int `json:"kolom_kanan"  validate:"required,numeric,min=0"`
		SeatBelakang	*int `json:"seat_belakang"  validate:"required,numeric,min=0"`
		TotalSeat		int `json:"total_seat"  validate:"required,numeric,min=0"`
}
type LayoutResponse struct{
	LayoutID    uint `gorm:"column:layout_id; PRIMARY_KEY" json:"-"`
	HashID 		string `json:"id"  validate:""`
	Nama 		string `json:"nama"  validate:"required"`
	BrsKiri 	int `json:"baris_kiri"  validate:"required,numeric,min=0"`
	KlmKiri 	int `json:"kolom_kiri"  validate:"required,numeric,min=0"`
	BrsKanan 	int `json:"baris_kanan"  validate:"required,numeric,min=0"`
	KlmKanan 	int `json:"kolom_kanan"  validate:"required,numeric,min=0"`
	SeatBelakang	*int `json:"seat_belakang"  validate:"required,numeric,min=0"`
	TotalSeat		int `json:"total_seat"  validate:"numeric,min=0"`
}

func (LayoutPayload) TableName() string {
	return "layout_kursi"
}
func (LayoutResponse) TableName() string {
	return "layout_kursi"
}


func (data *LayoutPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	result := trx.Select("nama","brs_kiri","brs_kanan","klm_kiri","klm_kanan","seat_belakang","total_seat").Create(&tmp)
	if result.Error != nil {
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.LayoutID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *LayoutPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := LayoutResponse{}
	if err := db.Where("layout_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("layout_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}

func (data *LayoutResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("layout_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}

func (data *LayoutPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("layout_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("layout_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *LayoutResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	var result []LayoutResponse
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
		trans = trans.Where("layout_id > ?",id).Find(&result)
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


func (data *LayoutPayload) defineValue()  (tmp LayoutResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	tmp.BrsKiri = data.BrsKiri
	tmp.BrsKanan = data.BrsKanan
	tmp.KlmKiri = data.KlmKiri
	tmp.KlmKanan = data.KlmKanan
	tmp.SeatBelakang = data.SeatBelakang
	tmp.TotalSeat = (tmp.BrsKiri * tmp.KlmKiri) + (tmp.BrsKanan * tmp.KlmKanan) + *tmp.SeatBelakang
	return tmp,nil
}

func (data *LayoutResponse) switchValue(tmp *LayoutResponse) {
	// hanya digunakan untuk update
	data.Nama = tmp.Nama
	data.BrsKiri = tmp.BrsKiri
	data.BrsKanan = tmp.BrsKanan
	data.KlmKiri = tmp.KlmKiri
	data.KlmKanan = tmp.KlmKanan
	data.SeatBelakang = tmp.SeatBelakang
	data.TotalSeat = tmp.TotalSeat
}

func (data *LayoutPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *LayoutResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *LayoutPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&LayoutResponse{}).Where("layout_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *LayoutPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("layout_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
