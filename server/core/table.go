package core

// Table ...
type Table struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	DID         int64  `json:"did"`
	Description string `json:"description" validate:"max=200"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// TableStore ...
type TableStore interface {
	Create(*Table) error

	Update(*Table) error

	FindNameAndDID(int64, string) (*Table, error)

	List(name string) ([]*Table, error)

	Delete(int64) error
}
