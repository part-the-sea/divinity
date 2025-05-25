package main

import "time"

type School struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Zip            string    `json:"zip"`
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
