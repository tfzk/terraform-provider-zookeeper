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

func NewMissingEnvVarError(varName string) *MissingEnvVarError {
	return &MissingEnvVarError{varName}
}

// NonSeqZNodeCannotEndWithPathSeparatorError returned when a non sequenzial ZNode ends with a path separator.
type NonSeqZNodeCannotEndWithPathSeparatorError struct {
	path string
}

func (e *NonSeqZNodeCannotEndWithPathSeparatorError) Error() string {
	return fmt.Sprintf("non-sequential ZNode cannot have path '%s' because it ends in '%c'", e.path, zNodePathSeparator)
}

func NewNonSeqZNodeCannotEndWithPathSeparatorError(path string) *NonSeqZNodeCannotEndWithPathSeparatorError {
	return &NonSeqZNodeCannotEndWithPathSeparatorError{path}
}

// CannotUpdateDoesNotExistError returned when a ZNode cannot be updated because it does not exist.
type CannotUpdateDoesNotExistError struct {
	path string
}

func (e *CannotUpdateDoesNotExistError) Error() string {
	return fmt.Sprintf("failed to update ZNode '%s': does not exist", e.path)
}

func NewCannotUpdateDoesNotExistError(path string) *CannotUpdateDoesNotExistError {
	return &CannotUpdateDoesNotExistError{path}
}
