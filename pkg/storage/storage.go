package storage

import (
	"github.com/JTGlez/GoWeb-IT/pkg/models"
)

var (
	Store     = make(map[int]models.Product)
	LastID    int
	CodeIndex = make(map[string]int)
)
