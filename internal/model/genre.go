package model

type Genre struct {
	GenreID   int64
	GenreName string
}

func GenreIDs(g []Genre) []int64 {
	if len(g) == 0 {
		return nil
	}
	res := make([]int64, 0, len(g))
	for _, v := range g {
		res = append(res, v.GenreID)
	}
	return res
}

func GenreNames(g []Genre) []string {
	if len(g) == 0 {
		return nil
	}
	res := make([]string, 0, len(g))
	for _, v := range g {
		res = append(res, v.GenreName)
	}
	return res
}
