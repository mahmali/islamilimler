package main

type Hadis struct {
	Konu   string   `json:"konu"`
	Numara int      `json:"numara"`
	Metin  []string `json:"metin"`
}
