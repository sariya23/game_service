package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func main() {
	d := []int{1, 2}
	s1 := sq.Select("*").From("t1").Where(sq.Eq{"id": d})
	s2 := sq.Select("*").From("tw").Where(sq.Eq{"ge": d})

	tagSQL, tagArgs, _ := s1.ToSql() // тут будут ?
	genreSQL, genreArgs, _ := s2.ToSql()

	intersectSQL := fmt.Sprintf("(%s intersect %s)", tagSQL, genreSQL)
	args := append(tagArgs, genreArgs...)

	finalQuery := sq.Select("*").
		From("game").
		Where(sq.Expr("qwe in "+intersectSQL, args...)).Where(fmt.Sprintf("extract(year from %s)=?", "asd")).Limit(4)

	sqlStr, finalArgs, err := finalQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		panic(err)
	}

	finalArgs = append(args, finalArgs...)
	fmt.Println(sqlStr, finalArgs)
}
