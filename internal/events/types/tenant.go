package types

// TenantContext represents tenant-specific context for multi-tenancy
type TenantContext struct {
	OrganizationID string `json:"organization_id"`
	Domain         string `json:"domain"`
}
