package client

import (
	"fmt"
)

// MissingEnvVarError returned when a necessary Environment variable is missing.
type MissingEnvVarError struct {
	varName string
}

func (e *MissingEnvVarError) Error() string {
	return "missing environment variable: " + e.varName
}

// NewMissingEnvVarError creates a new MissingEnvVarError.
//
// varName is the name of the missing environment variable.
//
// Example:
//
//	NewMissingEnvVarError("ZOOKEEPER_SERVERS")
func NewMissingEnvVarError(varName string) *MissingEnvVarError {
	return &MissingEnvVarError{varName}
}

// NonSeqZNodeCannotEndWithPathSeparatorError returned when a non sequenzial ZNode ends with a path separator.
type NonSeqZNodeCannotEndWithPathSeparatorError struct {
	path string
}

func (e *NonSeqZNodeCannotEndWithPathSeparatorError) Error() string {
	return fmt.Sprintf(
		"non-sequential ZNode cannot have path '%s' because it ends in '%c'",
		e.path,
		zNodePathSeparator,
	)
}

// NewNonSeqZNodeCannotEndWithPathSeparatorError creates a new NonSeqZNodeCannotEndWithPathSeparatorError.
//
// path is the path of the non-sequential ZNode.
//
// Example:
//
//	NewNonSeqZNodeCannotEndWithPathSeparatorError("/path/to/znode")
func NewNonSeqZNodeCannotEndWithPathSeparatorError(
	path string,
) *NonSeqZNodeCannotEndWithPathSeparatorError {
	return &NonSeqZNodeCannotEndWithPathSeparatorError{path}
}

// CannotUpdateDoesNotExistError returned when a ZNode cannot be updated because it does not exist.
type CannotUpdateDoesNotExistError struct {
	path string
}

func (e *CannotUpdateDoesNotExistError) Error() string {
	return fmt.Sprintf("failed to update ZNode '%s': does not exist", e.path)
}

// NewCannotUpdateDoesNotExistError creates a new CannotUpdateDoesNotExistError.
//
// path is the path of the non-sequential ZNode.
//
// Example:
//
//	NewCannotUpdateDoesNotExistError("/path/to/znode")
func NewCannotUpdateDoesNotExistError(path string) *CannotUpdateDoesNotExistError {
	return &CannotUpdateDoesNotExistError{path}
}
