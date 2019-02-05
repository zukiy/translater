package v1

type (
	// TranslateRequest ...
	TranslateRequest struct {
		Word string `json:"word"`
		Lang string `json:"lang"`
	}

	// TranslateResponse ...
	TranslateResponse struct {
		Results []string `json:"results"`
	}

	// Word ...
	Word struct {
		Text string `json:"text"`
		Lang string `json:"lang"`
	}

	// SaveRequest ...
	SaveRequest struct {
		Word      Word `json:"word"`
		Translate Word `json:"translate"`
	}

	// ListRequest
	ListRequest struct {
		PerPage int `json:"per_page"`
		Page    int `json:"page"`
	}
)
