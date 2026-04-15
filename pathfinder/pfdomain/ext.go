package pfdomain

type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type Extensions struct {
	PersistedQuery *PersistedQuery `json:"persistedQuery,omitempty"`

	/*
			[RequestIds] has structure like:
			{
		 		"<endpoint path>": {
		 			"<name of api used>": "<uuid4>/<request-scoped id>"
		    	}
			}
	*/
	RequestIds map[string]map[string]string `json:"requestIds,omitempty"`
}
