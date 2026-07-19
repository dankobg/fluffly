package dbcustom

type AnalyticsUserStats struct {
	Total   int `json:"total"`
	ByState struct {
		Active   int `json:"active"`
		Inactive int `json:"inactive"`
	} `json:"by_state"`
}

type AnalyticsOrganizationStats struct {
	Total    int `json:"total"`
	ByStatus struct {
		Approved int `json:"approved"`
		Pending  int `json:"pending"`
		Rejected int `json:"rejected"`
	} `json:"by_status"`
}

type AnalyticsAdoptionStats struct {
	Total    int `json:"total"`
	ByStatus struct {
		Approved int `json:"approved"`
		Pending  int `json:"pending"`
		Rejected int `json:"rejected"`
	} `json:"by_status"`
}

type AnalyticsAnimalStats struct {
	Total    int `json:"total"`
	ByStatus struct {
		Adoptable int `json:"adoptable"`
		Pending   int `json:"pending"`
		Adopted   int `json:"adopted"`
		Reserved  int `json:"reserved"`
		Rejected  int `json:"rejected"`
	} `json:"by_status"`
}

type AnalyticsStats struct {
	Users         AnalyticsUserStats
	Organizations AnalyticsOrganizationStats
	Animals       AnalyticsAnimalStats
	Adoptions     AnalyticsAdoptionStats
}

type MyAnalyticsOrganizationStats struct {
	Total    int `json:"total"`
	ByStatus struct {
		Approved int `json:"approved"`
		Pending  int `json:"pending"`
		Rejected int `json:"rejected"`
	} `json:"by_status"`
}

type MyAnalyticsAdoptionStats struct {
	Total    int `json:"total"`
	ByStatus struct {
		Approved int `json:"approved"`
		Pending  int `json:"pending"`
		Rejected int `json:"rejected"`
	} `json:"by_status"`
}

type MyAnalyticsAnimalStats struct {
	Posted struct {
		Total    int `json:"total"`
		ByStatus struct {
			Adoptable int `json:"adoptable"`
			Pending   int `json:"pending"`
			Adopted   int `json:"adopted"`
			Reserved  int `json:"reserved"`
			Rejected  int `json:"rejected"`
		} `json:"by_status"`
	} `json:"posted"`
	Favorites int `json:"favorites"`
	Adoptions int `json:"adoptions"`
}

type MyAnalyticsStats struct {
	Organizations MyAnalyticsOrganizationStats
	Animals       MyAnalyticsAnimalStats
	Adoptions     MyAnalyticsAdoptionStats
}
