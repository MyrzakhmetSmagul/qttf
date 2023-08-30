package jsonsaver

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveJson(path string, body interface{}) error {
	fmt.Printf("Saving json to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to save json file: %w", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(body)
	return nil
}
