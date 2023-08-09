package models

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type Users struct {
	Users UsersConfig
}

type Product struct {
	Product ProductConfig
}

type UsersConfig struct {
	User_id int
	Name    string
	Email   string
}
type ProductConfig struct {
	Product_id  int
	Name        string
	Description string
	Amount      string
	Stok        int
}
