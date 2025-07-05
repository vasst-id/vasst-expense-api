package entities

// Role represents a user role in the system
type Role struct {
	RoleID      int    `json:"role_id" db:"role_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// CreateRoleInput is used for creating a new role
type CreateRoleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// UpdateRoleInput is used for updating an existing role
type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
