package column

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
	"fmt"
	"time"
)

// New returns a new ColumnStore.
func New(db *db.DB) core.ColumnStore {
	return &columnStore{db}
}

type columnStore struct {
	db *db.DB
}

func (s *columnStore) Create(column *core.Column) error {
	column.Created = time.Now().Format("2006-01-02 15:04:05")
	column.Updated = column.Created
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(column)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		column.ID, err = res.LastInsertId()
		return err
	})
}

func (s *columnStore) Update(column *core.Column) error {
	column.Updated = time.Now().Format("2006-01-02 15:04:05")
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(column)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		return err
	})
}

func (s *columnStore) FindNameAndTID(tid int64, name string) (*core.Column, error) {
	out := &core.Column{}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"column_name": name,
			"column_tid":  tid,
		}
		query, args, err := binder.BindNamed(queryNameAndTID, params)
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

func (s *columnStore) List(q *core.ColumnQuery) ([]*core.Column, int, error) {
	var out []*core.Column
	var total int
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := map[string]interface{}{
			"column_name": "%" + q.Name + "%",
			"column_tid":  q.TID,
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

func getQueryCountSqlite(q *core.ColumnQuery) (querySQL string) {
	querySQL = " Select Count(1) FROM columns Where 1=1 "
	if q.Name != "" {
		querySQL += " And column_name like :column_name "
	}
	if q.TID > 0 {
		querySQL += " And column_tid = :column_tid "
	}
	return querySQL
}

func getQueryListSqlite(q *core.ColumnQuery) (querySQL string) {
	querySQL = queryBase + " FROM columns Where 1=1 "
	if q.Name != "" {
		querySQL += " And column_name like :column_name "
	}
	if q.TID > 0 {
		querySQL += " And column_tid = :column_tid "
	}
	if q.OrderBy != "" {
		querySQL += fmt.Sprintf(" ORDER BY %s ", q.OrderBy)
	} else {
		querySQL += " ORDER BY column_created DESC "
	}

	querySQL += fmt.Sprintf(" limit %d offset %d", q.Size, q.Index*q.Size)
	return querySQL
}

func (s *columnStore) Delete(id int64) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := map[string]interface{}{
			"column_id": id,
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
column_id
,column_name
,column_tid
,column_title
,column_data_type
,column_pk
,column_null
,column_length
,column_index
,column_enum
,column_description
,column_created
,column_updated
`

const queryNameAndTID = queryBase + `
FROM columns
WHERE column_name = :column_name and column_tid = :column_tid 
`

const stmtInsert = `
INSERT INTO columns (
 column_name
,column_tid
,column_title
,column_data_type
,column_pk
,column_null
,column_length
,column_index
,column_enum
,column_description
,column_created
,column_updated
) VALUES (
 :column_name
,:column_tid
,:column_title
,:column_data_type
,:column_pk
,:column_null
,:column_length
,:column_index
,:column_enum
,:column_description
,:column_created
,:column_updated
)
`

const stmtUpdate = `
UPDATE columns
SET
column_name         = :column_name
,column_title   	= :column_title
,column_data_type   = :column_data_type
,column_pk   		= :column_pk
,column_null   		= :column_null
,column_length   	= :column_length
,column_index   	= :column_index
,column_enum   		= :column_enum
,column_description = :column_description
,column_updated     = :column_updated
WHERE column_id = :column_id
`

const stmtDelete = `
DELETE FROM columns WHERE column_id = :column_id
`
