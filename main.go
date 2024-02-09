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
	Blue   = "\033[34m"
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red = "\033[31m"
	Cyan = "\033[36m"
)

var sessionLog strings.Builder

func main() {
	var cards []Card
	importFile := flag.String("import_from", defaultCLImport, "Enter a file name to import cards")
	exportFile := flag.String("export_to", defaultCLExport, "Enter a file name to export cards to on exit")
	flag.Parse()
	checkForCLImport(&cards, *importFile)

	for {
		switch action := readLine("\n" + Blue + menuPrompt + Reset); action {
		case "add":
			add(&cards)
		case "remove":
			remove(&cards)
		case "import":
			import_(&cards, readLine(Blue + "File to import cards from:" + Reset))
		case "export":
			export(cards, *exportFile)
		case "ask":
			ask(&cards)
		case "exit":
			logPrintln(Blue + "Bye bye!")
			export(cards, *exportFile)
			return
		case "log":
			log_()
		case "hardest card":
			hardestCard(cards)
		case "reset stats":
			resetStats(&cards)
		default:
			logPrintln(Red + `Command "` + action + `" not recognized. Please enter another command.` + Reset)
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

	logPrintln(Green + "Card statistics have been reset" + Reset)
	return
}

func hardestCard(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].missed > cards[j].missed
	})

	noCardsMsg := Green + "There are no cards with errors." + Reset
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
		logPrintln(Green + `The hardest card is "` + terms[0] + `". You have ` + strconv.Itoa(cards[0].missed) + ` errors answering it.` + Reset)
		return
	}

	hardestMsg := fmt.Sprintf("The hardest cards are \"%s\"", terms[0])
	for _, term := range terms[1:] {
		hardestMsg += fmt.Sprintf(", \"%s\"", term)
	}

	logPrintln(Green + hardestMsg + Reset)
	return
}

func log_() {
	fileName := readLine(Blue + "File name:" + Reset)
	file, _ := os.Create(fileName)
	defer file.Close()

	file.WriteString(sessionLog.String())
	logPrintln(Green + "The log has been saved." + Reset)
	return
}

func ask(cards *[]Card) {
	if len(*cards) == 0 {
		logPrintln(Red + "No cards yet!" + Reset)
		return
	}
	count := readInt(Blue + "How many cards?" + Reset)
	quizCards(cards, count)
	return
}

func quizCards(cards *[]Card, count int) {
	for i := 0; i < count; i++ {
		index := i % len(*cards)

		answer := readLine(Blue + "Card: \"" + Cyan + (*cards)[index].term + Reset + Blue + "\"\nAnswer:" + Reset)
		if answer == (*cards)[index].definition {
			logPrintln(Green + "Correct!" + Reset)
			continue
		}

		(*cards)[index].missed++
		if definitionExists(answer, *cards) {
			val2 := getTermFor(answer, *cards)
			msg := Red + "Wrong. The right answer is \"" + (*cards)[index].definition + "\", but your definition is correct for \"" + val2 + "\"." + Reset
			logPrintln(msg)
		} else {
			logPrintln(Red + "Wrong. The right answer is \"" + (*cards)[index].definition + "\"." + Reset)
		}
	}
	return
}

func export(cards []Card, exportFile string) {
	if exportFile == defaultCLExport {
		exportFile = readLine(Blue + "\nFile path to export cards to:" + Reset)
	}
	file, _ := os.Create(exportFile)

	for _, card := range cards {
		fmt.Fprintln(file, card.term, card.definition, card.missed)
	}


	logPrintln(Green + strconv.Itoa(len(cards)) + " cards have been saved." + Reset)
	return
}

func import_(cards *[]Card, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		logPrintln(Red + "File not found" + Reset)
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

	logPrintln(Green + strconv.Itoa(loadedCards) + " cards have been loaded." + Reset)
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
	cardToDelete := readLine(Blue + "Which card? (enter the front term/sentence of the card)" + Reset)
	var newCards []Card

	for _, card := range *cards {
		if card.term != cardToDelete {
			newCards = append(newCards, card)
		}
	}

	if len(*cards) == len(newCards) {
		logPrintln(Red + "Can't remove \"" + cardToDelete + "\": there is no such card." + Reset)
	} else {
		logPrintln(Green + "The card has been removed." + Reset)
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

	logPrintln(Green + `The pair ("` + newCard.term + `":"` + newCard.definition + `") has been added` + Reset)
	return
}

func getNewTerm(cards []Card) (newTerm string) {
	newTerm = readLine(Blue + "\nNew card front:" + Reset)
	for {
		if !termExists(newTerm, cards) {
			return newTerm
		}
		newTerm = readLine(Red + `The term "` + newTerm + `" already exists. Try a different term:` + Reset)
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
	newDefinition = readLine(Blue + "New card back:" + Reset)

	for {
		if !definitionExists(newDefinition, cards) {
			return newDefinition
		}

		newDefinition = readLine(Red + `The definition "` + newDefinition + `" already exists. Try a different definition:` + Reset)
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
