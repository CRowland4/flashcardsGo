package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Card struct {
	term       string
	definition string
	missed     int
}

const (
	menuPrompt      = "Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):"
	defaultCLImport = ""
	defaultCLExport = ""
)

var sessionLog strings.Builder

func main() {
	var cards []Card
	importFile := flag.String("import_from", defaultCLImport, "Enter a file name to import cards")
	exportFile := flag.String("export_to", defaultCLExport, "Enter a file name to export cards to on exit")
	flag.Parse()
	checkForCLImport(&cards, *importFile)

	for {
		switch action := readLine("\n" + menuPrompt); action {
		case "add":
			add(&cards)
		case "remove":
			remove(&cards)
		case "import":
			import_(&cards, readLine("File name:"))
		case "export":
			export(cards, *exportFile)
		case "ask":
			ask(&cards)
		case "exit":
			logPrintln("Bye bye!")
			export(cards, *exportFile)
			return
		case "log":
			log_()
		case "hardest card":
			hardestCard(cards)
		case "reset stats":
			resetStats(&cards)
		default:
			msg := fmt.Sprintf("Command %s not recognized. Please enter another command.", action)
			logPrintln(msg)
		}
	}
}

func checkForCLImport(cards *[]Card, importFile string) {
	if importFile != defaultCLImport {
		import_(cards, importFile)
	}
	return
}

func resetStats(cards *[]Card) {
	for i := range *cards {
		(*cards)[i].missed = 0
	}

	logPrintln("Card statistics have been reset")
	return
}

func hardestCard(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].missed > cards[j].missed
	})

	noCardsMsg := "There are no cards with errors."
	if len(cards) == 0 || cards[0].missed == 0 {
		logPrintln(noCardsMsg)
		return
	}

	var terms []string
	for _, card := range cards {
		if card.missed != cards[0].missed {
			break
		}
		terms = append(terms, card.term)
	}

	if len(terms) == 1 {
		msg := fmt.Sprintf("The hardest card is \"%s\". You have %d errors answering it.", terms[0], cards[0].missed)
		logPrintln(msg)
		return
	}

	hardestMsg := fmt.Sprintf("The hardest cards are \"%s\"", terms[0])
	for _, term := range terms[1:] {
		hardestMsg += fmt.Sprintf(", \"%s\"", term)
	}

	logPrintln(hardestMsg)
	return
}

func log_() {
	fileName := readLine("File name:")
	file, _ := os.Create(fileName)
	defer file.Close()

	file.WriteString(sessionLog.String())
	logPrintln("The log has been saved")
	return
}

func ask(cards *[]Card) {
	if len(*cards) == 0 {
		logPrintln("No cards yet!")
	}
	count := readInt("How many times to ask?")
	quizCards(cards, count)
}

func quizCards(cards *[]Card, count int) {
	for i := 0; i < count; i++ {
		index := i % len(*cards)

		answerPrompt := fmt.Sprintf("Print the definition of \"%s\"", (*cards)[index].term)
		answer := readLine(answerPrompt)

		if answer == (*cards)[index].definition {
			logPrintln("Correct!")
			continue
		}

		(*cards)[index].missed++
		if definitionExists(answer, *cards) {
			val2 := getTermFor(answer, *cards)
			format := "Wrong. The right answer is \"%s\", but your definition is correct for \"%s\"."
			msg := fmt.Sprintf(format, (*cards)[index].definition, val2)
			logPrintln(msg)
		} else {
			msg := fmt.Sprintf("Wrong. The right answer is \"%s\".", (*cards)[index].definition)
			logPrintln(msg)
		}
	}
	return
}

func export(cards []Card, exportFile string) {
	if exportFile == defaultCLExport {
		exportFile = readLine("File name:")
	}
	file, _ := os.Create(exportFile)

	for _, card := range cards {
		fmt.Fprintln(file, card.term, card.definition, card.missed)
	}

	msg := fmt.Sprintf("%d cards have been saved.", len(cards))
	logPrintln(msg)
	return
}

func import_(cards *[]Card, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		logPrintln("File not found")
		return
	}
	defer file.Close()
	var loadedCards int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cardString := scanner.Text()
		addImportCard(cards, cardString)
		loadedCards++
	}

	msg := fmt.Sprintf("%d cards have been loaded.", loadedCards)
	logPrintln(msg)
	return
}

func addImportCard(cards *[]Card, cardText string) {
	var importCard Card
	fmt.Sscanf(cardText, "%s %s %d", &(importCard.term), &(importCard.definition), &(importCard.missed))

	if !termExists(importCard.term, *cards) {
		*cards = append(*cards, importCard)
		return
	}

	for i, card := range *cards {
		if importCard.term == card.term {
			(*cards)[i] = importCard
			break
		}
	}

	return
}

func remove(cards *[]Card) {
	cardToDelete := readLine("Which card?")
	var newCards []Card

	for _, card := range *cards {
		if card.term != cardToDelete {
			newCards = append(newCards, card)
		}
	}

	if len(*cards) == len(newCards) {
		msg := fmt.Sprintf("Can't remove \"%s\": there is no such card.", cardToDelete)
		logPrintln(msg)
	} else {
		logPrintln("The card has been removed.")
		*cards = newCards
	}

	return
}

func add(cards *[]Card) {
	newCard := Card{
		term:       getNewTerm(*cards),
		definition: getNewDefinition(*cards),
		missed:     0,
	}
	*cards = append(*cards, newCard)

	msg := fmt.Sprintf("The pair (\"%s\":\"%s\") has been added", newCard.term, newCard.definition)
	logPrintln(msg)
	return
}

func getNewTerm(cards []Card) (newTerm string) {
	newTerm = readLine("The card")
	for {
		if !termExists(newTerm, cards) {
			return newTerm
		}
		newPrompt := fmt.Sprintf("The term \"%s\" already exists. Try again:", newTerm)
		newTerm = readLine(newPrompt)
	}
}

func termExists(term string, cards []Card) (exists bool) {
	for _, card := range cards {
		if card.term == term {
			return true
		}
	}

	return false
}

func getNewDefinition(cards []Card) (newDefinition string) {
	newDefinition = readLine("The definition of the card")

	for {
		if !definitionExists(newDefinition, cards) {
			return newDefinition
		}

		newPrompt := fmt.Sprintf("The definition \"%s\" already exists. Try again:", newDefinition)
		newDefinition = readLine(newPrompt)
	}
}

func definitionExists(definition string, cards []Card) (exists bool) {
	for _, card := range cards {
		if card.definition == definition {
			return true
		}
	}

	return false
}

func getTermFor(definition string, cards []Card) (term string) {
	for _, card := range cards {
		if card.definition == definition {
			term = card.term
			break
		}
	}

	return term
}

func readLine(prompt string) (line string) {
	logPrintln(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	line = scanner.Text()
	sessionLog.WriteString(line)

	return line
}

func readInt(prompt string) (num int) {
	logPrintln(prompt)
	fmt.Scanln(&num)

	sessionLog.WriteString(strconv.Itoa(num))
	return num
}

func logPrintln(msg string) {
	fmt.Println(msg)
	sessionLog.WriteString(msg)
	return
}
