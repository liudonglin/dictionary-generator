package table

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
)

// New returns a new UserStore.
func New(db *db.DB) core.TableStore {
	return &tableStore{db}
}

type tableStore struct {
	db *db.DB
}

func (s *tableStore) Create(table *core.Table) error {
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
			"table_name": "%" + name + "%",
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

func (s *tableStore) List(name string) ([]*core.Table, error) {
	panic("not implement")
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

const queryBase = `
SELECT
table_id
,table_name
,table_did
,table_description
,table_created
,table_updated
`

const queryNameAndDID = queryBase + `
FROM tables
WHERE table_name = :table_name and table_did = :table_did 
`

const stmtInsert = `
INSERT INTO tables (
 table_name
,table_did
,table_description
,table_created
,table_updated
) VALUES (
 :table_name
,:table_did
,:table_description
,:table_created
,:table_updated
)
`

const stmtUpdate = `
UPDATE tables
SET
table_name         	= :table_name
,table_description   = :table_description
,table_updated       = :table_updated
WHERE table_id = :table_id
`

const stmtDelete = `
DELETE FROM tables WHERE table_id = :table_id
`
