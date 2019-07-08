package dbase

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/db"
	"fmt"
	"time"
)

// New returns a new DataBaseStore.
func New(db *db.DB) core.DataBaseStore {
	return &dataBaseStore{db}
}

type dataBaseStore struct {
	db *db.DB
}

func (s *dataBaseStore) Create(_db *core.DataBase) error {
	_db.Created = time.Now().Format("2006-01-02 15:04:05")
	_db.Updated = _db.Created
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
	_db.Updated = time.Now().Format("2006-01-02 15:04:05")
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
			"database_name": name,
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

func (s *dataBaseStore) FindID(id int64) (*core.DataBase, error) {
	out := &core.DataBase{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_id": id,
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

func (s *dataBaseStore) List(q *core.DBQuery) ([]*core.DataBase, int, error) {
	var out []*core.DataBase
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_name": "%" + q.Name + "%",
			"database_pid":  q.PID,
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

func getQueryCountSQL(q *core.DBQuery) (queryAll string) {
	queryAll = " Select Count(1) FROM dbases Where 1=1 "
	if q.Name != "" {
		queryAll += " And database_name like :database_name "
	}
	if q.PID > 0 {
		queryAll += " And database_pid = :database_pid "
	}
	return queryAll
}

func getQueryListSQL(q *core.DBQuery, driver db.Driver) (querySQL string) {
	querySQL = queryBase + " FROM dbases Where 1=1 "
	if q.Name != "" {
		querySQL += " And database_name like :database_name "
	}
	if q.PID > 0 {
		querySQL += " And database_pid = :database_pid "
	}
	if q.OrderBy != "" {
		querySQL += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		querySQL += " ORDER BY database_created DESC "
	}

	if driver == db.Sqlite {
		querySQL += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	} else {
		querySQL += fmt.Sprintf(" limit %d, %d", q.Index*q.Size, q.Size)
	}

	return querySQL
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

func (s *dataBaseStore) DeleteByPID(pid int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"database_pid": pid,
		}
		stmt, args, err := binder.BindNamed(stmtDeleteByPID, params)
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
FROM dbases
WHERE database_name = :database_name and database_pid = :database_pid 
`

const queryID = queryBase + `
FROM dbases
WHERE database_id = :database_id 
`

const stmtInsert = `
INSERT INTO dbases (
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
UPDATE dbases
SET
database_name         	= :database_name
,database_description   = :database_description
,database_updated       = :database_updated
WHERE database_id = :database_id
`

const stmtDelete = `
DELETE FROM dbases WHERE database_id = :database_id
`

const stmtDeleteByPID = `
DELETE FROM dbases WHERE database_pid = :database_pid
`
