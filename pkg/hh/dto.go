package hh

type GetVacanciesResponse struct {
	Items []struct {
		ID                     string `json:"id"`
		Premium                bool   `json:"premium"`
		Name                   string `json:"name"`
		Department             any    `json:"department"`
		HasTest                bool   `json:"has_test"`
		ResponseLetterRequired bool   `json:"response_letter_required"`
		Area                   struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"area"`
		Salary      any `json:"salary"`
		SalaryRange any `json:"salary_range"`
		Type        struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
		Address           any    `json:"address"`
		ResponseURL       any    `json:"response_url"`
		SortPointDistance any    `json:"sort_point_distance"`
		PublishedAt       string `json:"published_at"`
		CreatedAt         string `json:"created_at"`
		Archived          bool   `json:"archived"`
		ApplyAlternateURL string `json:"apply_alternate_url"`
		ShowLogoInSearch  any    `json:"show_logo_in_search,omitempty"`
		ShowContacts      bool   `json:"show_contacts"`
		InsiderInterview  any    `json:"insider_interview"`
		URL               string `json:"url"`
		AlternateURL      string `json:"alternate_url"`
		Relations         []any  `json:"relations"`
		Employer          struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			URL          string `json:"url"`
			AlternateURL string `json:"alternate_url"`
			LogoUrls     struct {
				Num90    string `json:"90"`
				Num240   string `json:"240"`
				Original string `json:"original"`
			} `json:"logo_urls"`
			VacanciesURL         string `json:"vacancies_url"`
			CountryID            int    `json:"country_id"`
			AccreditedItEmployer bool   `json:"accredited_it_employer"`
			Trusted              bool   `json:"trusted"`
		} `json:"employer"`
		Snippet struct {
			Requirement    string `json:"requirement"`
			Responsibility string `json:"responsibility"`
		} `json:"snippet"`
		Contacts any `json:"contacts"`
		Schedule struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"schedule"`
		WorkingDays          []any `json:"working_days"`
		WorkingTimeIntervals []any `json:"working_time_intervals"`
		WorkingTimeModes     []any `json:"working_time_modes"`
		AcceptTemporary      bool  `json:"accept_temporary"`
		FlyInFlyOutDuration  []any `json:"fly_in_fly_out_duration"`
		WorkFormat           []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"work_format"`
		WorkingHours []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"working_hours"`
		WorkScheduleByDays []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"work_schedule_by_days"`
		AcceptLaborContract bool  `json:"accept_labor_contract"`
		CivilLawContracts   []any `json:"civil_law_contracts"`
		NightShifts         bool  `json:"night_shifts"`
		ProfessionalRoles   []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"professional_roles"`
		AcceptIncompleteResumes bool `json:"accept_incomplete_resumes"`
		Experience              struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"experience"`
		Employment struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"employment"`
		EmploymentForm struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"employment_form"`
		Internship     bool `json:"internship"`
		AdvResponseURL any  `json:"adv_response_url"`
		IsAdvVacancy   bool `json:"is_adv_vacancy"`
		AdvContext     any  `json:"adv_context"`
		Branding       struct {
			Type   string `json:"type"`
			Tariff any    `json:"tariff"`
		} `json:"branding,omitempty"`
	} `json:"items"`
	Found        int    `json:"found"`
	Pages        int    `json:"pages"`
	Page         int    `json:"page"`
	PerPage      int    `json:"per_page"`
	Clusters     any    `json:"clusters"`
	Arguments    any    `json:"arguments"`
	Fixes        any    `json:"fixes"`
	Suggests     any    `json:"suggests"`
	AlternateURL string `json:"alternate_url"`
}

type GetVacanciesRequest struct {
	Text              string   `json:"text,omitempty" url:"text,omitempty"`
	SearchField       []string `json:"search_field,omitempty" url:"search_field,omitempty"`
	Area              []string `json:"area,omitempty" url:"area,omitempty"`
	ProfessionalRole  []string `json:"professional_role,omitempty" url:"professional_role,omitempty"`
	Industry          []string `json:"industry,omitempty" url:"industry,omitempty"`
	Experience        string   `json:"experience,omitempty" url:"experience,omitempty"`
	Employment        []string `json:"employment,omitempty" url:"employment,omitempty"`
	Schedule          []string `json:"schedule,omitempty" url:"schedule,omitempty"`
	Salary            int      `json:"salary,omitempty" url:"salary,omitempty"`
	Currency          string   `json:"currency,omitempty" url:"currency,omitempty"`
	OnlyWithSalary    bool     `json:"only_with_salary,omitempty" url:"only_with_salary,omitempty"`
	Period            int      `json:"period,omitempty" url:"period,omitempty"`
	DateFrom          string   `json:"date_from,omitempty" url:"date_from,omitempty"`
	DateTo            string   `json:"date_to,omitempty" url:"date_to,omitempty"`
	EmployerId        []string `json:"employer_id,omitempty" url:"employer_id,omitempty"`
	ExcludeEmployerId []string `json:"exclude_employer_id,omitempty" url:"exclude_employer_id,omitempty"`
	Label             []string `json:"label,omitempty" url:"label,omitempty"`
	OrderBy           string   `json:"order_by,omitempty" url:"order_by,omitempty"`
	Page              int      `json:"page,omitempty" url:"page,omitempty"`
	PerPage           int      `json:"per_page,omitempty" url:"per_page,omitempty"`
	TopLat            float64  `json:"top_lat,omitempty" url:"top_lat,omitempty"`
	BottomLat         float64  `json:"bottom_lat,omitempty" url:"bottom_lat,omitempty"`
	LeftLng           float64  `json:"left_lng,omitempty" url:"left_lng,omitempty"`
	RightLng          float64  `json:"right_lng,omitempty" url:"right_lng,omitempty"`
}
