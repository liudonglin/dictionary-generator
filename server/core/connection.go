package core

// Connection ...
type Connection struct {
	ID          int64  `json:"id"`
	PID         int64  `json:"pid"`
	Name        int64  `json:"name" validate:"max=20"`
	DataBase    string `json:"data_base" validate:"required,max=20"`
	Host        string `json:"host" validate:"required,max=20"`
	Port        string `json:"port" validate:"required,max=20"`
	User        string `json:"user" validate:"required,max=20"`
	Password    string `json:"password" validate:"required,max=20"`
	Description string `json:"description" validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
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

	FindNameAndPID(int64, string) (*Connection, error)

	List(*ConnectionQuery) ([]*Connection, int, error)

	Delete(int64) error

	DeleteByPID(int64) error
}
