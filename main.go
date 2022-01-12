package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	tty "github.com/mattn/go-tty"
)

const (
	tail = "o"
	body = "O"
	head = "@"
)

var version = "v0.1"

var options struct {
	Version bool `short:"v" long:"version" description:"Show version\n"`
	Size    int  `short:"s" long:"size" description:"Set the size of playing field, from 5 to 40\n default value = 10, means 10 rows and (10 * 2 =) 20 columns\n"`
	Tempo   int  `short:"t" long:"tempo" description:"Set the tempo of moving the snake, from 1 to 20\n default value = 2, means (60 sec / 2 =) 30 sec for a step\n"`
}

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
	_, err := flags.Parse(&options)
	if err != nil {
		e, _ := err.(*flags.Error)
		if e.Type != flags.ErrHelp {
			os.Exit(1)
		}
		os.Exit(0)
	}
	if options.Version {
		fmt.Printf("go-threadsnake version: %s\n", version)
		os.Exit(0)
	}
}

// function to generate playing field
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

// function to generate the position of a fruit
func frug(a *[]string) []int {
	return []int{rand.Intn(len((*a)[0])-2) + 1, rand.Intn(len((*a))-2) + 1}
}

func cleanScreen() {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[0;0H")
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
	pfs := 10 // playing field size
	speed := time.Second / 2

	// Checking input args and applaing new values for pfs and speed
	if s := options.Size; s >= 5 && s <= 40 {
		pfs = s
	}
	if t := options.Tempo; t >= 1 && t <= 20 {
		speed = time.Second / time.Duration(t)
	}

	arr := genArr(pfs) // playing field
	sec := 0
	scor := 0
	pos := [][]int{ // snake
		{3, 1}, // Head
		{2, 1}, // Body
		{1, 1}, // Tail
	}
	fru := frug(arr) // fruit
	dir := "d"       // direction could be:
	// w - UP
	// s - DOWN
	// a - LEFT
	// d - RIGHT

	// channel for direction
	ch := make(chan string)

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

		// checking pressed keys and set new direction
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

			// drawing fruit
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

			// moving snake
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

			// checking WIN or LOSS
			if pos[0][0] < 1 || pos[0][1] < 1 || pos[0][0] > len((*arr)[0])-2 || pos[0][1] > len((*arr))-2 || cros {
				fmt.Printf("Game over! Your score is %d!\n", scor)
				return
			}
			if len(pos) >= len((*arr))*len((*arr)[0])/2 {
				fmt.Printf("You WIN!\nYour score is %d!\n", scor)
				return
			}

			// drawing playing field and snake
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
				// printing playing field and snake by line
				fmt.Println(str)
			}
			// printing statistics
			fmt.Printf("Steps: %d\n", sec)
			fmt.Printf("Score: %d\n", scor)

			sec++
			time.Sleep(speed)
		}
	}
}
