package model

import (
	"errors"
	"github.com/digikarya/helper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type KendaraanPayload struct{
	KendaraanID    		uint `gorm:"column:kendaraan_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	JenisKendaraan 		string `json:"jenis_kendaraan"  validate:""`
	NoKendaraan 		string `json:"no_kendaraan"  validate:"required"`
	NoMesin 			string `json:"no_mesin"  validate:"required"`
	NoRangka 			string `json:"no_rangka"  validate:"required"`
	Pemilik 			string `json:"pemilik"  validate:"required"`
	MaxSeat 			uint `json:"max_seat"  validate:"required"`
	DayaAngkut 			uint `json:"daya_angkut"  validate:"required"`
	Merk 				string `json:"merk"  validate:"required"`
	TahunPembuatan 		string `json:"tahun_pembuatan"  validate:"required"`
	KapasitasMesin 		string `json:"kapasitas_mesin"  validate:"required"`
	KodeUnit 			string `json:"kode_unit"  validate:"required"`
	NoBody 				string `json:"no_body"  validate:"required"`
	TrayekID 			string `json:"trayek_id"  validate:"required"`
	KategoriKendaraanID string `json:"kategori_kendaraan_id"  validate:"required"`
	Status 				string `json:"status"  validate:""`
	Kategori 			string `json:"kategori"  validate:""`
	Trayek	 			string `json:"trayek"  validate:""`
	LayoutID   			string  `json:"layout_id"  validate:""`
	JumlahSeat			uint `json:"jumlah_seat"  validate:""`
}
type KendaraanResponse struct{
	KendaraanID    		uint `gorm:"column:kendaraan_id; PRIMARY_KEY" json:"-"`
	HashID 				string `json:"id"  validate:""`
	JenisKendaraan 		string `json:"jenis_kendaraan"  validate:""`
	NoKendaraan 		string `json:"no_kendaraan"  validate:"required"`
	NoMesin 			string `json:"no_mesin"  validate:"required"`
	NoRangka 			string `json:"no_rangka"  validate:"required"`
	Pemilik 			string `json:"pemilik"  validate:"required"`
	MaxSeat 			uint `json:"max_seat"  validate:"required"`
	DayaAngkut 			uint `json:"daya_angkut"  validate:"required"`
	Merk 				string `json:"merk"  validate:"required"`
	TahunPembuatan 		string `json:"tahun_pembuatan"  validate:"required"`
	KapasitasMesin 		string `json:"kapasitas_mesin"  validate:"required"`
	KodeUnit 			string `json:"kode_unit"  validate:"required"`
	NoBody 				string `json:"no_body"  validate:"required"`
	TrayekID 			uint `json:"trayek_id"  validate:"required"`
	KategoriKendaraanID uint `json:"kategori_kendaraan_id"  validate:"required"`
	Status 				string `json:"status"  validate:""`
	Kategori 			string `json:"kategori"  validate:""`
	Trayek	 			string `json:"trayek"  validate:""`
	LayoutID   			string  `json:"layout_id"  validate:""`
	JumlahSeat			uint `json:"jumlah_seat"  validate:""`
}

func (KendaraanPayload) TableName() string {
	return "kendaraan"
}
func (KendaraanResponse) TableName() string {
	return "kendaraan"
}


func (data *KendaraanPayload) Create(db *gorm.DB,r *http.Request) (interface{},error){
	err := data.setPayload(r)
	if err != nil {
		return nil, err
	}
	tmp,err := data.defineValue()
	if err != nil {
		return nil, err
	}
	trx := db.Begin()
	result := trx.Select("jenis_kendaraan", "no_kendaraan", "no_mesin", "no_rangka", "pemilik", "max_seat", "daya_angkut", "merk", "tahun_pembuatan", "kapasitas_mesin", "kode_unit", "no_body", "trayek_id", "kategori_kendaraan_id").Create(&tmp)
	if result.Error != nil {
		trx.Rollback()
		return nil,result.Error
	}
	if result.RowsAffected < 1 {
		trx.Rollback()
		return nil,errors.New("failed to create data")
	}
	//log.Print(tmp.AgenID)
	if err := data.updateHashId(trx,int(tmp.KendaraanID));err != nil{
		trx.Rollback()
		return nil, err
	}
	trx.Commit()
	return data,nil
}

func (data *KendaraanPayload) Update(db *gorm.DB,r *http.Request,string ...string)  (interface{},error) {
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
	tmpUpdate := KendaraanResponse{}
	if err := db.Where("kendaraan_id = ?", id).First(&tmpUpdate).Error; err != nil {
		return nil,err
	}
	tmpUpdate.switchValue(&tmp)
	result := db.Select("jenis_kendaraan", "no_kendaraan", "no_mesin", "no_rangka", "pemilik", "max_seat", "daya_angkut", "merk", "tahun_pembuatan", "kapasitas_mesin", "kode_unit", "no_body", "trayek_id", "kategori_kendaraan_id").Updates(&tmpUpdate)
	if result.Error != nil {
		return nil,errors.New("gagal update")
	}
	return data,nil
}


func (data *KendaraanResponse) Find(db *gorm.DB,string ...string) (interface{},error){
	tmp := KendaraanPayload{}
	id,err := helper.DecodeHash(string[0])
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	sql := `SELECT
			kendaraan.*,trayek.trayek_id ,
				trayek.hash_id 'trayek_id',CONCAT(trayek.asal,' - ',trayek.tujuan) 'trayek',
				kategori_kendaraan.katid 'katid',kategori_kendaraan.kategori_id 'kategori_kendaraan_id', CONCAT(kategori_kendaraan.nama,' - ',kategori_kendaraan.kode) 'kategori',
                kategori_kendaraan.total_seat 'jumlah_seat', kategori_kendaraan.layout_id 'layout_id'
			FROM kendaraan
			JOIN trayek ON kendaraan.trayek_id=trayek.trayek_id
			JOIN (SELECT k.kategori_id 'katid',k.nama,k.kode,k.hash_id 'kategori_id',l.layout_id 'layid',  l.hash_id 'layout_id',l.total_seat FROM kategori_kendaraan k JOIN layout_kursi l ON k.layout_kursi_id=l.layout_id) kategori_kendaraan ON kendaraan.kategori_kendaraan_id=kategori_kendaraan.katid
			WHERE kendaraan_id = ?`
	result := db.Raw(sql+" LIMIT 1", id).Scan(&tmp)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	getSurat := SuratKendaraanResponse{}
	activeSurat,err := getSurat.FindByKendaraanActive(db,tmp.HashID)
	//if err !=  nil{
	//	activeSurat = []SuratKendaraanResponse{}
	//}
	var count int64
	count = 0
	db.Model(&SuratKendaraanResponse{}).Where("status = ?", "aktif").Where("NOW() > kadaluarsa").Count(&count)
	return struct {
		JumlahSuratKadaluarsa int64 `json:"jumlah_surat_kadaluarsa"`
		KendaraanPayload
		Surat interface{} `json:"surat"`
	}{KendaraanPayload: tmp, JumlahSuratKadaluarsa: count,Surat: activeSurat},nil
}



func (data *KendaraanPayload) Delete(db *gorm.DB,string ...string) (interface{},error){
	hashid := string[0]
	id,err := helper.DecodeHash(hashid)
	if err != nil {
		return nil,errors.New("data tidak sesuai")
	}
	result := db.Where("kendaraan_id = ?", id).Find(&data)
	if err := result.Error; err != nil {
		return nil,err
	}
	if result.RowsAffected < 1 {
		return nil,errors.New("data tidak ditemukan")
	}
	response := db.Where("kendaraan_id = ?",id).Delete(&data)
	if response.Error != nil {
		log.Print(response.Error)
		return nil,errors.New("gagal di hapus")
	}
	db.Commit()
	return data,nil
}


func (data *KendaraanResponse) All(db *gorm.DB,string ...string) (interface{}, error) {
	result := []KendaraanPayload{}
	limit,err := strconv.Atoi(string[1])
	if err != nil {
		return nil, err
	}
	//trans := db.Limit(limit).Find(&result)
	sql := `SELECT
			kendaraan.*,trayek.trayek_id ,
				trayek.hash_id 'trayek_id',CONCAT(trayek.asal,' - ',trayek.tujuan) 'trayek',
				kategori_kendaraan.kategori_id 'katid',kategori_kendaraan.hash_id 'kategori_kendaraan_id', CONCAT(kategori_kendaraan.nama,' - ',kategori_kendaraan.kode) 'kategori'
			FROM kendaraan
			JOIN trayek ON kendaraan.trayek_id=trayek.trayek_id
			JOIN kategori_kendaraan ON kendaraan.kategori_kendaraan_id=kategori_kendaraan.kategori_id`
	hashID := string[0]
	param1 := limit
	param2 := limit
	if hashID != "" {
		id,err := helper.DecodeHash(hashID)
		if err != nil {
			return nil,err
		}
		sql += " WHERE kendaraan_id > ?"
		param1 = int(id)
		//trans = trans.Where("kendaraan_id > ?",id).Find(&result)
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


func (data *KendaraanPayload) defineValue()  (tmp KendaraanResponse,err error) {
	// ambil data dari payload menjadi data siap insert atau update
	tmp.KendaraanID = data.KendaraanID
	tmp.JenisKendaraan = data.JenisKendaraan
	tmp.NoKendaraan = data.NoKendaraan
	tmp.NoMesin = data.NoMesin
	tmp.NoRangka = data.NoRangka
	tmp.Pemilik = data.Pemilik
	tmp.MaxSeat = data.MaxSeat
	tmp.DayaAngkut = data.DayaAngkut
	tmp.Merk = data.Merk
	tmp.TahunPembuatan = data.TahunPembuatan
	tmp.KapasitasMesin = data.KapasitasMesin
	tmp.KodeUnit = data.KodeUnit
	tmp.NoBody = data.NoBody
	tmp.Status = data.Status
	tmp.TrayekID,err = helper.DecodeHash(data.TrayekID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	tmp.KategoriKendaraanID,err = helper.DecodeHash(data.KategoriKendaraanID)
	if err != nil {
		return tmp,errors.New("data tidak sesuai")
	}
	return tmp,nil
}

func (data *KendaraanResponse) switchValue(tmp *KendaraanResponse) {
	data.JenisKendaraan = tmp.JenisKendaraan
	data.NoKendaraan = tmp.NoKendaraan
	data.NoMesin = tmp.NoMesin
	data.NoRangka = tmp.NoRangka
	data.Pemilik = tmp.Pemilik
	data.MaxSeat = tmp.MaxSeat
	data.DayaAngkut = tmp.DayaAngkut
	data.Merk = tmp.Merk
	data.TahunPembuatan = tmp.TahunPembuatan
	data.KapasitasMesin = tmp.KapasitasMesin
	data.KodeUnit = tmp.KodeUnit
	data.NoBody = tmp.NoBody
	data.Status = tmp.Status
	data.TrayekID = tmp.TrayekID
	data.KategoriKendaraanID = tmp.KategoriKendaraanID
}

func (data *KendaraanPayload) setPayload(r *http.Request)  (err error)  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}



func (data *KendaraanResponse) setPayload(r *http.Request)  error  {
	if err := helper.DecodeJson(r,&data);err != nil {
		return errors.New("invalid payload")
	}
	if err := helper.ValidateData(data);err != nil {
		return err
	}
	return nil
}

func (data *KendaraanPayload) countData(db *gorm.DB,id uint) (int64,error) {
	var count int64
	db.Model(&KendaraanResponse{}).Where("kendaraan_id = ?", id).Count(&count)
	if count < 1 {
		return count, errors.New("data tidak ditemukan")
	}
	return count,nil

}


func (data *KendaraanPayload) updateHashId(db *gorm.DB, id int)  error {
	hashID,err := helper.EncodeHash(id)
	if err != nil {
		return err
	}
	//log.Print(tmp.DaerahID)
	response := db.Model(&data).Where("kendaraan_id",id).Update("hash_id", hashID)
	if response.Error != nil{
		return response.Error
	}
	if response.RowsAffected < 1 {
		return errors.New("gagal rubah id")
	}
	return nil
}
