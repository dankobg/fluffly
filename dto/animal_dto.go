package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
)

func AnimalToResponse(data models.Animal, upl media.Uploader) api.Animal {
	getURL := func(name *string, kind string, upl media.Uploader) string {
		if name == nil {
			return ""
		}

		u, err := upl.URL(*name, kind)
		if err != nil {
			return ""
		}

		return u
	}

	var properties *dbcustom.Properties

	jsonProperties := data.Properties.Ptr()
	if jsonProperties != nil {
		properties = &jsonProperties.Val
	}

	return api.Animal{
		ID:             data.ID,
		OrganizationID: data.OrganizationID.Ptr(),
		Age:            api.AnimalAge(data.Age),
		Description:    data.Description.Ptr(),
		Distance:       data.Distance.Ptr(),
		Gender:         (*api.AnimalGender)(data.Gender.Ptr()),
		Hermaphrodite:  data.Hermaphrodite,
		ImageSmallURL:  getURL(&data.ImageObjectRefSmall, data.ImageObjectKind, upl),
		ImageMediumURL: getURL(&data.ImageObjectRefMedium, data.ImageObjectKind, upl),
		ImageLargeURL:  getURL(&data.ImageObjectRefLarge, data.ImageObjectKind, upl),
		ImageFullURL:   getURL(&data.ImageObjectRefFull, data.ImageObjectKind, upl),
		Name:           data.Name,
		Size:           api.AnimalSize(data.Size),
		Status:         api.AnimalStatus(data.Status),
		Properties:     properties,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func AnimalWithJoinDataToResponse(data dbtype.AnimalWithJoinData, upl media.Uploader) api.Animal {
	resp := AnimalToResponse(data.Animal, upl)
	resp.Type = AnimalTypeToResponse(data.Type)
	resp.Specie = AnimalSpeciesMinToResponse(data.Specie)

	if data.Breeds.Val != nil {
		breeds := make([]api.AnimalBreed, len(*data.Breeds.Val))
		for i, breed := range *data.Breeds.Val {
			breeds[i] = AnimalBreedWithJoinDataToResponse(breed)
		}

		resp.Breeds = &breeds
	}

	if data.Tags.Val != nil {
		tags := make([]api.AnimalTag, len(*data.Tags.Val))
		for i, tag := range *data.Tags.Val {
			tags[i] = AnimalTagToResponse(tag)
		}

		resp.Tags = &tags
	}

	if data.Photos.Val != nil {
		photos := make([]api.AnimalPhoto, len(*data.Photos.Val))
		for i, photo := range *data.Photos.Val {
			photos[i] = AnimalPhotoToResponse(photo, upl)
		}

		resp.Photos = &photos
	}

	if data.Videos.Val != nil {
		videos := make([]api.AnimalVideo, len(*data.Videos.Val))
		for i, video := range *data.Videos.Val {
			videos[i] = AnimalVideoToResponse(video, upl)
		}

		resp.Videos = &videos
	}

	resp.Likes = data.Likes
	resp.Liked = data.Liked.Ptr()

	if data.AdoptionID.IsValue() {
		resp.AdoptionID = data.AdoptionID.Ptr()
	}

	if data.Microchip != nil && data.Microchip.ID != 0 {
		resp.Microchip = new(MicrochipToResponse(*data.Microchip))
	}

	if data.Organization != nil && data.Organization.ID != 0 {
		organizationWithJoinData := OrganizationWithJoinDataToResponse(dbtype.OrganizationWithJoinData{
			Organization:          *data.Organization,
			WorkHour:              data.OrganizationWorkHour,
			Contact:               data.OrganizationContact,
			ContactAddress:        data.OrganizationContactAddress,
			ContactAddressCountry: data.OrganizationContactAddressCountry,
		}, upl)
		resp.Organization = &organizationWithJoinData
	}

	return resp
}
