package sheets

import (
	"fmt"
	"qttf/internal/models"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

func GetCitiesFromSheet(resp *sheets.Spreadsheet) ([]models.City, error) {
	if len(resp.Sheets) > 0 && len(resp.Sheets[0].Data) > 0 && len(resp.Sheets[0].Data[0].RowData) > 0 {
		cities := make([]models.City, 0)
		citiesMap := make(map[string]bool)
		for _, row := range resp.Sheets[0].Data[0].RowData[1:] {

			name := row.Values[2].FormattedValue
			if _, ok := citiesMap[name]; ok {
				continue
			}
			citiesMap[name] = true
			city := models.City{}
			city.Name = name
			city.Hyperlink = row.Values[2].Hyperlink
			cities = append(cities, city)
		}
		return cities, nil
	}
	return nil, fmt.Errorf("GetCitiesFromSheet: data from resp not valid")
}

func GetPlayersFromSheet(resp *sheets.Spreadsheet) ([]models.Player, error) {
	if len(resp.Sheets) > 0 && len(resp.Sheets[0].Data) > 0 && len(resp.Sheets[0].Data[0].RowData) > 0 {
		players := make([]models.Player, len(resp.Sheets[0].Data[0].RowData[1:]))
		for i, row := range resp.Sheets[0].Data[0].RowData[1:] {
			fullName := strings.Split(row.Values[1].FormattedValue, " ")
			players[i].Surname = fullName[0]
			if len(fullName) > 1 {
				players[i].Name = fullName[1]
			}

			players[i].Hyperlink = row.Values[1].Hyperlink
			players[i].City.Name = row.Values[2].FormattedValue
			players[i].City.Hyperlink = row.Values[2].Hyperlink
		}
		return players, nil
	}
	return nil, fmt.Errorf("GetPlayersFromSheet: data from resp not valid")
}

func GetRatingFromSheet(resp *sheets.Spreadsheet) ([]models.Rating, error) {
	if len(resp.Sheets) > 0 && len(resp.Sheets[0].Data) > 0 && len(resp.Sheets[0].Data[0].RowData) > 0 {
		rating := make([]models.Rating, len(resp.Sheets[0].Data[0].RowData[1:]))
		for i, row := range resp.Sheets[0].Data[0].RowData[1:] {
			num, err := strconv.Atoi(row.Values[0].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("GetRatingFromSheet: %w", err)
			}

			rating[i].Position = num
			fullName := strings.Split(row.Values[1].FormattedValue, " ")
			rating[i].Player.Surname = fullName[0]
			if len(fullName) > 1 {
				rating[i].Player.Name = fullName[1]
			}

			rating[i].Player.Hyperlink = row.Values[1].Hyperlink
			rating[i].Player.City.Name = row.Values[2].FormattedValue
			rating[i].Player.City.Hyperlink = row.Values[2].Hyperlink
			num, err = strconv.Atoi(row.Values[3].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("GetRatingFromSheet: %w", err)
			}

			rating[i].Rating = num
			t, err := time.Parse("2.1.2006", row.Values[4].FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("GetRatingFromSheet: %w", err)
			}

			rating[i].UpdateTime = t.Format("2006/01/02")
		}
		return rating, nil
	}
	return nil, fmt.Errorf("GetRatingFromSheet: data from resp not valid")
}
