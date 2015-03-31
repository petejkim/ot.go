package ot

type Range struct {
	Anchor int `json:"anchor"`
	Head   int `json:"head"`
}

type Selection struct {
	Ranges []Range `json:"ranges"`
}
