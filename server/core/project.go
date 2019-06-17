package core

// Project 项目
type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	Description string `json:"description" validate:"max=200"`
	Language    string `json:"language" validate:"oneof=java csharp"`
	DataBase    string `json:"data_base" validate:"oneof=mysql sqlserver"`
	Orm         string `json:"orm" validate:"oneof=mybatis smartSql"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// ProjectQuery 分页查询参数
type ProjectQuery struct {
	Pager
	Name string `json:"name"`
}

// ProjectStore 项目相关操作接口
type ProjectStore interface {

	// Create persists a new user to the datastore.
	Create(*Project) error

	Update(*Project) error

	FindName(string) (*Project, error)

	FindID(int64) (*Project, error)

	List(q *ProjectQuery) ([]*Project, int, error)

	Delete(int64) error
}
