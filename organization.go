package main

import "time"

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	OwnerUserID string    `json:"ownerUserId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
