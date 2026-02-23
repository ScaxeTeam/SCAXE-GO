package permission

import "strings"

const (
	DefaultOp    = "op"
	DefaultNotOp = "notop"
	DefaultTrue  = "true"
	DefaultFalse = "false"
)

type Permission struct {
	Name        string
	Description string
	Default     string
	Children    map[string]bool
}

func NewPermission(name, description, defaultValue string) *Permission {
	return &Permission{
		Name:        name,
		Description: description,
		Default:     defaultValue,
		Children:    make(map[string]bool),
	}
}

type PermissionManager struct {
	permissions map[string]*Permission
}

var GlobalManager = NewPermissionManager()

func NewPermissionManager() *PermissionManager {
	return &PermissionManager{
		permissions: make(map[string]*Permission),
	}
}

func (m *PermissionManager) AddPermission(perm *Permission) {
	m.permissions[strings.ToLower(perm.Name)] = perm
}

func (m *PermissionManager) GetPermission(name string) *Permission {
	return m.permissions[strings.ToLower(name)]
}

func (m *PermissionManager) HasPermission(name string, isOp bool) bool {
	perm := m.GetPermission(name)

	var defaultValue string
	if perm != nil {
		defaultValue = perm.Default
	} else {
		defaultValue = DefaultOp
	}

	switch defaultValue {
	case DefaultTrue:
		return true
	case DefaultFalse:
		return false
	case DefaultOp:
		return isOp
	case DefaultNotOp:
		return !isOp
	default:
		return isOp
	}
}

func (m *PermissionManager) GetPermissions() map[string]*Permission {
	return m.permissions
}
