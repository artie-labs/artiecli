package internal

import (
	"time"

	"github.com/google/uuid"
)

type ListDeploymentsResponse struct {
	Items []Deployment `json:"items"`
}

type Deployment struct {
	UUID                 uuid.UUID `json:"uuid"`
	DataPlaneName        string    `json:"dataPlaneName"`
	Name                 string    `json:"name"`
	LastUpdatedAt        time.Time `json:"lastUpdatedAt"`
	Status               string    `json:"status"`
	HasUndeployedChanges bool      `json:"hasUndeployedChanges"`
}

type GetDeploymentResponse struct {
	FullDeployment FullDeployment `json:"deployment"`
}

// [FullDeployment] comes from get deployment by uuid
type FullDeployment struct {
	Deployment
	Source Source `json:"source"`
}

type Source struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	UUID          uuid.UUID `json:"uuid"`
	Schema        string    `json:"schema"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	IsBackfilling bool      `json:"isBackfilling"`
}
