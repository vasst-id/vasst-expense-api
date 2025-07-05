package entities

// Medium represents a communication medium in the system
type Medium struct {
	MediumID    int    `json:"medium_id" db:"medium_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// CreateMediumInput is used for creating a new medium
type CreateMediumInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// UpdateMediumInput is used for updating an existing medium
type UpdateMediumInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
