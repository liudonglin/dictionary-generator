package table

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/db"
	"fmt"
	"time"
)

// New returns a new TableStore.
func New(db *db.DB) core.TableStore {
	return &tableStore{db}
}

type tableStore struct {
	db *db.DB
}

func (s *tableStore) Create(table *core.Table) error {
	table.Created = time.Now().Format("2006-01-02 15:04:05")
	table.Updated = table.Created
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(table)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		table.ID, err = res.LastInsertId()
		return err
	})
}

func (s *tableStore) Update(table *core.Table) error {
	table.Updated = time.Now().Format("2006-01-02 15:04:05")
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(table)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *tableStore) FindNameAndDID(did int64, name string) (*core.Table, error) {
	out := &core.Table{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_name": name,
			"table_did":  did,
		}
		query, args, err := binder.BindNamed(queryNameAndDID, params)
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

func (s *tableStore) FindID(id int64) (*core.Table, error) {
	out := &core.Table{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_id": id,
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

func (s *tableStore) List(q *core.TableQuery) ([]*core.Table, int, error) {
	var out []*core.Table
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_name": "%" + q.Name + "%",
			"table_did":  q.DID,
			"table_pid":  q.PID,
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

func getQueryCountSQL(q *core.TableQuery) (querySQL string) {
	querySQL = " Select Count(1) FROM tables Where 1=1 "
	if q.Name != "" {
		querySQL += " And table_name like :table_name "
	}
	if q.PID > 0 {
		querySQL += " And table_pid = :table_pid "
	}
	if q.DID > 0 {
		querySQL += " And table_did = :table_did "
	}
	return querySQL
}

func getQueryListSQL(q *core.TableQuery, driver db.Driver) (querySQL string) {
	querySQL = queryBase + " FROM tables Where 1=1 "
	if q.Name != "" {
		querySQL += " And table_name like :table_name "
	}
	if q.PID > 0 {
		querySQL += " And table_pid = :table_pid "
	}
	if q.DID > 0 {
		querySQL += " And table_did = :table_did "
	}
	if q.OrderBy != "" {
		querySQL += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		querySQL += " ORDER BY table_created DESC "
	}

	if driver == db.Sqlite {
		querySQL += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	} else {
		querySQL += fmt.Sprintf(" limit %d, %d", q.Index*q.Size, q.Size)
	}
	return querySQL
}

func (s *tableStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_id": id,
		}
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *tableStore) DeleteByDID(did int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_did": did,
		}
		stmt, args, err := binder.BindNamed(stmtDeleteByDID, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *tableStore) DeleteByPID(pid int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"table_pid": pid,
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
table_id
,table_name
,table_pid
,table_did
,table_title
,table_description
,table_created
,table_updated
`

const queryNameAndDID = queryBase + `
FROM tables
WHERE table_name = :table_name and table_did = :table_did 
`

const queryID = queryBase + `
FROM tables
WHERE table_id = :table_id
`

const stmtInsert = `
INSERT INTO tables (
 table_name
,table_pid
,table_did
,table_title
,table_description
,table_created
,table_updated
) VALUES (
 :table_name
,:table_pid
,:table_did
,:table_title
,:table_description
,:table_created
,:table_updated
)
`

const stmtUpdate = `
UPDATE tables
SET
table_name         		= :table_name
,table_title   			= :table_title
,table_description   	= :table_description
,table_updated       	= :table_updated
WHERE table_id = :table_id
`

const stmtDelete = `
DELETE FROM tables WHERE table_id = :table_id
`

const stmtDeleteByDID = `
DELETE FROM tables WHERE table_did = :table_did
`

const stmtDeleteByPID = `
DELETE FROM tables WHERE table_pid = :table_pid
`
