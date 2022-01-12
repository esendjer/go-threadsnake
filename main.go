package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tty "github.com/mattn/go-tty"
)

const (
	tail = "o"
	body = "O"
	head = "@"
)

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
}

func genArr(size int) *[]string {
	r := make([]string, size+2)
	r[0] = fmt.Sprintf("+%s+", strings.Repeat("-", (size*2)))
	r[size+1] = r[0]
	r[1] = fmt.Sprintf("|%s|", strings.Repeat(" ", (size*2)))
	for i := 2; i < size+1; i++ {
		r[i] = r[1]
	}
	return &r
}

func frug(a *[]string) []int {
	return []int{rand.Intn(len((*a)[0])-2) + 1, rand.Intn(len((*a))-2) + 1}
}

func cleanScreen() {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[0;0H")
}

func printHelp() {
	fmt.Println(`Using: go-threadsnake [1st arg] [2nd arg]

Where:
	1st arg - Size of playing field, is a number from 5 to 40
	          default value = 10, means 10 rows and (10 * 2 =) 20 columns

	2nd arg - Speed of moving the snake, is a number from 1 to 20
	          default value = 2, means (60 sec / 2 =) 30 sec for a step`)
}

func sayHello() {
	cleanScreen()
	fmt.Println(`

	Game will start in 3 sec!

	w - UP
	s - DOWN
	a - LEFT
	d - RIGHT

	Good luck!`)
}

func main() {
	as := 10
	speed := time.Second / 2

	// Checking input args and applaing new values for as and speed
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printHelp()
			return
		}
		a, err := strconv.Atoi(os.Args[1])
		if err == nil && a >= 5 && a <= 40 {
			as = a
		}
		if len(os.Args) == 3 {
			b, err := strconv.Atoi(os.Args[2])
			if err == nil && b >= 1 && b <= 20 {
				speed = time.Second / time.Duration(b)
			}
		}
	}

	arr := genArr(as)
	sec := 0
	scor := 0
	pos := [][]int{
		{3, 1}, // Head
		{2, 1}, // Body
		{1, 1}, // Tail
	}
	fru := frug(arr)
	dir := "d"
	// Direction could be:
	// w - UP
	// s - DOWN
	// a - LEFT
	// d - RIGHT

	ch := make(chan string)

	// Print playing field
	fmt.Println((*arr))

	// Reading pressed keys
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	go func() {
		for {
			k, err := tty.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			ch <- string(k)
		}
	}()

	sayHello()
	time.Sleep(time.Second * 3)

	for {
		cleanScreen()

		// Checking pressed keys
		select {
		case k := <-ch:
			if dir == "d" && (k == "s" || k == "w") {
				dir = k
			}
			if dir == "a" && (k == "s" || k == "w") {
				dir = k
			}
			if dir == "w" && (k == "a" || k == "d") {
				dir = k
			}
			if dir == "s" && (k == "a" || k == "d") {
				dir = k
			}
		default:
			pe := pos[:len(pos)-1]

			// Drawing fruit
			if pos[0][0] == fru[0] && pos[0][1] == fru[1] {
				scor++
				fc := true
				for fc {
					fru = frug(arr)
					fc = func() bool {
						for _, v := range pos[1:] {
							if fru[0] == v[0] && fru[1] == v[1] {
								return true
							}
						}
						return false
					}()
				}
				pe = pos
			}

			// Moving snake
			if dir == "d" {
				pos = append([][]int{{pos[0][0] + 1, pos[0][1]}}, pe...)
			}
			if dir == "a" {
				pos = append([][]int{{pos[0][0] - 1, pos[0][1]}}, pe...)
			}
			if dir == "s" {
				pos = append([][]int{{pos[0][0], pos[0][1] + 1}}, pe...)
			}
			if dir == "w" {
				pos = append([][]int{{pos[0][0], pos[0][1] - 1}}, pe...)
			}
			cros := func() bool {
				for _, v := range pos[1:] {
					if v[0] == pos[0][0] && v[1] == pos[0][1] {
						return true
					}
				}
				return false
			}()

			// Checking WIN or LOSS
			if pos[0][0] < 1 || pos[0][1] < 1 || pos[0][0] > len((*arr)[0])-2 || pos[0][1] > len((*arr))-2 || cros {
				fmt.Printf("Game over! Your score is %d!\n", scor)
				return
			}
			if len(pos) >= len((*arr))*len((*arr)[0])/2 {
				fmt.Printf("You WIN!\nYour score is %d!\n", scor)
				return
			}

			// Drawing playing field and snake
			for i, v := range *arr {
				str := v
				if fru[1] == i {
					str = str[:fru[0]] + "$" + str[fru[0]+1:]
				}
				for j, a := range pos {
					r := body
					if j == 0 {
						r = head
					}
					if j == len(pos)-1 {
						r = tail
					}
					if a[1] == i {
						str = str[:a[0]] + r + str[a[0]+1:]
					}

				}
				// Printing playing field and snake by line
				fmt.Println(str)
			}
			// Printing statistics
			fmt.Printf("Steps: %d\n", sec)
			fmt.Printf("Score: %d\n", scor)

			sec++
			time.Sleep(speed)
		}
	}
}
