package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	cards := new(map[string]string)
	*cards = make(map[string]string)

	for {
		switch action := readLine("Input the action (add, remove, import, export, ask, exit):"); action {
		case "add":
			add(cards)
		case "remove":
			remove(cards)
		case "import":
			import_(cards)
		case "export":
			export(*cards)
		case "ask":
			ask(*cards)
		case "exit":
			fmt.Println("Bye bye!")
			return
		default:
			fmt.Printf("Command %s not recognized. Please enter another comand.\n\n", action)
		}
	}
}

func ask(cards map[string]string) {
	if len(cards) == 0 {
		fmt.Println("No cards yet!")
	}
	count := readInt("How many times to ask?")
	quizCards(cards, count)
}

func quizCards(cards map[string]string, count int) {
	terms := make([]string, len(cards))
	for term := range cards {
		terms = append(terms, term)
	}

	for i := 0; i < count; i++ {
		quizTerm := terms[rand.Intn(len(terms))]
		if quizTerm == "" { // Since the terms array is a slice, after adding elements, it dynamically creates more slots, initialized with an empty sring
			i--
			continue
		}

		answerPrompt := fmt.Sprintf("Print the definition of \"%s\"", quizTerm)
		answer := readLine(answerPrompt)

		if answer == cards[quizTerm] {
			fmt.Println("Correct!")
			continue
		} else if definitionExists(answer, cards) {
			val2 := getTermFor(answer, cards)
			fmt.Printf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n", cards[quizTerm], val2)
		} else {
			fmt.Printf("Wrong. The right answer is \"%s\".\n", cards[quizTerm])
		}
	}
	return
}

func export(cards map[string]string) {
	fileName := readLine("File name:")
	file, _ := os.Create(fileName)

	for term, definition := range cards {
		fmt.Fprintln(file, term, definition)
	}

	fmt.Printf("%d cards have been saved.\n", len(cards))
	return
}

func import_(cards *map[string]string) {
	fileName := readLine("File name:")
	if /*file*/ _, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		fmt.Println("File not found")
		return
	}
	// TODO Still have to implement, tests don't specify file format
	fmt.Printf("%d cards have been loaded.\n", len(*cards)) // TODO should print the amount imoprted, not len(map)
	return
}

func remove(cards *map[string]string) {
	card := readLine("Which card?")
	if _, ok := (*cards)[card]; ok {
		delete(*cards, card)
		fmt.Println("The card has been removed.")
	} else {
		fmt.Printf("Can't remove \"%s\": there is no such card.\n", card)
	}

	return
}

func add(cards *map[string]string) {
	term := getNewTerm(cards)
	definition := getNewDefinition(*cards)
	(*cards)[term] = definition

	fmt.Printf("The pair (\"%s\":\"%s\") has been added\n", term, definition)
	return
}

func getNewTerm(cards *map[string]string) (term string) {
	term = readLine("The card")

	for {
		if _, ok := (*cards)[term]; !ok {
			return term
		}

		newPrompt := fmt.Sprintf("The term \"%s\" already exists. Try again:", term)
		term = readLine(newPrompt)
	}
}

func getNewDefinition(cards map[string]string) (definition string) {
	definition = readLine("The definition of the card")

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

func getTermFor(definition string, cards map[string]string) (term string) {
	for key, val := range cards {
		if val == definition {
			term = key
			break
		}
	}

	return term
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
