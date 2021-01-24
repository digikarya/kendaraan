package app

import (
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
type Kepegawaian struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *Kepegawaian) Initialize(config *config.Config,route *mux.Router) {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{
		PrepareStmt: true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	sqlDB, err := db.DB()
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
func (a *Kepegawaian) setRouters() {
	//// Routing for handling the projectsUserFind

	a.Post("/agen", a.guard(handler.AgenCreate))
	a.Get("/agen/all", a.guard(handler.AgenAll))
	a.Get("/agen/{hashid}", a.guard(handler.AgenFind))
	a.Put("/agen/{hashid}", a.guard(handler.AgenUpdate))
	a.Delete("/agen/{hashid}", a.guard(handler.AgenDelete))

}


// Get wraps the router for GET method
func (a *Kepegawaian) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Put wraps the router for PUT method
func (a *Kepegawaian) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Post wraps the router for POST method
func (a *Kepegawaian) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Delete wraps the router for DELETE method
func (a *Kepegawaian) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *Kepegawaian) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *Kepegawaian) guest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

func (a *Kepegawaian) guard(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//if err := authHelper.Authorization(a.DB,r,"*");err != nil {
		//	helper.RespondJSONError(w,http.StatusUnauthorized,err)
		//	return
		//}
		handler(a.DB, w, r)
	}
}
//
//func (a *Kepegawaian) guardAdmin(handler RequestHandlerFunction) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if err := authHelper.Authorization(a.DB,r,"admin");err != nil {
//			helper.RespondJSONError(w,http.StatusUnauthorized,err)
//			return
//		}
//		handler(a.DB, w, r)
//	}
//}
//
//func (a *Kepegawaian) guardClient(handler RequestHandlerFunction) http.HandlerFunc {
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
//func (a *Kepegawaian) guardSaksi(handler RequestHandlerFunction) http.HandlerFunc {
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


