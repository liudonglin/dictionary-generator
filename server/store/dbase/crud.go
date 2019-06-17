package dbase

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
	"fmt"
)

// New returns a new DataBaseStore.
func New(db *db.DB) core.DataBaseStore {
	return &dataBaseStore{db}
}

type dataBaseStore struct {
	db *db.DB
}

func (s *dataBaseStore) Create(_db *core.DataBase) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(_db)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		_db.ID, err = res.LastInsertId()
		return err
	})
}

func (s *dataBaseStore) Update(_db *core.DataBase) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(_db)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *dataBaseStore) FindNameAndPID(pid int64, name string) (*core.DataBase, error) {
	out := &core.DataBase{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_name": "%" + name + "%",
			"database_pid":  pid,
		}
		query, args, err := binder.BindNamed(queryNameAndPID, params)
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

func (s *dataBaseStore) List(q *core.DBQuery) ([]*core.DataBase, int, error) {
	var out []*core.DataBase
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_name": "%" + q.Name + "%",
			"database_pid":  q.PID,
		}
		queryAll := ""
		if s.db.Driver() == db.Sqlite {
			queryAll = getQueryListSqlite(q)
		} else {
			panic("mysql query not implement")
		}
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
		queryCount := ""
		if s.db.Driver() == db.Sqlite {
			queryCount = getQueryCountSqlite(q)
		} else {
			panic("mysql query not implement")
		}
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

func getQueryCountSqlite(q *core.DBQuery) (queryAll string) {
	queryAll = " Select Count(1) FROM database Where 1=1 "
	if q.Name != "" {
		queryAll += " And database_name like :database_name "
	}
	if q.PID > 0 {
		queryAll += " And database_pid = :database_pid "
	}
	return queryAll
}

func getQueryListSqlite(q *core.DBQuery) (queryAll string) {
	queryAll = queryBase + " FROM database Where 1=1 "
	if q.Name != "" {
		queryAll += " And database_name like :database_name "
	}
	if q.PID > 0 {
		queryAll += " And database_pid = :database_pid "
	}
	if q.OrderBy != "" {
		queryAll += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		queryAll += " ORDER BY database_created DESC "
	}

	queryAll += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	return queryAll
}

func (s *dataBaseStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_id": id,
		}
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

const queryBase = `
SELECT
database_id
,database_name
,database_pid
,database_description
,database_created
,database_updated
`

const queryNameAndPID = queryBase + `
FROM database
WHERE database_name = :database_name and database_pid = :database_pid 
`

const stmtInsert = `
INSERT INTO database (
 database_name
,database_pid
,database_description
,database_created
,database_updated
) VALUES (
 :database_name
,:database_pid
,:database_description
,:database_created
,:database_updated
)
`

const stmtUpdate = `
UPDATE database
SET
database_name         	= :database_name
,database_description   = :database_description
,database_updated       = :database_updated
WHERE database_id = :database_id
`

const stmtDelete = `
DELETE FROM database WHERE database_id = :database_id
`
