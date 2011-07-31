/*
Copyright (C) 2011 Flavio Regis de Arruda

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
*/


package sudoku

import "fmt"
import "strconv"
import "strings"

var digits string = "123456789"
var rows string = "ABCDEFGHI"
var cols string = digits
var squares []string = Cross(rows, cols)
var unitlist [][]string = CreateUnitList(rows, cols)
var units map[string][][]string = CreateUnits(squares, unitlist)
var peers map[string]map[string]bool = CreatePeers(units)

func Cross(A, B string) []string {
	resp := make([]string, len(A)*len(B))
	p := 0
	for a := 0; a < len(A); a++ {
		for b := 0; b < len(B); b++ {
			resp[p] = string(A[a]) + string(B[b])
			p++
		}
	}

	return resp
}

func CreateUnitList(rows, cols string) [][]string {
	resp := make([][]string, len(rows)*3)
	p := 0
	for i := 0; i < len(cols); i++ {
		cr := Cross(rows, string(cols[i]))
		list := make([]string, len(cr))
		for j := 0; j < len(cr); j++ {
			list[j] = cr[j]
		}
		resp[p] = list
		p++
	}
	for i := 0; i < len(rows); i++ {
		cr := Cross(string(rows[i]), cols)
		list := make([]string, len(cr))
		for j := 0; j < len(cr); j++ {
			list[j] = cr[j]
		}
		resp[p] = list
		p++
	}
	rs := []string{`ABC`, `DEF`, `GHI`}
	cs := []string{`123`, `456`, `789`}
	for i := 0; i < len(rs); i++ {
		for j := 0; j < len(cs); j++ {
			cr := Cross(string(rs[i]), string(cs[j]))
			//fmt.Printf("cr: %v\n",cr)
			list := make([]string, len(cr))
			for k := 0; k < len(cr); k++ {
				list[k] = cr[k]
			}
			resp[p] = list
			p++
		}
	}
	return resp
}

func CreateUnits(squares []string, unitlist [][]string) map[string][][]string {
	units := make(map[string][][]string, len(squares))
	c := make(chan bool)
	for _, s := range squares {
		unit := make([][]string, 3)
		i := 0
		for _, u := range unitlist {
			for _, uu := range u {
				if s == uu {
					unit[i] = u
					units[s] = unit
					i++
					break
				}
			}
		}
	}
	close(c)
	return units
}

func CreatePeers(units map[string][][]string) map[string]map[string]bool {
	peers := make(map[string]map[string]bool, len(units))
	c := make(chan bool)
	for k, v := range units {
		peer := make(map[string]bool, 20)
		for _, u := range v {
			for _, uu := range u {
				if k != uu {
					peer[uu] = true
				}
			}
		}
		peers[k] = peer

	}

	close(c)
	return peers
}

func GridValues(grid string) map[string]string {
	values := make(map[string]string, len(grid))
	chars := make([]string, len(grid))
	for i := 0; i < len(grid); i++ {
		str := grid[i : i+1]
		if strings.Contains(digits, str) || strings.Contains("0.", str) {
			chars[i] = str
		}
	}
	if len(chars) != 81 {
		panic("Invalid grid size: expected grid of size 81 found grid of size " + strconv.Itoa(len(chars)))
	}

	for i := 0; i < len(grid); i++ {
		values[squares[i]] = chars[i]
	}
	return values
}

//  Eliminate removes d from values[s]; propagate when values or places <= 2.
//  Return values, except return False if a contradiction is detected.
func Eliminate(values map[string]string, s string, d string) map[string]string {
	if !strings.Contains(values[s], d) {
		return values
	}
	values[s] = strings.Replace(values[s], d, "", -1)
	// (1) If a square s is reduced to one value d2, then eliminate d2 from the peers.
	if len(values[s]) == 0 {
		return nil
	} else if len(values[s]) == 1 {
		d2 := values[s]
		for s2, _ := range peers[s] {
			if Eliminate(values, s2, d2) == nil {
				return nil
			}
		}
	}
	// (2) If a unit u is reduced to only one place for a value d, then put it there.
	for _, unit := range units[s] {
		if len(unit) == 0 {
			panic("Error ZERO length unit!")
		}
		dplaces := []string{}
		for _, square := range unit {
			if strings.Contains(values[square], d) {
				dplaces = append(dplaces, square)
			}
		}
		if len(dplaces) == 0 {
			return nil
		}
		if len(dplaces) == 1 {
			if Assign(values, dplaces[0], d) == nil {
				return nil
			}
		}
	}
	return values
}
// Assign 
// Eliminate all the other values (except d) from values[s] and propagate.
func Assign(values map[string]string, s string, d string) map[string]string {
	other_values := strings.Replace(values[s], d, "", -1)
	for i := 0; i < len(other_values); i++ {
		if Eliminate(values, s, string(other_values[i])) == nil {
			return nil
		}
	}
	return values
}


// ParseGrid convert grid to a dict of possible values, {square: digits}, or
// return False if a contradiction is detected.
func ParseGrid(grid string) (values map[string]string) {
	values = make(map[string]string, len(squares))
	for _, s := range squares {
		values[s] = digits
	}
	//k, v 
	for s, d := range GridValues(grid) {
		if strings.Contains(digits, d) {
			values = Assign(values, s, d)
			if values == nil {
				return nil
			}
		}
	}
	return values
}

// Using depth-first search and propagation, try all possible values.
func Search(values map[string]string) map[string]string {
	//fmt.Println("deep:",deep)
	if values == nil {
		return nil
	}
	solved := true
	for s, _ := range values {
		if len(values[s]) != 1 {
			solved = false
		}
	}
	if solved {
		return values
	}
	// Chose the unfilled square s with the fewest possibilities
	min := len(digits) + 1
	sq := ""
	for _, s := range squares {
		l := len(values[s])
		if l > 1 {
			if l < min {
				sq = s
				min = l
			}
		}
	}
	ch := make(chan map[string]string)
	for _, d := range values[sq] {
		go func(dd int) {
			nvalues := cloneValues(values)
			v := Search(Assign(nvalues, sq, string(dd)))
			if v != nil {
				ch <- v
			}
		}(d)
	}

	return <-ch
}

func cloneValues(m map[string]string) map[string]string {
	nm := make(map[string]string, len(m))
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

func Solve(grid string) map[string]string {
	return Search(ParseGrid(grid))
}

func Display(values map[string]string) {
	for r, row := range rows {
		for c, col := range digits {
			if c == 3 || c == 6 {
				fmt.Printf("| ")
			}
			fmt.Printf("%v ", values[string(row)+string(col)])
		}
		fmt.Println()
		if r == 2 || r == 5 {
			fmt.Println("------+-------+-------")
		}
	}
}
