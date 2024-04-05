package entity

type CResource struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Occupied   string `json:"occupied"`
}
