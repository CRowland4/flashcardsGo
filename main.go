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
	var cards []flashcard
	cardCount := readInt("Input the number of cards:")

	for i := 1; i <= cardCount; i++ {
		cards = append(cards, createCard(i))
	}

	for _, card := range cards {
		answerPrompt := fmt.Sprintf("Print the definition of \"%s\":", card.term)
		answer := readLine(answerPrompt)
		checkAnswer(card, answer)
	}

	return
}

func checkAnswer(card flashcard, answer string) {
	if answer == card.definition {
		fmt.Println("Correct!")
		return
	}

	fmt.Printf("Wrong. The right answer is \"%s\"\n", card.definition)
	return
}

func createCard(i int) (card flashcard) {
	termPrompt := fmt.Sprintf("The term for card #%d:", i)
	definitionPrompt := fmt.Sprintf("The definition for card #%d:", i)

	card.term = readLine(termPrompt)
	card.definition = readLine(definitionPrompt)

	return card
}

func readLine(prompt string) (line string) {
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func readInt(prompt string) (num int) {
	fmt.Println(prompt)
	fmt.Scanln(&num)
	return num
}
