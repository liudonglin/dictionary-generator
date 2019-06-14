package core

// DataBase 数据库
type DataBase struct {
	ID          int64  `json:"id"`
	PID         int64  `json:"pid"`
	Name        string `json:"name" validate:"required,max=20"`
	Description string `json:"description"  validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// DataBaseStore ...
type DataBaseStore interface {
	Create(*DataBase) error

	Update(*DataBase) error

	FindNameAndPID(int64, string) (*DataBase, error)

	List(name string) ([]*DataBase, error)

	Delete(int64) error
}
