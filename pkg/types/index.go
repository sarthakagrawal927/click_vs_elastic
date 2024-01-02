package types

type Person struct {
	ID       int     `json:"id"`
	Age      int     `json:"age"`
	Name     string  `json:"name"`
	Height   float64 `json:"height"`
	Weight   int     `json:"weight"`
	IsActive int     `json:"is_active"`
}
