package errorz

import "errors"

var ErrInvalidExpression = errors.New("invalid expression")
var ErrExpressionNotFound = errors.New("expression not found")
var ErrAccessForbidden = errors.New("access forbidden")
var ErrNoExpressions = errors.New("no expressions")
var ErrEvaluationInProgress = errors.New("expression evaluation in progress")
var ErrEvaluation = errors.New("evaluation error")
