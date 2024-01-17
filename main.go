package main

import (
	"bufio"
	"fmt"
	"os"
)

type flashcard struct {
	term       string
	definition string
}

func main() {
	card := createCard()
	answer := getInputLine()

	checkAnswer(card, answer)
}

func checkAnswer(card flashcard, answer string) {
	if answer == card.definition {
		fmt.Println("You're right!")
		return
	}

	fmt.Println("wrong WRONG WROOOOOOOONG")
	return
}

func createCard() (card flashcard) {
	card.term = getInputLine()
	card.definition = getInputLine()

	return card
}

func getInputLine() (line string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}
