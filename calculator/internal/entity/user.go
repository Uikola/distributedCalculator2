package entity

type User struct {
	ID             uint   `json:"id"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	Addition       uint   `json:"addition"`
	Subtraction    uint   `json:"subtraction"`
	Multiplication uint   `json:"multiplication"`
	Division       uint   `json:"division"`
}
