package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type SuratKendaraanPayload struct{
		SuratKendaraanID    uint `gorm:"column:surat_id; PRIMARY_KEY" json:"-"`
		HashID 				string `json:"id"  validate:""`
		NoSurat 			string `json:"no_surat"  validate:"required"`
		JenisSurat 			string `json:"jenis_surat"  validate:"required"`
		Kadaluarsa 			string `json:"kadaluarsa"  validate:"required"`
		TanggalPembuatan 	string  `json:"tanggal_pembuatan" validate:"required"`
		KendaraanID 		string  `json:"kendaraan_id" validate:"required"`
		Status 				string `json:"status"  validate:""`
		Img 				string `json:"image"  validate:""`
		helper.TimeModel

}
type SuratKendaraanResponse struct{
	SuratKendaraanID    uint `gorm:"column:surat_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	NoSurat 			string `json:"no_surat"  validate:"required"`
	JenisSurat 			string `json:"jenis_surat"  validate:"required"`
	Kadaluarsa 			string `json:"kadaluarsa"  validate:"required"`
	TanggalPembuatan 	string  `json:"tanggal_pembuatan" validate:"required"`
	KendaraanID 		uint  `json:"kendaraan_id" validate:"required"`
	Status 				string `json:"status"  validate:""`
	Img 				string `json:"image"  validate:""`
	helper.TimeModel
}

func (SuratKendaraanPayload) TableName() string {
	return "surat_kendaraan"
}
func (SuratKendaraanResponse) TableName() string {
	return "surat_kendaraan"
}


func (data *SuratKendaraanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	tmp,err := data.defineValue()
	if err != nil {
		return nil,err
	}
	result := trx.Select( "no_surat", "kadaluarsa", "tanggal_pembuatan", "jenis_surat", "keterangan", "kendaraan_id","created_at","update_at").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("failed to create data")
	}
	if err := data.updateHashId(trx,r,int(tmp.SuratKendaraanID), strconv.Itoa(int(tmp.KendaraanID)),tmp.JenisSurat);err != nil{
		trx.Rollback()
		return nil, err
	}

	trx.Commit()
	return data,nil
}

func (data *SuratKendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := SuratKendaraanResponse{}
	if err := db.Where("surat_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Where("surat_id = ?", id).Save(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *SuratKendaraanResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	sql :=  "SELECT " +
		"	surat_kendaraan.*,kendaraan.kendaraan_id, kendaraan.hash_id 'kendaraan_id' " +
		"		FROM surat_kendaraan " +
		"	JOIN kendaraan ON surat_kendaraan.kendaraan_id=kendaraan.kendaraan_id WHERE surat_id = ?"
	result := db.Raw(sql+" LIMIT 1", id).Scan(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return data,nil
}


func (data *SuratKendaraanResponse) FindByKendaraanAll(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	tmp := []SuratKendaraanPayload{}
	sql :=  "SELECT " +
		"	surat_kendaraan.*,kendaraan.kendaraan_id, kendaraan.hash_id 'kendaraan_id' " +
		"		FROM surat_kendaraan " +
		"	JOIN kendaraan ON surat_kendaraan.kendaraan_id=kendaraan.kendaraan_id WHERE surat_kendaraan.status <> 'aktif' AND surat_kendaraan.kendaraan_id = ? ORDER BY status, created_at"
	result := db.Raw(sql, id).Scan(&tmp)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmp,nil
}


func (data *SuratKendaraanResponse) FindByKendaraanActive(db *gorm.DB,string ...string) (interface{},error){
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	tmp := []SuratKendaraanPayload{}
	sql :=  "SELECT " +
		"	surat_kendaraan.*,kendaraan.kendaraan_id, kendaraan.hash_id 'kendaraan_id' " +
		"		FROM surat_kendaraan " +
		"	JOIN kendaraan ON surat_kendaraan.kendaraan_id=kendaraan.kendaraan_id WHERE surat_kendaraan.status = 'aktif' AND surat_kendaraan.kendaraan_id = ? "
	result := db.Raw(sql, id).Scan(&tmp)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	return tmp,nil
}


func (data *SuratKendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("surat_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("surat_id = ?",id).Delete(&data)
	if response.Error != nil {
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *SuratKendaraanResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	result := []SuratKendaraanResponse{}
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
		sql += " WHERE surat_id > ?"
		param1 = int(id)
		//trans = trans.Where("surat_id > ?",id).Find(&result)
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


func (data *SuratKendaraanPayload) defineValue()  (tmp SuratKendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update

	kadaluarsa,err := time.Parse("2006-01-02",  data.Kadaluarsa)
	if err != nil {
		return tmp,errors.New("invalid payload")
	}
	pembuatan,err := time.Parse("2006-01-02", data.TanggalPembuatan)
	if err != nil {
		return tmp,errors.New("invalid payload")
	}
	tmp.KendaraanID,err = helper.DecodeHash(data.KendaraanID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	days := kadaluarsa.Sub(pembuatan).Hours() / 24
	if days < 0 {
		return tmp,errors.New("tanggal kadalurasa harus lebih besar dari tanggal pembuatan")
	}
	tmp.JenisSurat = data.JenisSurat
	tmp.NoSurat = data.NoSurat
	tmp.Kadaluarsa =  kadaluarsa.Format("2006-01-02")
	tmp.TanggalPembuatan =  pembuatan.Format("2006-01-02")
	tmp.GenerateTime(false)
	return tmp,nil
}

func (data *SuratKendaraanResponse) switchValue(tmp *SuratKendaraanResponse) {
	data.JenisSurat = tmp.JenisSurat
	data.NoSurat = tmp.NoSurat
	data.Kadaluarsa =  tmp.Kadaluarsa
	data.TanggalPembuatan = tmp.TanggalPembuatan
	tmp.GenerateTime(true)
}

func (data *SuratKendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	//data.JenisSurat = r.FormValue("jenis_surat")
	//data.NoSurat = r.FormValue("no_surat")
	//data.Kadaluarsa =  r.FormValue("kadaluarsa")
	//data.TanggalPembuatan =  r.FormValue("tanggal_pembuatan")
	//data.KendaraanID =  r.FormValue("kendaraan_id")
	//data.Img =  "good"
	//data.Status = "1"
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *SuratKendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *SuratKendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&SuratKendaraanResponse{}).Where("surat_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *SuratKendaraanPayload) updateHashId(db *gorm.DB,r *http.Request, id int,param ...string)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)

	//filename := helper.CreateHashMd5(helper.GenerateTime().Format(time.RFC3339)+strconv.Itoa(int(data.SuratKendaraanID)))
	//checkImg := r.FormValue("image")
	//if checkImg != ""{
	//	filename,err = helper.ImgProcess(r,"img",filename)
	//}else{
	//	filename = ""
	//}
	//if err != nil{
	//	return errors.New("gagal upload file")
	//}


	var count int64
	db.Model(&SuratKendaraanResponse{}).Where("kendaraan_id",param[0]).Where("jenis_surat",param[1]).Count(&count)
	if count > 0 {
		updateStatus := db.Model(&data).Select("status").Where("kendaraan_id",param[0]).Where("jenis_surat",param[1]).Updates(SuratKendaraanPayload{Status: "tergantikan"})
		if updateStatus.Error != nil{
			return updateStatus.Error
		}
		if updateStatus.RowsAffected < 1 {
			return errors.New("gagal rubah id")
		}
	}

	response := db.Model(&data).Select("hash_id","status").Where("surat_id",id).Updates(SuratKendaraanPayload{HashID: hashID,Status: "aktif"})
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}


	return nil
}
