package models

import "RIP/internal/app/ds"

type FinesListWithRes struct {
	Fines    *[]ds.Fines `json:"id"`
	ResCount int         `json:"resCount"`
	ResID    int         `json:"resId"`
}

type FineWithCount struct {
	Fines *ds.Fines `json:"fines"`
	Count int       `json:"count"`
}
