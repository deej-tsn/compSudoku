package models

import "fmt"

type (
	Square struct {
		Value     int
		Confirmed bool
		Potential []int
	}

	Row []*Square

	Grid []Row
)

func (grid Grid) Print() {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}
