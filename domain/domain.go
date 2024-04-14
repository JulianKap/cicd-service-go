package domain

type Response struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	AllRowCount int `json:"all_row_count"`
}
