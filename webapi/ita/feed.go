package ita

type EmptyPayload struct{}

type FeedResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Report Report `json:"report"`
	Paging Paging `json:"paging"`
}

type Paging struct {
	Page        int `json:"page"`
	PageSize    int `json:"pageSize"`
	StartRecord int `json:"startRecord"`
	Previous    int `json:"previous"`
	Next        int `json:"next"`
}

type Report struct {
	ReportName     string         `json:"reportName"`
	RenderType     string         `json:"renderType"`
	ReportMetadata ReportMetadata `json:"reportMetadata"`
	ReportData     []ReportData   `json:"reportData"`
}

type ReportMetadata struct {
	TotalRows int `json:"totalRows"`
}

type ReportData map[string]Column

type Column struct {
	Value any `json:"value"`
}

// func GetFeed() error {
// 	var response FeedResponse

// 	resp, err := GenerateReportTemplate(os.Getenv("ITA_API_URI"), os.Getenv("ITA_FEED_ID"), os.Getenv("ITA_TOKEN"), 0, 10000)
// 	if err != nil {
// 		return fmt.Errorf("cannot get feed report: %v", err)
// 	}

// 	if err := json.Unmarshal(resp, &response); err != nil {
// 		return fmt.Errorf("cannot unmarshall feed report: %v", err)
// 	}
// }
