package models

// Country は国を表す構造体です
type Country struct {
	Name string `json:"name"`
}

// Capital は首都を表す構造体です
type Capital struct {
	Name string `json:"name"`
}

// CountryCapitalRelation は国と首都の関係を表す構造体です
type CountryCapitalRelation struct {
	Country Country `json:"country"`
	Capital Capital `json:"capital"`
}
