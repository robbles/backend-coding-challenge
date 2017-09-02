package models

type Result struct {
	Name  string  `json:"name"`
	Lat   float64 `json:"latitude"`
	Long  float64 `json:"longitude"`
	Score float64 `json:"score"`
}

type ResultsByScore []Result

func (a ResultsByScore) Len() int           { return len(a) }
func (a ResultsByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ResultsByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }
