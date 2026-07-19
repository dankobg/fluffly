package shared

import "fmt"

func AuthzIdentityID(id string) string {
	return fmt.Sprintf("identity:%s", id)
}

func AuthzCourierMessageID(id string) string {
	return fmt.Sprintf("courier_message:%s", id)
}

func AuthzSchemaID(id string) string {
	return fmt.Sprintf("schema:%s", id)
}

func AuthzSessionID(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func AuthzCountryID(id int64) string {
	return fmt.Sprintf("country:%d", id)
}

func AuthzAnimalID(id int64) string {
	return fmt.Sprintf("animal:%d", id)
}

func AuthzAnimalTypeID(id int64) string {
	return fmt.Sprintf("animal_type:%d", id)
}

func AuthzAnimalSpecieID(id int64) string {
	return fmt.Sprintf("animal_specie:%d", id)
}

func AuthzBreedID(id int64) string {
	return fmt.Sprintf("breed:%d", id)
}

func AuthzOrganizationID(id int64) string {
	return fmt.Sprintf("organization:%d", id)
}

func AuthzAnalyticsID(userID string) string {
	return fmt.Sprintf("analytics:%s", userID)
}

func AuthzAdoptionID(id int64) string {
	return fmt.Sprintf("adoptions:%d", id)
}
