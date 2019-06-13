package project

import (
	"database/sql"

	"code-server/core"
	"code-server/store/base/db"
)

// New returns a new UserStore.
func New(db *db.DB) core.ProjectStore {
	return &projectStore{db}
}

type projectStore struct {
	db *db.DB
}

func (s *projectStore) Create(project *core.Project) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(project)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		project.ID, err = res.LastInsertId()
		return err
	})
}

func (s *projectStore) Update(project *core.Project) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(project)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

// FindName 根据name返回Project
func (s *projectStore) FindName(name string) (*core.Project, error) {
	out := &core.Project{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"project_name": "%" + name + "%",
		}
		query, args, err := binder.BindNamed(queryName, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		err = scanRow(row, out)
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	})
	return out, err
}

func (s *projectStore) List(name string) ([]*core.Project, error) {
	var out []*core.Project
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"project_name": "%" + name + "%",
		}
		queryAll := getQueryAllSQL(name)
		query, args, err := binder.BindNamed(queryAll, params)
		if err != nil {
			return err
		}
		rows, err := queryer.Query(query, args...)
		if err != nil {
			return err
		}
		out, err = scanRows(rows)
		return err
	})
	return out, err
}

func (s *projectStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"project_id": id,
		}
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func getQueryAllSQL(name string) (queryAll string) {
	queryAll = queryBase + " FROM projects "
	if name != "" {
		queryAll += " Where project_name like :project_name"
	}
	queryAll += " ORDER BY project_created DESC "
	return queryAll
}

const queryBase = `
SELECT
project_id
,project_name
,project_language
,project_data_base
,project_orm
,project_description
,project_created
,project_updated
`

const queryName = queryBase + `
FROM projects
WHERE project_name = :project_name
`

const stmtInsert = `
INSERT INTO projects (
 project_name
,project_language
,project_data_base
,project_orm
,project_description
,project_created
,project_updated
) VALUES (
 :project_name
,:project_language
,:project_data_base
,:project_orm
,:project_description
,:project_created
,:project_updated
)
`

const stmtUpdate = `
UPDATE projects
SET
project_name         	= :project_name
,project_language       = :project_language
,project_data_base      = :project_data_base
,project_orm        	= :project_orm
,project_description    = :project_description
,project_updated        = :project_updated
WHERE project_id = :project_id
`

const stmtDelete = `
DELETE FROM projects WHERE project_id = :project_id
`