package helpers

import "database/sql"

func ClearTableData(db *sql.DB, tables ...string) error {
	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}
	return nil
}
