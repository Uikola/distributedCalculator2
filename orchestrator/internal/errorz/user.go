package errorz

import "errors"

var ErrInvalidLogin = errors.New("invalid login")
var ErrInvalidPassword = errors.New("invalid password")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidSigningMethod = errors.New("invalid signing method")
var ErrInvalidClaimsType = errors.New("invalid claims type")
var ErrInvalidOperation = errors.New("invalid operation")
var ErrInvalidOperationTime = errors.New("invalid operation time")
