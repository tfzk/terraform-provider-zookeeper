package provider

import (
	"fmt"
)

// ACLPermissionsValueOutOfRangeError returned when an attempt is made
// to set ACL permissions using an out of range integer.
type ACLPermissionsValueOutOfRangeError struct {
	permValue int
}

func (e *ACLPermissionsValueOutOfRangeError) Error() string {
	return fmt.Sprintf("ACL permissions value %d is out of int32 range", e.permValue)
}

// NewACLPermissionsValueOutOfRangeError creates a new ACLPermissionsValueOutOfRangeError.
//
// permValue is the value of the ACL permissions.
//
// Example:
//
//	NewACLPermissionsValueOutOfRangeError(1000000000)
func NewACLPermissionsValueOutOfRangeError(permValue int) *ACLPermissionsValueOutOfRangeError {
	return &ACLPermissionsValueOutOfRangeError{permValue}
}
