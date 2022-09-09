package responseRepository

type Result struct {
	Status int         `json:"status"`
	Datas  interface{} `json:"datas"`
	Errors interface{} `json:"errors"`
}
