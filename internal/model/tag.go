package model

type Tag struct {
	TagID   int64
	TagName string
}

func TagIDs(t []Tag) []int64 {
	if t == nil {
		return nil
	}
	res := make([]int64, 0, len(t))
	for _, v := range t {
		res = append(res, v.TagID)
	}
	return res
}
