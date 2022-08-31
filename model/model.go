package model

type (
	DataStruct struct {
		ReceivedBy string            `json:"receivedBy"`
		Histories  []HistoriesStruct `json:"histories"`
	}

	HistoriesStruct struct {
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt"`
		Formatted   struct {
			CreatedAt string `json:"createdAt"`
		} `json:"formatted"`
	}
)
