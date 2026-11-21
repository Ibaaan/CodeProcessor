package domain

import "errors"

var ErrNotFound = errors.New("Resource is not found")
var ErrTaskCreation = errors.New("task creation failed")
