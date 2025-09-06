package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/persistence/dbtype"
)

func AnimalToResponse(data model.Animal) api.Animal {
	return api.Animal{
		ID:              data.ID,
		Age:             api.AnimalAge(data.Age),
		Description:     data.Description,
		Distance:        data.Distance,
		Gender:          (*api.AnimalGender)(data.Gender),
		Hermaphrodite:   data.Hermaphrodite,
		ImageURL:        data.ImageURL,
		Name:            data.Name,
		Size:            api.AnimalSize(data.Size),
		Status:          (*api.AnimalStatus)(data.Status),
		StatusChangedAt: data.StatusChangedAt,
		Properties:      &data.Properties,
		AdoptedAt:       data.AdoptedAt,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}

func AnimalWithJoinDataToResponse(data dbtype.AnimalWithJoinData) api.Animal {
	resp := AnimalToResponse(data.Animal)
	resp.Type = api.AnimalType{
		ID:        data.Type.ID,
		Name:      data.Type.Name,
		CreatedAt: data.Type.CreatedAt,
		UpdatedAt: data.Type.UpdatedAt,
	}
	resp.Species = api.AnimalSpecies{
		AnimalTypeID: data.Species.AnimalTypeID,
		ID:           data.Species.ID,
		Name:         data.Species.Name,
		CreatedAt:    data.Species.CreatedAt,
		UpdatedAt:    data.Species.UpdatedAt,
	}
	resp.Breeds = make([]api.Breed, len(data.Breeds))
	for i, breed := range data.Breeds {
		resp.Breeds[i] = AnimalBreedWithJounDataToResp(breed)
	}
	resp.Tags = make([]api.AnimalTag, len(data.Tags))
	for i, tag := range data.Tags {
		resp.Tags[i] = AnimalTagToResp(tag)
	}
	resp.Photos = make([]api.AnimalPhoto, len(data.Photos))
	for i, photo := range data.Photos {
		resp.Photos[i] = AnimalPhotoToResp(photo)
	}
	resp.Videos = make([]api.AnimalVideo, len(data.Videos))
	for i, video := range data.Videos {
		resp.Videos[i] = AnimalVideoToResp(video)
	}

	organizationWithJoinData := OrganizationWithJoinDataToResponse(dbtype.OrganizationWithJoinData{
		Organization: data.Organization.Organization,
		WorkHour:     data.Organization.WorkHour,
		Contact:      data.Organization.Contact,
	})
	resp.Organization = &organizationWithJoinData

	return resp
}
