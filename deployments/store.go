package deployments

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StoreDeployment stores a deployment in the database and retrieves it
func StoreDeployment(db *sqlx.DB, deployment *Deployment) error {
	query := `
		INSERT INTO deployments (id, name, status, target_replica_count, replica_count, github_repo, github_branch)
		VALUES (:id, :name, :status, :target_replica_count, :replica_count, :github_repo, :github_branch)
		RETURNING *`

	rows, err := db.NamedQuery(query, deployment)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(deployment)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDeployment retrieves a deployment from the database by ID
func GetDeployment(db *sqlx.DB, id uuid.UUID) (*Deployment, error) {
	var deployment Deployment
	err := db.Get(&deployment, `
		SELECT * FROM deployments WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &deployment, nil
}

func UpdateDeployment(db *sqlx.DB, deployment *Deployment) error {
	result, err := db.Exec(`
		UPDATE deployments SET
			status = $1,
			target_replica_count = $2,
			replica_count = $3,
			github_repo = $4,
			github_branch = $5
		WHERE id = $6
	`, deployment.Status, deployment.TargetReplicaCount, deployment.ReplicaCount, deployment.GithubRepo, deployment.GithubBranch, deployment.ID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
