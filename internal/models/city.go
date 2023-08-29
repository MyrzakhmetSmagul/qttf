package models

import "fmt"

type City struct {
	Id        int
	Name      string
	Hyperlink string
}

func (c City) ToInsertScript() string {
	return fmt.Sprintf("INSERT INTO city (city_name, city_link) VALUES('%s', '%s') ON CONFLICT DO NOTHING;\n", c.Name, c.Hyperlink)
}
