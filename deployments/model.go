package deployments

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	DeploymentStatusPending   = "pending"
	DeploymentStatusSucceeded = "succeeded"
	DeploymentStatusFailed    = "failed"
)

type Deployment struct {
	ID                 uuid.UUID   `db:"id" json:"id"`
	Name               string      `db:"name" json:"name"`
	Status             string      `db:"status" json:"status"`
	TargetReplicaCount int         `db:"target_replica_count" json:"targetReplicaCount"`
	ReplicaCount       int         `db:"replica_count" json:"replicaCount"`
	GithubRepo         string      `db:"github_repo" json:"githubRepo"`
	GithubBranch       string      `db:"github_branch" json:"githubBranch"`
	CreatedAt          time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time   `db:"updated_at" json:"updatedAt"`
	DeletedAt          pq.NullTime `db:"deleted_at" json:"deletedAt"`
}

type DeploymentReplica struct {
	ID           uuid.UUID `db:"id" json:"id"`
	DeploymentID uuid.UUID `db:"deployment_id" json:"deploymentId"`
	Status       string    `db:"status" json:"status"`
}

func NewDeployment(name string, githubRepo string, githubBranch string, targetReplicaCount int) *Deployment {
	return &Deployment{
		ID:                 uuid.New(),
		Name:               name,
		Status:             DeploymentStatusPending,
		TargetReplicaCount: targetReplicaCount,
		ReplicaCount:       0,
		GithubRepo:         githubRepo,
		GithubBranch:       githubBranch,
	}
}
