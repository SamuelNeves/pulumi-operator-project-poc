package responses

type ListBucketsResponse struct {
	IDs []string `json:"ids"`
}

type BucketResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
