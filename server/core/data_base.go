package core

// DataBase 数据库
type DataBase struct {
	ID          int64    `json:"id"`
	PID         int64    `json:"pid" validate:"gt=0"`
	Name        string   `json:"name" validate:"required,max=40"`
	Description string   `json:"description"  validate:"max=200"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
	Tables      []*Table `json:"tables"`
}

// DBQuery 分页查询参数
type DBQuery struct {
	Pager
	PID  int64  `json:"pid"`
	Name string `json:"name"`
}

// DataBaseStore ...
type DataBaseStore interface {
	Create(*DataBase) error

	Update(*DataBase) error

	FindNameAndPID(int64, string) (*DataBase, error)

	FindID(int64) (*DataBase, error)

	List(q *DBQuery) ([]*DataBase, int, error)

	Delete(int64) error

	DeleteByPID(int64) error
}
