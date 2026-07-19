package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/persistence/dbcustom"
)

func AnalyticsStatsToResponse(data dbcustom.AnalyticsStats) api.AnalyticsStats {
	resp := api.AnalyticsStats{
		Adoptions: api.AnalyticsAdoptionsStats{
			ByStatus: struct {
				Approved int "json:\"approved\""
				Pending  int "json:\"pending\""
				Rejected int "json:\"rejected\""
			}{
				Approved: data.Adoptions.ByStatus.Approved,
				Pending:  data.Adoptions.ByStatus.Pending,
				Rejected: data.Adoptions.ByStatus.Rejected,
			},
			Total: data.Adoptions.Total,
		},
		Users: api.AnalyticsUserStats{
			ByState: struct {
				Active   int "json:\"active\""
				Inactive int "json:\"inactive\""
			}{
				Active:   data.Users.ByState.Active,
				Inactive: data.Users.ByState.Inactive,
			},
			Total: data.Users.Total,
		},
		Animals: api.AnalyticsAnimalStats{
			ByStatus: struct {
				Adoptable int "json:\"adoptable\""
				Adopted   int "json:\"adopted\""
				Pending   int "json:\"pending\""
				Rejected  int "json:\"rejected\""
				Reserved  int "json:\"reserved\""
			}{
				Adoptable: data.Animals.ByStatus.Adoptable,
				Adopted:   data.Animals.ByStatus.Adopted,
				Pending:   data.Animals.ByStatus.Pending,
				Rejected:  data.Animals.ByStatus.Rejected,
				Reserved:  data.Animals.ByStatus.Reserved,
			},
			Total: data.Animals.Total,
		},
		Organizations: api.AnalyticsOrganizationStats{
			ByStatus: struct {
				Approved int "json:\"approved\""
				Pending  int "json:\"pending\""
				Rejected int "json:\"rejected\""
			}{
				Approved: data.Organizations.ByStatus.Approved,
				Rejected: data.Organizations.ByStatus.Rejected,
				Pending:  data.Organizations.ByStatus.Pending,
			},
			Total: data.Organizations.Total,
		},
	}

	return resp
}

func MyAnalyticsStatsToResponse(data dbcustom.MyAnalyticsStats) api.MyAnalyticsStats {
	resp := api.MyAnalyticsStats{
		Adoptions: api.MyAnalyticsAdoptionsStats{
			ByStatus: struct {
				Approved int "json:\"approved\""
				Pending  int "json:\"pending\""
				Rejected int "json:\"rejected\""
			}{
				Approved: data.Adoptions.ByStatus.Approved,
				Pending:  data.Adoptions.ByStatus.Pending,
				Rejected: data.Adoptions.ByStatus.Rejected,
			},
			Total: data.Adoptions.Total,
		},
		Organizations: api.MyAnalyticsOrganizationStats{
			ByStatus: struct {
				Approved int "json:\"approved\""
				Pending  int "json:\"pending\""
				Rejected int "json:\"rejected\""
			}{
				Approved: data.Organizations.ByStatus.Approved,
				Pending:  data.Organizations.ByStatus.Pending,
				Rejected: data.Organizations.ByStatus.Rejected,
			},
			Total: data.Organizations.Total,
		},

		Animals: api.MyAnalyticsAnimalStats{
			Adoptions: data.Animals.Adoptions,
			Favorites: data.Animals.Favorites,
			Posted: struct {
				ByStatus struct {
					Adoptable int "json:\"adoptable\""
					Adopted   int "json:\"adopted\""
					Pending   int "json:\"pending\""
					Rejected  int "json:\"rejected\""
					Reserved  int "json:\"reserved\""
				} "json:\"by_status\""
				Total int "json:\"total\""
			}{
				ByStatus: struct {
					Adoptable int "json:\"adoptable\""
					Adopted   int "json:\"adopted\""
					Pending   int "json:\"pending\""
					Rejected  int "json:\"rejected\""
					Reserved  int "json:\"reserved\""
				}{
					Adoptable: data.Animals.Posted.ByStatus.Adoptable,
					Adopted:   data.Animals.Posted.ByStatus.Adopted,
					Pending:   data.Animals.Posted.ByStatus.Pending,
					Rejected:  data.Animals.Posted.ByStatus.Rejected,
					Reserved:  data.Animals.Posted.ByStatus.Reserved,
				},
				Total: data.Animals.Posted.Total,
			},
		},
	}

	return resp
}
