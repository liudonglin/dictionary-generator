package templete

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/db"
	"fmt"
	"strings"
	"time"
)

// New returns a new TempleteStore.
func New(db *db.DB) core.TempleteStore {
	return &templeteStore{db}
}

type templeteStore struct {
	db *db.DB
}

func (s *templeteStore) Create(tpl *core.Templete) error {
	tpl.Created = time.Now().Format("2006-01-02 15:04:05")
	tpl.Updated = tpl.Created
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(tpl)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		tpl.ID, err = res.LastInsertId()
		return err
	})
}

func (s *templeteStore) Update(tpl *core.Templete) error {
	tpl.Updated = time.Now().Format("2006-01-02 15:04:05")
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(tpl)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *templeteStore) FindName(name string) (*core.Templete, error) {
	out := &core.Templete{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"templete_name": name,
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

func (s *templeteStore) FindID(id int64) (*core.Templete, error) {
	out := &core.Templete{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"templete_id": id,
		}
		query, args, err := binder.BindNamed(queryID, params)
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

func (s *templeteStore) List(q *core.TempleteQuery) ([]*core.Templete, int, error) {
	var out []*core.Templete
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"templete_name":      "%" + q.Name + "%",
			"templete_language":  q.Language,
			"templete_data_base": q.DataBase,
			"templete_orm":       q.Orm,
			"templete_type":      q.Type,
			"ptr_project_id":     q.PID,
		}
		queryAll := getQueryListSQL(q, s.db.Driver())
		query, args, err := binder.BindNamed(queryAll, params)
		if err != nil {
			return err
		}
		rows, err := queryer.Query(query, args...)
		if err != nil {
			return err
		}
		out, err = scanRows(rows)

		//查询count
		queryCount := getQueryCountSQL(q)
		query, args, err = binder.BindNamed(queryCount, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		scanSingle(row, &total)

		return err
	})
	return out, total, err
}

func getQueryCountSQL(q *core.TempleteQuery) (querySQL string) {
	querySQL = " Select Count(1) FROM templetes #Join# Where 1=1 "
	// 有pid才做联合查询
	if q.PID > 0 {
		querySQL = strings.Replace(querySQL,"#Join#","LEFT JOIN project_templete_relations ON templete_id = ptr_templete_id",-1)
		querySQL += " And ptr_project_id = :ptr_project_id "
	} else {
		querySQL = strings.Replace(querySQL,"#Join#","",-1)
	}

	if q.Name != "" {
		querySQL += " And templete_name like :templete_name "
	}
	if q.Language != "" {
		querySQL += " And templete_language = :templete_language "
	}
	if q.DataBase != "" {
		querySQL += " And templete_data_base = :templete_data_base "
	}
	if q.Orm != "" {
		querySQL += " And templete_orm = :templete_orm "
	}
	if q.Type != "" {
		querySQL += " And templete_type = :templete_type "
	}
	return querySQL
}

func getQueryListSQL(q *core.TempleteQuery, driver db.Driver) (querySQL string) {
	querySQL = queryBase + " FROM templetes #Join# Where 1=1 "

	if q.PID > 0 {
		querySQL = strings.Replace(querySQL,"#Join#","LEFT JOIN project_templete_relations ON templete_id = ptr_templete_id",-1)
		querySQL += " And ptr_project_id = :ptr_project_id "
	} else {
		querySQL = strings.Replace(querySQL,"#Join#","",-1)
	}

	if q.Name != "" {
		querySQL += " And templete_name like :templete_name "
	}
	if q.Language != "" {
		querySQL += " And templete_language = :templete_language "
	}
	if q.DataBase != "" {
		querySQL += " And templete_data_base = :templete_data_base "
	}
	if q.Orm != "" {
		querySQL += " And templete_orm = :templete_orm "
	}
	if q.Type != "" {
		querySQL += " And templete_type = :templete_type "
	}
	if q.OrderBy != "" {
		querySQL += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		querySQL += " ORDER BY templete_created DESC "
	}

	if driver == db.Sqlite {
		querySQL += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	} else {
		querySQL += fmt.Sprintf(" limit %d, %d", q.Index*q.Size, q.Size)
	}
	return querySQL
}

func (s *templeteStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"templete_id": id,
		}
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *templeteStore) CreateProjectTempleteRelation(pid, tid int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"ptr_project_id":  pid,
			"ptr_templete_id": tid,
		}
		stmt, args, err := binder.BindNamed(stmtInsertProjectTempleteRelations, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		return err
	})
}

func (s *templeteStore) DeleteProjectTempleteRelationByPID(pid int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"ptr_project_id": pid,
		}
		stmt, args, err := binder.BindNamed(stmtDeleteProjectTempleteRelationByPID, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

const queryBase = `
SELECT
templete_id
,templete_name
,templete_content
,templete_language
,templete_data_base
,templete_orm
,templete_type
,templete_created
,templete_updated
`

const queryName = queryBase + `
FROM templetes
WHERE templete_name = :templete_name
`

const queryID = queryBase + `
FROM templetes
WHERE templete_id = :templete_id
`

const stmtInsert = `
INSERT INTO templetes (
templete_name
,templete_content
,templete_language
,templete_data_base
,templete_orm
,templete_type
,templete_created
,templete_updated
) VALUES (
:templete_name
,:templete_content
,:templete_language
,:templete_data_base
,:templete_orm
,:templete_type
,:templete_created
,:templete_updated
)
`

const stmtUpdate = `
UPDATE templetes
SET
templete_name         	= :templete_name
,templete_content		= :templete_content
,templete_language		= :templete_language
,templete_data_base		= :templete_data_base
,templete_orm			= :templete_orm
,templete_type			= :templete_type
,templete_updated       = :templete_updated
WHERE templete_id = :templete_id
`

const stmtDelete = `
DELETE FROM templetes WHERE templete_id = :templete_id
`

const stmtInsertProjectTempleteRelations = `
INSERT INTO project_templete_relations (
ptr_project_id
,ptr_templete_id
) VALUES (
:ptr_project_id
,:ptr_templete_id
)
`

const stmtDeleteProjectTempleteRelationByPID = `
DELETE FROM project_templete_relations WHERE ptr_project_id = :ptr_project_id
`
