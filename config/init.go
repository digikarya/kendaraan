package config

type Config struct {
	DB *DBConfig
	DSN string
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Charset  string
	Database  string

}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     "kendaraan-mysql",
			Port:     "33061",
			Username: "root",
			Password: "admin321",
			Charset:  "utf8",
			Database:  "kendaraan",
		},
	}
}

func (c *Config)  GetDSN() string {
	//DSN Format
	//dsn := "zhi:admin123@tcp(127.0.0.1:3306)/authHelper?charset=utf8mb4&parseTime=True&loc=Local"
	return c.DB.Username+":"+c.DB.Password+"@tcp("+c.DB.Host+")/"+c.DB.Database+"?charset=utf8mb4&parseTime=True&loc=Local"
}