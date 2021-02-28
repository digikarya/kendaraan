package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type KategoriKendaraanPayload struct{
		KategoriID    		uint `gorm:"column:kategori_id; PRIMARY_KEY" json:"-"`
		HashID 				string `json:"id"  validate:""`
		Nama 				string `json:"nama"  validate:"required"`
		Kode 				string `json:"kode"  validate:"required"`
		CheckListID 		string `json:"check_list_id"  validate:"required"`
		LayoutKursiID 		string `json:"layout_id"  validate:"required"`
		JenisKendaraanID	string `json:"jenis_kendaraan_id"  validate:"required"`
		Kapasitas			uint `json:"kapasitas"  validate:""`
		JenisKendaraan		string `json:"jenis_kendaraan"  validate:""`
		CheckList   		string `json:"check_list"  validate:""`
		Layout		   		string `json:"layout"  validate:""`

}
type KategoriKendaraanResponse struct{
	KategoriID    		uint `gorm:"column:kategori_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	Nama 				string `json:"nama"  validate:"required"`
	Kode 				string `json:"kode"  validate:"required"`
	CheckListID 		uint `json:"check_list_id"  validate:"required,numeric"`
	LayoutKursiID 		uint `json:"layout_id"  validate:"required,numeric"`
	JenisKendaraanID	uint `json:"jenis_kendaraan_id"  validate:"required,numeric"`
	Kapasitas			int `json:"kapasitas"  validate:""`
	JenisKendaraan		string `json:"jenis_kendaraan"  validate:""`
	CheckList   		string `json:"check_list"  validate:""`
	Layout		   		string `json:"layout"  validate:""`
}

func (KategoriKendaraanPayload) TableName() string {
	return "kategori_kendaraan"
}
func (KategoriKendaraanResponse) TableName() string {
	return "kategori_kendaraan"
}


func (data *KategoriKendaraanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	getLayout := LayoutResponse{}
	tmpGetKapasitas,err := getLayout.Find(trx,data.LayoutKursiID)
	if err != nil {
		trx.Rollback()
		return nil, errors.New("Data gagal ditambahkan")
	}
	tmp.Kapasitas = tmpGetKapasitas.(*LayoutResponse).TotalSeat
	result := trx.Select("nama","kode","check_list_id","layout_kursi_id","jenis_kendaraan_id","kapasitas").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}

	if err := data.updateHashId(trx,int(tmp.KategoriID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *KategoriKendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := KategoriKendaraanResponse{}
	result := db.Where("kategori_id = ?", id).First(&tmpUpdate)
	if  result.Error != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result = db.Select("nama","kode","check_list_id","layout_kursi_id","jenis_kendaraan_id","kapasitas").Where("kategori_id = ?", id).Updates(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return tmpUpdate,nil
}


func (data *KategoriKendaraanPayload) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	sql :=  "SELECT " +
		"	kategori_kendaraan.*," +
		"	jenis_kendaraan.hash_id 'jenis_kendaraan_id', concat(jenis_kendaraan.nama,' - ',jenis_kendaraan.kode) AS 'jenis_kendaraan', " +
		"	check_list_kendaraan.hash_id 'check_list_id',concat(check_list_kendaraan.jenis_kendaraan,' - ',check_list_kendaraan.merek) AS 'check_list', " +
		"	layout_kursi.hash_id 'check_list_id', layout_kursi.nama AS 'layout' " +
		"	FROM kategori_kendaraan" +
		"	JOIN jenis_kendaraan ON kategori_kendaraan.jenis_kendaraan_id=jenis_kendaraan.jenis_Id " +
		"	JOIN check_list_kendaraan ON kategori_kendaraan.check_list_id=check_list_kendaraan.check_list_id " +
		"	JOIN layout_kursi ON kategori_kendaraan.layout_kursi_id=layout_kursi.layout_id " +
		" WHERE kategori_id = ?"
	result := db.Raw(sql+" LIMIT 1", id).Scan(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}



func (data *KategoriKendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("kategori_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("kategori_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *KategoriKendaraanPayload) All(db *gorm.DB,string ...string) (interface{}, error) {
 	result :=  []KategoriKendaraanPayload{}
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	//trans := db.Limit(limit).Find(&result)
	sql :=  "SELECT " +
		"	kategori_kendaraan.*," +
		"	jenis_kendaraan.hash_id 'jenis_kendaraan_id', concat(jenis_kendaraan.nama,' - ',jenis_kendaraan.kode) AS 'jenis_kendaraan', " +
		"	check_list_kendaraan.hash_id 'check_list_id',concat(check_list_kendaraan.jenis_kendaraan,' - ',check_list_kendaraan.merek) AS 'check_list', " +
		"	layout_kursi.hash_id 'check_list_id', layout_kursi.nama AS 'layout' " +
		"	FROM kategori_kendaraan" +
		"	JOIN jenis_kendaraan ON kategori_kendaraan.jenis_kendaraan_id=jenis_kendaraan.jenis_Id " +
		"	JOIN check_list_kendaraan ON kategori_kendaraan.check_list_id=check_list_kendaraan.check_list_id " +
		"	JOIN layout_kursi ON kategori_kendaraan.layout_kursi_id=layout_kursi.layout_id "
		hashID := string[0]
	param1 := limit
	param2 := limit
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		sql += " WHERE kategori_id > ?"
		param1 = int(id)
		//trans = trans.Where("kategori_id > ?",id).Find(&result)
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


func (data *KategoriKendaraanPayload) defineValue()  (tmp KategoriKendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.Nama = data.Nama
	tmp.Kode = data.Kode
	tmp.CheckListID,err = helper.DecodeHash(data.CheckListID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.JenisKendaraanID,err = helper.DecodeHash(data.JenisKendaraanID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.LayoutKursiID,err = helper.DecodeHash(data.LayoutKursiID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}

	return tmp,nil
}

func (data *KategoriKendaraanResponse) switchValue(tmp *KategoriKendaraanResponse) {
	// hanya digunakan untuk update
	data.Nama = tmp.Nama
	data.Kode = tmp.Kode
	data.LayoutKursiID = tmp.LayoutKursiID
	data.KategoriID = tmp.KategoriID
	data.JenisKendaraanID = tmp.JenisKendaraanID
}

func (data *KategoriKendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *KategoriKendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *KategoriKendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&KategoriKendaraanResponse{}).Where("kategori_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *KategoriKendaraanPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("kategori_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
