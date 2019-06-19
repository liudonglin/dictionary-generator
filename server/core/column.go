package core

// Column ...
type Column struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	TID         int64  `json:"tid"`
	Title       string `json:"title" validate:"required,max=20"`
	DataType    string `json:"data_type"`
	PK          bool   `json:"pk"`
	AI          bool   `json:"ai"`
	Null        bool   `json:"null"`
	Length      string `json:"length"`
	Index       bool   `json:"index"`
	Enum        string `json:"enum"`
	Description string `json:"description" validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// ColumnQuery 分页查询参数
type ColumnQuery struct {
	Pager
	TID  int64  `json:"tid"`
	Name string `json:"name"`
}

// ColumnStore ...
type ColumnStore interface {
	Create(*Column) error

	Update(*Column) error

	FindNameAndTID(int64, string) (*Column, error)

	List(*ColumnQuery) ([]*Column, int, error)

	Delete(int64) error
}
