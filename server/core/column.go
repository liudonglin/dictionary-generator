package core

// Column ...
type Column struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required,max=40"`
	PID         int64  `json:"pid" validate:"gt=0"`
	DID         int64  `json:"did" validate:"gt=0"`
	TID         int64  `json:"tid" validate:"gt=0"`
	Title       string `json:"title" validate:"required,max=40"`
	DataType    string `json:"data_type"`
	ColumnType  string `json:"column_type"`
	PK          bool   `json:"pk"`
	AI          bool   `json:"ai"`
	Null        bool   `json:"null"`
	Length      string `json:"length"`
	Index       bool   `json:"index"`
	Unique      bool   `json:"unique"`
	Enum        string `json:"enum"`
	Description string `json:"description" validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// ColumnQuery 分页查询参数
type ColumnQuery struct {
	Pager
	PID  int64  `json:"pid"`
	DID  int64  `json:"did"`
	TID  int64  `json:"tid"`
	Name string `json:"name"`
}

// ColumnStore ...
type ColumnStore interface {
	Create(*Column) error

	Update(*Column) error

	FindNameAndTID(int64, string) (*Column, error)

	FindPK(int64) (*Column, error)

	List(*ColumnQuery) ([]*Column, int, error)

	Delete(int64) error

	DeleteByTID(int64) error

	DeleteByDID(int64) error

	DeleteByPID(int64) error
}
