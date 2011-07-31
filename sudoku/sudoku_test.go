package sudoku

import (
	"testing"
	"fmt"
	"os"
	"bufio"
)


func TestCross(t *testing.T) {
	squares := Cross(rows, cols)
	if len(squares) != 81 {
		t.Errorf("Wrong number of squares. Expected %v found %v\n", 81, len(squares))
	}
	fmt.Printf("squares: %v\n", squares)
}

func TestCreateUnitList(t *testing.T) {
	unitList := CreateUnitList(rows, cols)
	if len(unitList) != 27 {
		t.Errorf("Wrong number of unitLists. Expected %v found %v\n", len(cols)*3, len(unitList))
	}
	fmt.Printf("unitList: %v\n", unitList)
}

func TestCreateUnits(t *testing.T) {
	squares := Cross(rows, cols)
	unitList := CreateUnitList(rows, cols)
	units := CreateUnits(squares, unitList)
	if len(units) != len(squares) {
		t.Errorf("Wrong number of units. Expected %v found %v\n", len(squares), len(units))
	}
	fmt.Printf("units['G5']: %v\n", units["G5"])
}


func TestCreatePeers(t *testing.T) {
	squares := Cross(rows, cols)
	unitList := CreateUnitList(rows, cols)
	units := CreateUnits(squares, unitList)
	peers := CreatePeers(units)
	if len(peers) != len(units) {
		t.Errorf("Wrong number of peers. Expected %v found %v\n", len(units), len(peers))
	}
	fmt.Printf("peers['G5']: %v\n", peers["G5"])
}

func TestGrid_Values(t *testing.T) {
	g := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	values := GridValues(g)
	if len(values) != len(squares) {
		t.Errorf("Wrong number of grid_values. Expected %v found %v\n", len(squares), len(values))
		return
	}
	//test some values.
	if values["A1"] != "4" || values["B2"] != "3" || (values["C3"] != "." && values["C3"] != "0") {
		t.Errorf("Invalid grid_values for test %v\n", values)
		return
	}
	fmt.Println("grid_values:", values)
}

func TestParser_Grid(t *testing.T) {
	g := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	values := ParseGrid(g)
	fmt.Println("parser_grid:", values)
}

func TestSolve(t *testing.T) {
	g := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	values := Solve(g)
	fmt.Println("resolution:", values)
	g = "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
	values = Solve(g)
	fmt.Println("resolution:")
	Display(values)
}

func TestSolveHarvestFile(t *testing.T) {
	fmt.Println("---		TestSolveHarvest		---")
	//HARD
	file, err := os.Open("/home/regis/Documents/Projects/go/workspace/GoNuts/src/pkg/sudoku/hardest.txt")
	if err != nil {
		panic("Nao foi possível abrir arquivo hardest.txt : " + err.String())
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, _, lerr := reader.ReadLine()
		if lerr != nil {
			break
		}
		fmt.Println("problem(hard): ",string(line))
		values := Solve(string(line))
		fmt.Println("resolution(hard):")
		Display(values)
	}
}

func TestSolveHard(t *testing.T) {
    //line := ".....6....59.....82....8....45........3........6..3.54...325..6.................."   VERY VERY HARD!!!
    line := "....361.......2......15..296..8.....29....6.5..3...7.4..6.........34.91.9.....47."
	values := Solve(line)
	fmt.Println("problem(hard): ", line)
	fmt.Println("solution(hard) : ")
	Display(values)
}
