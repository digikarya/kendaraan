package helper


type EndpointService struct {
	Kepegawaian struct{
		URL string
		Agen string
	}
	Kendaraan string
	SPJ string
}
func GetEndpoint() *EndpointService {
	return &EndpointService{
		Kepegawaian: struct {
			URL  string
			Agen string
		}{URL: "http://pegawai-go:8080", Agen: "/agen"},
		Kendaraan: "http://kendaraan-go:8080",
	}
}

