package core

// Project 项目
type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	DataBase    string `json:"data_base"`
	Orm         string `json:"orm"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// ProjectStore 项目相关操作接口
type ProjectStore interface {

	// Create persists a new user to the datastore.
	Create(*Project) error

	Update(*Project) error

	FindName(string) (*Project, error)

	List(name string) ([]*Project, error)

	Delete(int64) error
}
