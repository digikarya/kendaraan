package app

import (
	"github.com/digikarya/helper"
	"github.com/digikarya/kendaraan/app/handler"
	"github.com/digikarya/kendaraan/config"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// App has router and db instances
type Kendaraan struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *Kendaraan) Initialize(config *config.Config,route *mux.Router) {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{
		PrepareStmt: true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Could not connect database")
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5*time.Minute)
	a.DB = db
	a.Router = route
	a.setRouters()
	log.Println("app server is running")
}

// setRouters sets the all required routers
func (a *Kendaraan) setRouters() {
	//// Routing for handling the projectsUserFind

	a.Post("/layout", a.guard(handler.LayoutCreate))
	a.Get("/layout/all", a.guard(handler.LayoutAll))
	a.Get("/layout/{hashid}", a.guard(handler.LayoutFind))
	a.Put("/layout/{hashid}", a.guard(handler.LayoutUpdate))
	a.Delete("/layout/{hashid}", a.guard(handler.LayoutDelete))

	a.Post("/jenis_kendaraan", a.guard(handler.JenisKendaraanCreate))
	a.Get("/jenis_kendaraan/all", a.guard(handler.JenisKendaraanAll))
	a.Get("/jenis_kendaraan/{hashid}", a.guard(handler.JenisKendaraanFind))
	a.Put("/jenis_kendaraan/{hashid}", a.guard(handler.JenisKendaraanUpdate))
	a.Delete("/jenis_kendaraan/{hashid}", a.guard(handler.JenisKendaraanDelete))

	a.Post("/kategori_kendaraan", a.guard(handler.KategoriKendaraanCreate))
	a.Get("/kategori_kendaraan/all", a.guard(handler.KategoriKendaraanAll))
	a.Get("/kategori_kendaraan/{hashid}", a.guard(handler.KategoriKendaraanFind))
	a.Put("/kategori_kendaraan/{hashid}", a.guard(handler.KategoriKendaraanUpdate))
	a.Delete("/kategori_kendaraan/{hashid}", a.guard(handler.KategoriKendaraanDelete))

	a.Post("/trayek", a.guard(handler.TrayekCreate))
	a.Get("/trayek/all", a.guard(handler.TrayekAll))
	a.Get("/trayek/{hashid}", a.guard(handler.TrayekFind))
	a.Put("/trayek/{hashid}", a.guard(handler.TrayekUpdate))
	a.Delete("/trayek/{hashid}", a.guard(handler.TrayekDelete))
	a.Post("/trayek/detail/", a.guard(handler.DetailTrayekCreate))
	a.Put("/trayek/detail/{hashid}", a.guard(handler.DetailTrayekUpdate))
	a.Delete("/trayek/detail/{hashid}", a.guard(handler.DetailTrayekDelete))


	a.Post("/check_list", a.guard(handler.CheckListKendaraanCreate))
	a.Get("/check_list/all", a.guard(handler.CheckListKendaraanAll))
	a.Get("/check_list/{hashid}", a.guard(handler.CheckListKendaraanFind))
	a.Get("/check_list/byKategori/{hashid}", a.guard(handler.CheckListKendaraanFindByKategori))
	a.Put("/check_list/{hashid}", a.guard(handler.CheckListKendaraanUpdate))
	a.Delete("/check_list/{hashid}", a.guard(handler.CheckListKendaraanDelete))

	a.Post("/check_list/detail/", a.guard(handler.DetailCheckListKendaraanCreate))
	a.Put("/check_list/detail/{hashid}", a.guard(handler.DetailCheckListKendaraanUpdate))
	a.Delete("/check_list/detail/{hashid}", a.guard(handler.DetailCheckListKendaraanDelete))

	a.Post("/jadwal", a.guard(handler.JadwalCreate))
	a.Get("/jadwal/all", a.guard(handler.JadwalAll))
	a.Get("/jadwal/{hashid}", a.guard(handler.JadwalFind))
	a.Put("/jadwal/{hashid}", a.guard(handler.JadwalUpdate))
	a.Delete("/jadwal/{hashid}", a.guard(handler.JadwalDelete))

	a.Post("/kendaraan", a.guard(handler.KendaraanCreate))
	a.Get("/kendaraan/all", a.guard(handler.KendaraanAll))
	a.Get("/kendaraan/{hashid}", a.guard(handler.KendaraanFind))
	a.Put("/kendaraan/{hashid}", a.guard(handler.KendaraanUpdate))
	a.Delete("/kendaraan/{hashid}", a.guard(handler.KendaraanDelete))
	a.Post("/kendaraan/surat", a.guard(handler.SuratKendaraanCreate))
	a.Get("/kendaraan/surat/{hashid}", a.guard(handler.SuratKendaraanFindByKendaraanAll))
	a.Get("/kendaraan/surat/{hashid}/active", a.guard(handler.SuratKendaraanFindByKendaraanActive))
	a.Delete("/kendaraan/surat/{hashid}", a.guard(handler.SuratKendaraanDelete))
	a.Put("/kendaraan/surat/{hashid}", a.guard(handler.SuratKendaraanUpdate))

	a.Post("/layout/search", a.guard(handler.SearchLayout))
	a.Post("/jenis_kendaraan/search", a.guard(handler.SearchJenisKendaraan))
	a.Post("/kategori_kendaraan/search", a.guard(handler.SearchKategoriKendaraan))
	a.Post("/check_list/search", a.guard(handler.SearchCheckList))
	a.Post("/trayek/search", a.guard(handler.SearchTrayek))
	a.Post("/kendaraan/search", a.guard(handler.SearchKendaraan))
	a.Post("/jadwal/search", a.guard(handler.SearchJadwal))
}


// Get wraps the router for GET method
func (a *Kendaraan) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Put wraps the router for PUT method
func (a *Kendaraan) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Post wraps the router for POST method
func (a *Kendaraan) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Delete wraps the router for DELETE method
func (a *Kendaraan) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *Kendaraan) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *Kendaraan) guest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

func (a *Kendaraan) guard(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := helper.AuthorizeRole(a.DB,r,"admin");err != nil {
			helper.RespondJSONError(w,http.StatusUnauthorized,err)
			return
		}
		handler(a.DB, w, r)
	}
}
//
//func (a *Kendaraan) guardAdmin(handler RequestHandlerFunction) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if err := authHelper.Authorization(a.DB,r,"admin");err != nil {
//			helper.RespondJSONError(w,http.StatusUnauthorized,err)
//			return
//		}
//		handler(a.DB, w, r)
//	}
//}
//
//func (a *Kendaraan) guardClient(handler RequestHandlerFunction) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		helper.CorsHelper(w,r)
//		if r.Method == http.MethodOptions {
//			return
//		}
//		if err := authHelper.Authorization(a.DB,r,"client");err != nil {
//			helper.RespondJSONError(w,http.StatusUnauthorized,err)
//			return
//		}
//		handler(a.DB, w, r)
//	}
//}
//
//func (a *Kendaraan) guardSaksi(handler RequestHandlerFunction) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		helper.CorsHelper(w,r)
//		if r.Method == http.MethodOptions {
//			return
//		}
//		if err := authHelper.Authorization(a.DB,r,"saksi");err != nil {
//			helper.RespondJSONError(w,http.StatusUnauthorized,err)
//			return
//		}
//		handler(a.DB, w, r)
//	}
//}
//
//


