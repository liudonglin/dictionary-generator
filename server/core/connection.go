package core

// Connection ...
type Connection struct {
	DataBase string `json:"data_base" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     string `json:"port" validate:"required"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ConnectionQuery 分页查询参数
type ConnectionQuery struct {
	Pager
	PID  int64  `json:"pid"`
	Name string `json:"name"`
}

// ConnectionStore ...
type ConnectionStore interface {
	Create(*Connection) error

	Update(*Connection) error

	FindNameAndPID(int64, string) (*Column, error)

	List(*ColumnQuery) ([]*Column, int, error)

	Delete(int64) error

	DeleteByPID(int64) error
}
