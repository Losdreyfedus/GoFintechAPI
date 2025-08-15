package auth

import (
	"fmt"
	"sync"
)

// Permission represents a specific action on a resource
type Permission struct {
	Resource   string `json:"resource"`              // e.g., "user", "transaction", "balance"
	Action     string `json:"action"`                // e.g., "read", "write", "delete", "admin"
	ResourceID string `json:"resource_id,omitempty"` // e.g., "self", "all", specific ID
}

// Role represents a user role with permissions
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	Inherits    []string     `json:"inherits,omitempty"` // Inherit from other roles
}

// Policy represents an access control policy
type Policy struct {
	Resource   string                 `json:"resource"`
	Actions    []string               `json:"actions"`
	Roles      []string               `json:"roles"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
}

// RBACManager manages role-based access control
type RBACManager struct {
	roles    map[string]*Role
	policies map[string]*Policy
	mu       sync.RWMutex
}

// NewRBACManager creates a new RBAC manager
func NewRBACManager() *RBACManager {
	rbac := &RBACManager{
		roles:    make(map[string]*Role),
		policies: make(map[string]*Policy),
	}

	// Initialize default roles
	rbac.initializeDefaultRoles()
	rbac.initializeDefaultPolicies()

	return rbac
}

// initializeDefaultRoles creates default system roles
func (rbac *RBACManager) initializeDefaultRoles() {
	// Super Admin - Full access
	rbac.AddRole(&Role{
		Name: "super_admin",
		Permissions: []Permission{
			{Resource: "*", Action: "*"}, // All resources, all actions
		},
	})

	// Admin - Administrative access
	rbac.AddRole(&Role{
		Name: "admin",
		Permissions: []Permission{
			{Resource: "user", Action: "*"},
			{Resource: "transaction", Action: "*"},
			{Resource: "balance", Action: "*"},
			{Resource: "audit", Action: "read"},
		},
	})

	// Manager - Limited administrative access
	rbac.AddRole(&Role{
		Name: "manager",
		Permissions: []Permission{
			{Resource: "user", Action: "read"},
			{Resource: "transaction", Action: "read"},
			{Resource: "balance", Action: "read"},
			{Resource: "audit", Action: "read"},
		},
	})

	// Support - Customer support access
	rbac.AddRole(&Role{
		Name: "support",
		Permissions: []Permission{
			{Resource: "user", Action: "read"},
			{Resource: "transaction", Action: "read"},
			{Resource: "balance", Action: "read"},
		},
	})

	// User - Basic user access
	rbac.AddRole(&Role{
		Name: "user",
		Permissions: []Permission{
			{Resource: "user", Action: "read", ResourceID: "self"},
			{Resource: "transaction", Action: "read", ResourceID: "self"},
			{Resource: "balance", Action: "read", ResourceID: "self"},
			{Resource: "transaction", Action: "create", ResourceID: "self"},
		},
	})
}

// initializeDefaultPolicies creates default access control policies
func (rbac *RBACManager) initializeDefaultPolicies() {
	// User management policies
	rbac.AddPolicy(&Policy{
		Resource: "user",
		Actions:  []string{"create", "read", "update", "delete"},
		Roles:    []string{"admin", "super_admin"},
	})

	// Transaction policies
	rbac.AddPolicy(&Policy{
		Resource: "transaction",
		Actions:  []string{"create", "read", "update", "delete"},
		Roles:    []string{"admin", "super_admin"},
		Conditions: map[string]interface{}{
			"max_amount": 1000000, // 1M limit for admins
		},
	})

	// Balance policies
	rbac.AddPolicy(&Policy{
		Resource: "balance",
		Actions:  []string{"read", "update"},
		Roles:    []string{"admin", "super_admin"},
	})

	// Audit policies
	rbac.AddPolicy(&Policy{
		Resource: "audit",
		Actions:  []string{"read"},
		Roles:    []string{"admin", "manager", "super_admin"},
	})
}

// AddRole adds a new role
func (rbac *RBACManager) AddRole(role *Role) {
	rbac.mu.Lock()
	defer rbac.mu.Unlock()
	rbac.roles[role.Name] = role
}

// GetRole retrieves a role by name
func (rbac *RBACManager) GetRole(name string) (*Role, error) {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	role, exists := rbac.roles[name]
	if !exists {
		return nil, fmt.Errorf("role not found: %s", name)
	}

	return role, nil
}

// AddPolicy adds a new access control policy
func (rbac *RBACManager) AddPolicy(policy *Policy) {
	rbac.mu.Lock()
	defer rbac.mu.Unlock()
	rbac.policies[policy.Resource] = policy
}

// CheckPermission checks if a user with given roles has permission for an action on a resource
func (rbac *RBACManager) CheckPermission(userRoles []string, resource, action string, resourceID interface{}) bool {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	// Check if user has super admin role
	for _, role := range userRoles {
		if role == "super_admin" {
			return true
		}
	}

	// Check policies
	policy, exists := rbac.policies[resource]
	if !exists {
		return false
	}

	// Check if user has required role
	hasRole := false
	for _, userRole := range userRoles {
		for _, requiredRole := range policy.Roles {
			if userRole == requiredRole {
				hasRole = true
				break
			}
		}
		if hasRole {
			break
		}
	}

	if !hasRole {
		return false
	}

	// Check if action is allowed
	actionAllowed := false
	for _, allowedAction := range policy.Actions {
		if allowedAction == "*" || allowedAction == action {
			actionAllowed = true
			break
		}
	}

	if !actionAllowed {
		return false
	}

	// Check conditions if any
	if len(policy.Conditions) > 0 {
		return rbac.checkConditions(policy.Conditions, resourceID)
	}

	return true
}

// checkConditions checks policy conditions
func (rbac *RBACManager) checkConditions(conditions map[string]interface{}, resourceID interface{}) bool {
	// Implement condition checking logic here
	// For now, return true (can be extended with complex condition logic)
	return true
}

// GetUserPermissions returns all permissions for a user with given roles
func (rbac *RBACManager) GetUserPermissions(userRoles []string) []Permission {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	var permissions []Permission
	seen := make(map[string]bool)

	for _, roleName := range userRoles {
		role, exists := rbac.roles[roleName]
		if !exists {
			continue
		}

		for _, permission := range role.Permissions {
			key := fmt.Sprintf("%s:%s", permission.Resource, permission.Action)
			if !seen[key] {
				permissions = append(permissions, permission)
				seen[key] = true
			}
		}
	}

	return permissions
}

// ValidateRole validates if a role name is valid
func (rbac *RBACManager) ValidateRole(roleName string) bool {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	_, exists := rbac.roles[roleName]
	return exists
}

// GetRoles returns all available roles
func (rbac *RBACManager) GetRoles() []string {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	var roles []string
	for roleName := range rbac.roles {
		roles = append(roles, roleName)
	}

	return roles
}
