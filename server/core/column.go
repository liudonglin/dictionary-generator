package core

// Column ...
type Column struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	TID         int64  `json:"tid"`
	Title       string `json:"title" validate:"required,max=20"`
	DataType    string `json:"data_type"`
	PK          bool   `json:"pk"`
	Null        bool   `json:"null"`
	Length      string   `json:"length"`
	Index       bool   `json:"index"`
	Enum        string `json:"enum"`
	Description string `json:"description" validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// ColumnStore ...
type ColumnStore interface {
	Create(*Column) error

	Update(*Column) error

	FindNameAndTID(int64, string) (*Column, error)

	List(name string) ([]*Column, error)

	Delete(int64) error
}