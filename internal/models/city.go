package models

type City struct {
	Id        int
	Name      string
	Hyperlink string
}

type CityList struct {
	TotalCount int     `json:"total_count"`
	TotalPage  int     `json:"total_page"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Cities     []*City `json:"cities"`
}
