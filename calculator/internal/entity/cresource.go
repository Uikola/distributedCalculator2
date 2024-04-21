package entity

type CResource struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Expression        string `json:"expression"`
	Occupied          bool   `json:"occupied"`
	OrchestratorAlive bool   `json:"orchestrator_alive"`
}
