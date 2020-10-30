package cache

type Keyword struct {
	Pos int    `json:"pos,omitempty"`
	Key string `json:"key,omitempty"`
}

type Item struct {
	V       interface{} `json:"V,omitempty"`
	Expired int64       `json:"Expired,omitempty"`
}
