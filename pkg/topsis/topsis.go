package topsis

const (
	TotalCriteria    = 5
	TotalSubCriteria = 4
	TotalAlternative = 7
)

var (
	AnswerToFloat = map[string]float32{
		"9":   9.0,
		"7":   7.0,
		"5":   5.0,
		"3":   2.0,
		"1":   1.0,
		"1/3": 1.0 / 3.0,
		"1/5": 1.0 / 5.0,
		"1/7": 1.0 / 7.0,
		"1/9": 1.0 / 9.0,
	}
	AlternativeToRow = map[string]int{
		"Multimedia":                    0,
		"Teknik Elektro":                1,
		"Teknik Listrik":                2,
		"Tata Busana":                   3,
		"Tata Boga":                     4,
		"Tata Kecantikan":               5,
		"Teknik Kendaraan Ringan/Motor": 6,
	}
	CriteriaToCol = map[string]int{
		"total_open_jobs":              0,
		"salary":                       1,
		"entrepreneurship_opportunity": 2,
		"interest":                     3,
		"supporting_facilites":         4,
	}
)
