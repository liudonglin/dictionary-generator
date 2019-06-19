package core

// Connection ...
type Connection struct {
	DataBase string `json:"data_base" validate:"required"`
	Host string `json:"host" validate:"required"`
	Port string `json:"port" validate:"required"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}