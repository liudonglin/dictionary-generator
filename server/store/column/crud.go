package column

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
)

// New returns a new UserStore.
func New(db *db.DB) core.ColumnStore {
	return &columnStore{db}
}

type columnStore struct {
	db *db.DB
}

func (s *columnStore) Create(column *core.Column) error {
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
			"column_name": "%" + name + "%",
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

func (s *columnStore) List(name string) ([]*core.Column, error) {
	panic("not implement")
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
INSERT INTO tables (
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
