package errorz

import "errors"

var ErrNoAvailableResources = errors.New("all resources is occupied")
