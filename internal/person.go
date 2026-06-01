package internal

type Person struct {
	Name     string		`json:"name"`
	Year     int		`json:"year"`
	Month    int		`json:"month"`
	Day      int		`json:"day"`
	Hour     int		`json:"hour"`
	Minute   int		`json:"minute"`
	Second   int		`json:"second"`
	Lat      float64	`json:"lat"`
	Lon      float64	`json:"lon"`
	TZOffset int		`json:"tzoffset"`
}

