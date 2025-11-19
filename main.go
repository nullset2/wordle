package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	reset    = "\033[0m"
	white    = "\033[97m"
	bold     = "\033[1m"
	BgBlack  = "\033[40m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	StatsFilename := filepath.Join(home, "bin", "stats")
	dictionaryFilename := filepath.Join(home, "bin", "dictionary.txt")
	input, err := os.Open(dictionaryFilename)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	reader := bufio.NewReader(input)
	dict := make(map[int]string)
	i := 0

	for {
		line, _, _ := reader.ReadLine()
		if len(line) == 0 {
			break
		}

		s := string(line)
		dict[i] = s
		i++
	}

	stats := []int{0, 0}

	if _, err := os.Stat(StatsFilename); os.IsNotExist(err) {
		err := os.WriteFile(StatsFilename, []byte("0 0"), 0644)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		line, err := os.ReadFile(StatsFilename)
		if err != nil {
			panic(err)
		}
		parts := strings.Split(string(line), " ")
		won, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		lost, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		stats = []int{won, lost}
	}

	rand := rand.Intn(len(dict))
	word := dict[rand]
	won := false
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Wordle")

	for tries := 6; tries > 0; tries-- {
		var line string
		fmt.Print("Guess? >")
		for scanner.Scan() {
			line = strings.ToUpper(scanner.Text())

			if len(line) != 5 {
				fmt.Printf("Five letter words only. Guess? >")
			} else {
				break
			}
		}

		results := make([]int, 0)

		for i, c := range word {
			if c == rune(line[i]) {
				results = append(results, 1)
			} else if strings.ContainsRune(word, rune(line[i])) {
				results = append(results, 0)
			} else {
				results = append(results, -1)
			}
		}

		score := 0

		for i, x := range results {
			score += x
			if x == -1 {
				fmt.Print(BgBlack + white + " " + string(line[i]) + " " + reset)
			} else if x == 0 {
				fmt.Print(BgYellow + white + " " + string(line[i]) + " " + reset)
			} else {
				fmt.Print(BgGreen + white + " " + string(line[i]) + " " + reset)
			}
		}
		fmt.Print("\n")

		if score == 5 {
			won = true
			break
		}
	}

	if won {
		fmt.Println("Congratulations! You Won!")
		stats[0]++
	} else {
		fmt.Printf("The word was: %s\n", word)
		fmt.Println("Try again!")
		stats[1]++
	}
	fmt.Printf("Total Games Played: %v\n", stats[0]+stats[1])
	fmt.Printf("Games Won: %v\n", stats[0])
	fmt.Printf("Games Lost: %v\n", stats[1])
	fmt.Printf("Success Rate: %.2f%%\n", (float64(stats[0])/float64(stats[0]+stats[1]))*100)
	err = os.WriteFile(StatsFilename, []byte(fmt.Sprintf("%d %d", stats[0], stats[1])), 0644)
	if err != nil {
		panic(err)
	}
}
