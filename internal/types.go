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
