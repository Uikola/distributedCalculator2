package entity

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddExpressionRequest struct {
	Expression string `json:"expression"`
}

type UpdateOperationRequest struct {
	Operation string `json:"operation"`
	Time      uint   `json:"time"`
}
