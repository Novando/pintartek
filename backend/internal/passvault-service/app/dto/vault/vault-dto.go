package vault

import "time"

type VaultRequest struct {
	Name       string     `json:"name" validate:"required"`
	Credential Credential `json:"credential"`
}

type VaultResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VaultEditRequest struct {
	Name string `json:"name" validate:"required"`
}
