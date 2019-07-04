package connection

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/db"
	"fmt"
	"time"
)

// New returns a new ColumnStore.
func New(db *db.DB) core.ConnectionStore {
	return &connectionStore{db}
}

type connectionStore struct {
	db *db.DB
}

func (s *connectionStore) Create(connection *core.Connection) error {
	connection.Created = time.Now().Format("2006-01-02 15:04:05")
	connection.Updated = connection.Created
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(connection)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		connection.ID, err = res.LastInsertId()
		return err
	})
}

func (s *connectionStore) Update(connection *core.Connection) error {
	connection.Updated = time.Now().Format("2006-01-02 15:04:05")
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(connection)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *connectionStore) FindNameAndPID(pid int64, name string) (*core.Connection, error) {
	out := &core.Connection{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"connection_name": name,
			"connection_pid":  pid,
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

func (s *connectionStore) List(q *core.ConnectionQuery) ([]*core.Connection, int, error) {
	var out []*core.Connection
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"connection_name": "%" + q.Name + "%",
			"connection_pid":  q.PID,
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

func getQueryCountSqlite(q *core.ConnectionQuery) (querySQL string) {
	querySQL = " Select Count(1) FROM connections Where 1=1 "
	if q.Name != "" {
		querySQL += " And connection_name like :connection_name "
	}
	if q.PID > 0 {
		querySQL += " And connection_pid = :connection_pid "
	}
	return querySQL
}

func getQueryListSqlite(q *core.ConnectionQuery) (querySQL string) {
	querySQL = queryBase + " FROM connections Where 1=1 "
	if q.Name != "" {
		querySQL += " And connection_name like :connection_name "
	}
	if q.PID > 0 {
		querySQL += " And connection_pid = :connection_pid "
	}
	if q.OrderBy != "" {
		querySQL += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		querySQL += " ORDER BY connection_id ASC "
	}

	querySQL += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	return querySQL
}

func (s *connectionStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"connection_id": id,
		}
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *connectionStore) DeleteByPID(pid int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"connection_pid": pid,
		}
		stmt, args, err := binder.BindNamed(stmtDeletePID, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

const queryBase = `
SELECT
connection_id
,connection_name
,connection_pid
,connection_data_base
,connection_host
,connection_port
,connection_user
,connection_password
,connection_description
,connection_created
,connection_updated
`

const queryNameAndPID = `
FROM connections
WHERE connection_name = :connection_name and connection_pid = :connection_pid 
`

const stmtInsert = `
INSERT INTO connections (
 connection_name
,connection_pid
,connection_data_base
,connection_host
,connection_port
,connection_user
,connection_password
,connection_description
,connection_created
,connection_updated
) VALUES (
 :connection_name
,:connection_pid
,:connection_data_base
,:connection_host
,:connection_port
,:connection_user
,:connection_password
,:connection_description
,:connection_created
,:connection_updated
)
`

const stmtUpdate = `
UPDATE connections
SET
connection_name         	= :connection_name
,connection_pid   			= :connection_pid
,connection_data_base   	= :connection_data_base
,connection_host   			= :connection_host
,connection_port   			= :connection_port
,connection_user   			= :connection_user
,connection_password   		= :connection_password
,connection_description   	= :connection_description
,connection_created   		= :connection_created
,connection_updated 		= :connection_updated
WHERE connection_id = :connection_id
`

const stmtDelete = `
DELETE FROM connections WHERE connection_id = :connection_id
`

const stmtDeletePID = `
DELETE FROM connections WHERE connection_pid = :connection_pid
`
