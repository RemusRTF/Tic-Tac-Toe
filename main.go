package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func clearBoard(m map[int]string) {
	for i := 1; i <= 9; i++ {
		m[i] = "-"
	}
}

func displayBoard(m map[int]string) {
	for j := 2; j >= 0; j-- {
		for i := 1; i <= 3; i++ {
			fmt.Print("	")
			fmt.Print(m[i+3*j])
		}
		fmt.Println()
	}
}

func getScore(m map[int]string) int {
	for i := 0; i < 3; i++ {
		if m[1+i*3] == m[2+i*3] && m[2+i*3] == m[3+i*3] {

			if m[1+i*3] == "X" {
				return 1
			} else if m[1+i*3] == "O" {
				return -1
			}

		}
	}

	for i := 0; i < 3; i++ {
		if m[1+i] == m[4+i] && m[4+i] == m[7+i] {

			if m[1+i] == "X" {
				return 1
			}
			if m[1+i] == "O" {
				return -1
			}

		}
	}

	if m[1] == m[5] && m[5] == m[9] {
		if m[1] == "X" {
			return 1
		}
		if m[1] == "O" {
			return -1
		}
	}

	if m[3] == m[5] && m[5] == m[7] {
		if m[3] == "X" {
			return 1
		}
		if m[3] == "O" {
			return -1
		}
	}

	return 0
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func playAI(available map[int]bool, m map[int]string, activeChar string) int {
	//results := make(map[int]int)
	highestValue := 1

	//depth: 8 is max
	depth := 9

	start := time.Now()

	position, _ := doSomething(available, m, activeChar, highestValue, depth)

	elapsed := time.Since(start)

	fmt.Println("Calculating took: ", elapsed)

	return position
}

var counter int

func doSomething(available map[int]bool, m map[int]string, activeChar string, highestValue int, depth int) (position int, highestYet int) {
	counter++

	if len(available) > 1 && depth > 0 {
		//fmt.Println("doing Something")
		depth--

		bestValue := 0
		highestPos := 0

		first := true //Helps that the highestValue gets set as a base value

		for a := range available {

			newM := make(map[int]string)
			for tile, value := range m {
				newM[tile] = value
			}

			newM[a] = activeChar

			/*clearConsole()
			fmt.Println("Depth: ", depth)
			displayBoard(newM)*/

			if activeChar == "X" {
				//Check, if won by next turn
				if getScore(newM) == highestValue {
					//fmt.Println("Reached highest Value at:", a)
					return a, highestValue
				}
			} else {
				//Check, if won by next turn
				if getScore(newM) == -highestValue {
					//fmt.Println("Reached lowest Value at:", a)
					return a, -highestValue //----------------------------- false --------------------------------
				}
			}

			if first {
				_, highestYet := doSomething(removeElement(available, a), newM, changePlayer(activeChar), highestValue, depth)
				bestValue = highestYet
				highestPos = a
				first = false
			} else {
				_, highestYet := doSomething(removeElement(available, a), newM, changePlayer(activeChar), highestValue, depth)

				if activeChar == "X" {
					if highestYet >= bestValue {
						bestValue = highestYet
						highestPos = a
					}
				} else {
					if highestYet <= bestValue {
						bestValue = highestYet
						highestPos = a
					}
				}
			}

			//fmt.Println(bestValue)
			if depth > 6 {
				clearConsole()
				displayBoard(newM)
				time.Sleep(1000)
			}

		}

		return highestPos, bestValue
	}

	var bestValue int
	var highestPos int
	first := true

	for a := range available {

		newM := make(map[int]string)
		for tile, value := range m {
			newM[tile] = value
		}
		newM[a] = activeChar

		highestYet := getScore(newM)

		if first {
			bestValue = highestYet
			highestPos = a
			first = false
		}

		if activeChar == "X" {
			//Check, if won by next turn
			if highestYet == highestValue {
				return a, highestValue
			}

		} else {
			//Check, if won by next turn
			if highestYet == -highestValue {
				return a, -highestValue //--------------------- false -----------------------------------
			}
		}

		if activeChar == "X" && highestYet > bestValue {
			bestValue = highestYet
			highestPos = a
		} else if highestYet < bestValue {
			bestValue = highestYet
			highestPos = a
		}
	}
	return highestPos, bestValue
}

func removeElement(available map[int]bool, value int) map[int]bool {
	newAvailable := make(map[int]bool)
	for number := range available {
		newAvailable[number] = true
	}
	delete(newAvailable, value)
	return newAvailable
}

func changePlayer(activeChar string) string {
	if activeChar == "X" {
		return "O"
	}
	return "X"
}

func main() {
	clearConsole()
	fmt.Println("Starting...")

	m := make(map[int]string)
	available := make(map[int]bool)
	for i := 1; i <= 9; i++ {
		available[i] = true
	}

	clearBoard(m)

	//available := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	reader := bufio.NewReader(os.Stdin)

	activeChar := "O"

	for getScore(m) == 0 && len(available) > 0 {
		//clearConsole()
		activeChar = changePlayer(activeChar)
		intInput := 0

		if activeChar == "O" {

			//clearConsole()
			displayBoard(m)

			goodInput := false

			for !goodInput {
				fmt.Print("Enter a number: ")
				input, _ := reader.ReadString('\n')
				input = input[:len(input)-2]
				intInput, _ = strconv.Atoi(input)
				for a := range available {
					if a == intInput {
						goodInput = true
						break
					}
				}
			}
		} else {
			counter = 0
			intInput = playAI(available, m, activeChar)
			fmt.Println("Branches: ", counter)
			fmt.Println("AI Output: ", intInput)
		}

		m[intInput] = activeChar
		delete(available, intInput)
		//fmt.Println("Score: ", getScore(m))
	}
	//clearConsole()
	displayBoard(m)
	fmt.Println()
	if getScore(m) == 0 {
		activeChar = "Nobody"
	}
	fmt.Println(activeChar + " won the game!")
}
