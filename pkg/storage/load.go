package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/JTGlez/GoWeb-IT/pkg/models"
)

func LoadProducts(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	for _, product := range products {
		Store[product.ID] = product
		CodeIndex[product.CodeValue] = product.ID
		if product.ID > LastID {
			LastID = product.ID
		}
	}

	return nil
}
