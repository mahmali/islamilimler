package model

type Hadis struct {
	Kitap  string   `json:"kitap"`
	Konu   string   `json:"konu"`
	Numara int      `json:"numara"`
	Metin  []string `json:"metin"`
}
