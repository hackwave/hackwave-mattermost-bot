package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

const (
	PLUGIN_NAME         = "DiceRoller"
	PLUGIN_REGEX        = `^roll (|\d{1,3})(d|D|)\d{1,3}(|\+\d{1,3})$`
	PLUGIN_HELP_COMMAND = "   **roll**"
	PLUGIN_HELP_TEXT    = "   *Example:* `roll 2d20` *(or)* `roll 100` *(or)* `roll 3d4+2`"
)

func RollDice(numberOfSides, diceSides int) []int {
	var results []int
	rand.Seed(time.Now().UTC().UnixNano())
	for count := 0; count < numberOfSides; count++ {
		results = append(results, RollDie(diceSides))
	}
	return results
}

func RollDie(diceSides int) int {
	return rand.Intn(diceSides)
}

func ParseDiceCommand(message string) string {
	regex := regexp.MustCompile(" \\d{1,3}")
	numberOfDiceMatch := regex.FindStringSubmatch(message)
	var numberOfDice string
	if len(numberOfDiceMatch) > 0 {
		numberOfDice = numberOfDiceMatch[0][1:len(numberOfDiceMatch[0])]
	}

	regex = regexp.MustCompile("d\\d{1,3}")
	diceMatch := regex.FindStringSubmatch(message)
	var dice string
	if len(diceMatch) > 0 {
		dice = diceMatch[0][1:len(diceMatch[0])]
	}

	// Handle roll # without any d or number of die
	numberOfDiceInt, _ := strconv.Atoi(numberOfDice)
	diceInt, _ := strconv.Atoi(dice)
	if diceInt == 0 {
		dice = numberOfDice
		diceInt = numberOfDiceInt
		numberOfDiceInt = 1
		numberOfDice = "1"
	}

	results := RollDice(numberOfDiceInt, diceInt)
	var resultString string
	var resultSum int
	for count, result := range results {
		resultSum += result
		if len(results) == (count + 1) {
			resultString += fmt.Sprintf("%v ", result)
		} else {
			resultString += fmt.Sprintf("%v + ", result)
		}
	}

	regex = regexp.MustCompile("\\+\\d{1,3}")
	plusMatch := regex.FindStringSubmatch(message)

	if len(plusMatch) > 0 {
		// Modifier
		plus := plusMatch[0][1:len(plusMatch[0])]
		resultString = resultString + "( + " + plus + " ) "
		dice = dice + "(+" + plus + ")"
		plusInt, _ := strconv.Atoi(plus)
		resultSum += plusInt

	}
	return " **Rolling " + numberOfDice + "d" + dice + "**: " + resultString + " **= " + strconv.Itoa(resultSum) + "**"

}
