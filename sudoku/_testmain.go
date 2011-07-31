package main

import "sudoku"
import "testing"
import __os__ "os"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"sudoku.TestCross", sudoku.TestCross},
	{"sudoku.TestCreateUnitList", sudoku.TestCreateUnitList},
	{"sudoku.TestCreateUnits", sudoku.TestCreateUnits},
	{"sudoku.TestCreatePeers", sudoku.TestCreatePeers},
	{"sudoku.TestGrid_Values", sudoku.TestGrid_Values},
	{"sudoku.TestParser_Grid", sudoku.TestParser_Grid},
	{"sudoku.TestSolve", sudoku.TestSolve},
	{"sudoku.TestSolveHarvestFile", sudoku.TestSolveHarvestFile},
	{"sudoku.TestSolveHard", sudoku.TestSolveHard},
}

var benchmarks = []testing.InternalBenchmark{}

var matchPat string
var matchRe *__regexp__.Regexp

func matchString(pat, str string) (result bool, err __os__.Error) {
	if matchRe == nil || matchPat != pat {
		matchPat = pat
		matchRe, err = __regexp__.Compile(matchPat)
		if err != nil {
			return
		}
	}
	return matchRe.MatchString(str), nil
}

func main() {
	testing.Main(matchString, tests, benchmarks)
}
