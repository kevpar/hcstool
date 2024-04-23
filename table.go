package main

import "fmt"

type colInfo struct {
	header string
	format string
}

func printTable[T any](colInfo []colInfo, rowData []T, rowExtract func(T) []any) error {
	var (
		cols = len(colInfo)
		max  = make([]int, 0, cols)
		data = make([][]string, 0, len(rowData))
	)
	for _, ci := range colInfo {
		max = append(max, len(ci.header))
	}
	for _, rd := range rowData {
		d := rowExtract(rd)
		if len(d) != cols {
			return fmt.Errorf("row did not match header column count")
		}
		values := make([]string, 0, cols)
		for i := 0; i < cols; i++ {
			s := fmt.Sprintf(colInfo[i].format, d[i])
			if l := len(s); l > max[i] {
				max[i] = l
			}
			values = append(values, s)
		}
		data = append(data, values)
	}
	for i := range colInfo {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%[1]*[2]s", max[i], colInfo[i].header)
	}
	fmt.Printf("\n")
	for _, d := range data {
		for i := range d {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%[1]*[2]s", max[i], d[i])
		}
		fmt.Printf("\n")
	}
	return nil
}
