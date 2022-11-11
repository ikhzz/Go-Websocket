package models

type (
	MainResponse struct {
		Status 			bool 		`json:"status"`
		StatusCode 		int 		`json:"-"`
		Message 		string		`json:"message"`
		Data 			interface{}	`json:"data"`
		TotalGroupData	int			`json:"total_group_data"`
		Payload			PayloadModel `json:"-"`
		Connections []WebSocketConnection `json:"-"`
	}

	PayloadModel struct {
		Id	string
	}
)