package core

// Table ...
type Table struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required,max=40"`
	PID         int64     `json:"pid" validate:"gt=0"`
	DID         int64     `json:"did" validate:"gt=0"`
	Description string    `json:"description" validate:"max=200"`
	Created     string    `json:"created"`
	Updated     string    `json:"updated"`
	Columns     []*Column `json:"columns"`
}

// TableQuery 分页查询参数
type TableQuery struct {
	Pager
	PID  int64  `json:"pid"`
	DID  int64  `json:"did"`
	Name string `json:"name"`
}

// TableStore ...
type TableStore interface {
	Create(*Table) error

	Update(*Table) error

	FindNameAndDID(int64, string) (*Table, error)

	FindID(int64) (*Table, error)

	List(*TableQuery) ([]*Table, int, error)

	Delete(int64) error

	DeleteByDID(int64) error

	DeleteByPID(int64) error
}
