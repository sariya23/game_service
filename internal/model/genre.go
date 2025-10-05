package model

type Genre struct {
	GenreID   int64
	GenreName string
}

func GenreIDs(g []Genre) []int64 {
	if g == nil {
		return nil
	}
	res := make([]int64, 0, len(g))
	for _, v := range g {
		res = append(res, v.GenreID)
	}
	return res
}
