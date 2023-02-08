// ((i%9 == pos%9 || i/9 == pos/9) || (i/27 == pos/27 && i%9/3 == pos%9/3)) && val != v

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {
	feM := true
Menu:
	for {
		if feM {
			fmt.Printf("\"start\" to play sudoku,\n\"exit\" to exit game\n")
			feM = false
		}
		var iMenu string
		fmt.Printf("command: ")
		fmt.Scanln(&iMenu)
		switch iMenu {
		case "start":
			feG := true
			s := Sudoku{}
			s.createField()
			printField(s.playfield, s.startfield)
		Game:
			for {
				if feG {
					fmt.Println("Use \"show\" to show whole field,\n\"input\" to input digit,\n\"reset\" to reset whole field,\n\"submit\" to submit your solution,\n\"exit\" to return to menu")
					feG = false
				}
				var iGame string
				fmt.Printf("command: ")
				fmt.Scanln(&iGame)
				switch iGame {
				case "show":
					printField(s.playfield, s.startfield)
				case "input":
					printField(s.playfield, s.startfield)
					s.inputNumber()
				case "reset":
					s.playfield = s.startfield
				case "submit":
					id, err := s.submit()
					if err != nil {
						printFieldSelect(s.playfield, s.startfield, id)
						fmt.Println(err)
					} else {
						fmt.Println("You have complited Sudoku!")
					}
				case "exit":
					feG = true
					feM = true
					break Game
				default:
					fmt.Println("Use \"show\" to show whole field,\n\"input\" to input digit,\n\"reset\" to reset whole field,\n\"submit\" to submit your solution,\n\"exit\" to return to menu")
				}
			}
		case "exit":
			break Menu
		default:
			fmt.Printf("\"start\" to play sudoku,\n\"exit\" to exit game\n")
		}
	}

}

type Sudoku struct {
	playfield  [81]uint8
	startfield [81]uint8
	cleanfield [81]uint8
}

func (s *Sudoku) submit() (int, error) {
	var zeros, err bool
	erId := 0
	for idx, val := range s.playfield {
		if val == 0 {
			erId = idx
			zeros = true
			break
		}
	}
	if !zeros {
		for idx, val := range s.playfield {
			for i, v := range s.playfield {
				// fmt.Println(idx, "-", val, "\t", i, "-", v)
				if idx != i && ((i%9 == idx%9 || i/9 == idx/9) || (i/27 == idx/27 && i%9/3 == idx%9/3)) && val == v {
					err = true
					erId = idx
					break
				}
			}

			if err {
				break
			}
		}
	}
	if zeros {
		return erId, errors.New("not all zeros replaced")
	} else if err {
		return erId, errors.New("not sudoku")
	} else {
		return 0, nil
	}
}

func (s *Sudoku) inputNumber() {
	idx, err := s.selectPoint()
	if err != nil {
		fmt.Println(err)
	} else if s.startfield[idx] != 0 {
		printFieldSelect(s.playfield, s.startfield, idx)
		fmt.Println("You can`t change this number")
	} else {
		printFieldSelect(s.playfield, s.startfield, idx)
		var v int
		fmt.Printf("Input digit: ")
		fmt.Scanln(&v)
		s.playfield[idx] = uint8(v)
		printFieldSelect(s.playfield, s.startfield, idx)
	}
}

func (s *Sudoku) createField() {
	// order := createOrder()
	var zeros int
	var input string
Dificulty:
	for {
		fmt.Printf("choose dificulty: 1/2/3: ")
		fmt.Scanln(&input)
		dificulty, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("wrong dificulty!")
		}
		switch dificulty {
		case 1:
			zeros = randomInt(18, 23)
			break Dificulty
		case 2:
			zeros = randomInt(23, 27)
			break Dificulty
		case 3:
			zeros = randomInt(27, 32)
			break Dificulty
		default:
			fmt.Println("wrong dificulty!")
		}

	}

	fmt.Printf("loading")
	field := [81]uint8{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
		4, 5, 6, 7, 8, 9, 1, 2, 3,
		7, 8, 9, 1, 2, 3, 4, 5, 6,
		2, 3, 4, 5, 6, 7, 8, 9, 1,
		5, 6, 7, 8, 9, 1, 2, 3, 4,
		8, 9, 1, 2, 3, 4, 5, 6, 7,
		3, 4, 5, 6, 7, 8, 9, 1, 2,
		6, 7, 8, 9, 1, 2, 3, 4, 5,
		9, 1, 2, 3, 4, 5, 6, 7, 8}

	order := createOrder()
	for idx, val := range field {
		field[idx] = order[val-1]
	}
	fmt.Printf(".")
	s.cleanfield = field
	s.mixLinesInSquares(10)
	fmt.Printf(".")
	s.mixColumsInSquares(10)
	fmt.Printf(".")
	s.mixSquaresHorizonrtal(10)
	fmt.Printf(".")
	s.mixSuqaresVertical(10)
	s.startfield = s.cleanfield
	fmt.Printf(".")
	s.addZeros(zeros)
	s.playfield = s.startfield
	fmt.Printf(".\n\n")
}

