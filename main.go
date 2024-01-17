package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cardCount := readInt("Input the number of cards:")

	cards := createCards(cardCount)
	quizCards(cards)

	return
}

func quizCards(cards map[string]string) {
	for key, val := range cards {
		answerPrompt := fmt.Sprintf("Print the definition of \"%s\"", key)
		answer := readLine(answerPrompt)

		if answer == val {
			fmt.Println("Correct!")
			continue
		} else if definitionExists(answer, cards) {
			val2 := getTermFor(answer, cards)
			fmt.Printf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n", val, val2)
		} else {
			fmt.Printf("Wrong. The right answer is \"%s\".\n", val)
		}
	}

	return
}

func getTermFor(definition string, cards map[string]string) (term string) {
	for key, val := range cards {
		if val == definition {
			term = key
			break
		}
	}

	return term
}

func createCards(cardCount int) (cards map[string]string) {
	cards = make(map[string]string)
	for i := 1; i <= cardCount; i++ {
		term := getNewTerm(cards, i)
		definition := getNewDefinition(cards, i)
		cards[term] = definition
	}

	return cards
}

func getNewDefinition(cards map[string]string, cardNum int) (definition string) {
	initialPrompt := fmt.Sprintf("The definition for card #%d:", cardNum)
	definition = readLine(initialPrompt)

	for {
		if !definitionExists(definition, cards) {
			return definition
		}

		newPrompt := fmt.Sprintf("The definition \"%s\" already exists. Try again:", definition)
		definition = readLine(newPrompt)
	}
}

func definitionExists(definition string, cards map[string]string) (exists bool) {
	for _, val := range cards {
		if definition == val {
			return true
		}
	}

	return false
}

func getNewTerm(cards map[string]string, cardNum int) (term string) {
	initialPrompt := fmt.Sprintf("The term for card #%d:", cardNum)
	term = readLine(initialPrompt)

	for {
		if _, ok := cards[term]; !ok {
			return term
		}

		newPrompt := fmt.Sprintf("The term \"%s\" already exists. Try again:", term)
		term = readLine(newPrompt)
	}
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
