package dbtype

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

type ListAnimalBreedsFilters struct {
	api.ListAnimalBreedsParams
}

type AnimalBreedWithJoinData struct {
	models.Breed
	Primary bool
}