func (s *Sudoku) selectPoint() (int, error) {
	var input string
	coltoint := map[rune]int{
		'A': 0,
		'B': 1,
		'C': 2,
		'D': 3,
		'E': 4,
		'F': 5,
		'G': 6,
		'H': 7,
		'I': 8,
	}
	var line, col int
	fmt.Printf("Select point: ")
	fmt.Scanln(&input)
	input = strings.ReplaceAll(input, " ", "")
	rd := regexp.MustCompile("\\d").FindAllStringSubmatch(input, -1)
	rw := regexp.MustCompile("[a-zA-Z]").FindAllStringSubmatch(input, -1)
	if len(rd) > 0 && len(rw) > 0 {
		l, err := strconv.Atoi(rd[0][0])
		if err != nil {
			fmt.Println(err)
		} else {
			line = l
		}
		if c, inMap := coltoint[unicode.ToUpper(rune(rw[0][0][0]))]; inMap {
			col = c
			return 9*(line-1) + col, nil

		} else {
			return 0, (errors.New("column is out of range. please use a-i or A-I letters"))
		}
	} else {
		return 0, (errors.New("unknown input. examles of correct input: \"A1\"; \"b3\"; \"5c\"; \"7H\""))
	}
}

func (s *Sudoku) mixLinesInSquares(rounds int) {
	for round := 0; round < rounds; round++ {
		square := randomInt(0, 2)
		line1, line2 := pickTwoOfThree()
		for symb := 0; symb < 9; symb++ {
			s.cleanfield[27*square+9*line1+symb], s.cleanfield[27*square+9*line2+symb] = s.cleanfield[27*square+9*line2+symb], s.cleanfield[27*square+9*line1+symb]
		}
	}
}

func (s *Sudoku) mixColumsInSquares(rounds int) {
	for round := 0; round < rounds; round++ {
		square := randomInt(0, 2)
		col1, col2 := pickTwoOfThree()
		for symb := 0; symb < 9; symb++ {
			s.cleanfield[square*3+col1+9*symb], s.cleanfield[square*3+col2+9*symb] = s.cleanfield[square*3+col2+9*symb], s.cleanfield[square*3+col1+9*symb]
		}
	}
}

func (s *Sudoku) mixSquaresHorizonrtal(rounds int) {
	for round := 0; round < rounds; round++ {
		square1, square2 := pickTwoOfThree()
		for line := 0; line < 3; line++ {
			for symb := 0; symb < 9; symb++ {
				s.cleanfield[27*square1+9*line+symb], s.cleanfield[27*square2+9*line+symb] = s.cleanfield[27*square2+9*line+symb], s.cleanfield[27*square1+9*line+symb]
			}
		}
	}
}

func (s *Sudoku) mixSuqaresVertical(rounds int) {
	for round := 0; round < rounds; round++ {
		square1, square2 := pickTwoOfThree()
		for col := 0; col < 3; col++ {
			for symb := 0; symb < 9; symb++ {
				s.cleanfield[3*square1+col+9*symb], s.cleanfield[3*square2+col+9*symb] = s.cleanfield[3*square2+col+9*symb], s.cleanfield[3*square1+col+9*symb]
			}
		}
	}
}

func (s *Sudoku) addZeros(ammount int) {
	for i := 0; i < ammount; i++ {
		for {
			idx := randomInt(0, 80)
			if s.startfield[idx] != 0 {
				s.startfield[idx] = 0
				break
			}
		}
	}
}

func printField(field, startfield [81]uint8) {
	fmt.Println("    A  B  C  D  E  F  G  H  I")
	for i := 0; i < 9; i++ {
		fmt.Printf("%d: ", i+1)
		for a := 0; a < 9; a++ {
			idx := 9*i + a
			if startfield[idx] == 0 {
				fmt.Printf("`%d`", field[idx])
			} else {
				fmt.Printf(" %d ", field[idx])
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func printFieldSelect(field, startfield [81]uint8, num int) {
	fmt.Println("\n    A  B  C  D  E  F  G  H  I")
	for i := 0; i < 9; i++ {
		fmt.Printf("%d: ", i+1)
		for a := 0; a < 9; a++ {
			idx := 9*i + a
			if idx == num {
				fmt.Printf("[%d]", field[idx])
			} else if startfield[idx] == 0 {
				fmt.Printf("`%d`", field[idx])
			} else {
				fmt.Printf(" %d ", field[idx])
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func randomInt(min, max int) int {
	time.Sleep(time.Duration(1) * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%(max+1-min) + min
}

func pickTwoOfThree() (int, int) {
	var a, b int
	a = randomInt(0, 2)
	switch a {
	case 0:
		b = randomInt(1, 2)
	case 1:
		act := randomInt(0, 1)
		if act == 0 {
			b = 0
		} else if act == 1 {
			b = 2
		}
	case 2:
		b = randomInt(0, 1)
	}
	return a, b
}

func createOrder() [9]uint8 {
	order := [9]uint8{}
	for order[0]+order[1]+order[2]+order[3]+order[4]+order[5]+order[6]+order[7]+order[8] != 45 {
		idx := randomInt(0, 8)
		val := randomInt(1, 9)
		var stbl bool
		for i := range order {
			if order[i] == uint8(val) {
				stbl = false
				break
			} else {
				stbl = true
			}
		}
		if stbl {
			order[idx] = uint8(val)
		}
	}
	return order
}
