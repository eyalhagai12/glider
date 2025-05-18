package projects

import "github.com/jmoiron/sqlx"

func StoreProject(db *sqlx.Tx, project *Project) error {
	_, err := db.NamedExec(`INSERT INTO projects (id, name, description) VALUES (:id, :name, :description)`, project)
	if err != nil {
		return err
	}

	return nil
}
