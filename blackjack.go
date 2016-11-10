package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Card struct {
	value  string // makes printing easier, and gives us a more elegant way to handle aces.
	suit   string
	hidden bool
}

func shuffle(deck []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 52; i++ {
		var i int = rand.Intn(52)
		var j int = rand.Intn(52)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func cardPrinter(hand []Card) {
	for range hand {
		fmt.Printf(" ---  ")
	}
	fmt.Println()
	for _, card := range hand {
		if card.hidden {
			fmt.Printf("|   | ")
		} else {
			fmt.Printf("|%s  | ", card.suit)
		}
	}
	fmt.Println()
	for range hand {
		fmt.Printf("|   | ")
	}
	fmt.Println()
	for _, card := range hand {
		if card.hidden {
			fmt.Printf("|   | ")
		} else {
			fmt.Printf("|  %s| ", card.value)
		}
	}
	fmt.Println()
	for range hand {
		fmt.Printf(" ---  ")
	}
	fmt.Println()
}

func parseCard(card Card) int {
	var ret int
	switch card.value {
	case "a":
		ret = 1
	case "j":
		ret = 10
	case "q":
		ret = 10
	case "k":
		ret = 10
	default:
		ret, _ = strconv.Atoi(card.value)
	}

	return ret
}

func cardMaker(card int) Card {
	var rank int = card % 12
	var ret Card
	switch {
	case card < 13:
		ret.suit = "c"
	case card < 26:
		ret.suit = "d"
	case card < 39:
		ret.suit = "h"
	case card <= 51:
		ret.suit = "s"
	}
	switch rank {
	case 0:
		ret.value = "a"
	case 9:
		ret.value = "j"
	case 10:
		ret.value = "q"
	case 11:
		ret.value = "k"
	default:
		ret.value = strconv.Itoa(rank + 1)
	}
	return ret
}

func totalDeck(deck []Card) int {
	var sum int
	var aces []int
	for _, card := range deck {
		if card.value == "a" {
			sum += 1
			aces = append(aces, 1)
		} else {
			sum += parseCard(card)
		}
	}
	for range aces {
		if sum+10 <= 21 {
			sum += 10
		}
	}
	return sum
}

func cls() {
	c := exec.Command("clear") // I've only tested this on macOS in Terminal.app
	c.Stdout = os.Stdout
	c.Run()
}

func drawHands(uhand []Card,dhand []Card) {
	cls()
	fmt.Println("\U0001f514  Bell Labs Casino")
	fmt.Println("My hand:")
	cardPrinter(dhand)
	fmt.Println("Your hand:")
	cardPrinter(uhand)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// some state vars
	var uhand []Card
	var dhand []Card

	fmt.Println("Welcome to the Bell Labs Casino!") // ironically I wrote this on Windows in a Microsoft text editor.
	fmt.Println("Shuffling the deck...just a moment.")

	var deck []Card = make([]Card, 52)
	for i := 0; i < 52; i++ {
		deck[i] = cardMaker(i)
	}
	deck = shuffle(deck)

	for i := 0; i < 2; i++ {
		uhand = append(uhand, deck[0])
		deck = deck[1:]
		dhand = append(dhand, deck[0])
		deck = deck[1:]
	}

	time.Sleep(4 * time.Second)


	if totalDeck(uhand) == 21 && totalDeck(dhand) == 21 {
		drawHands(uhand,dhand)
		fmt.Println("A natural tie!")
		os.Exit(0)
	} else if totalDeck(uhand) == 21 {
		drawHands(uhand,dhand)
		fmt.Println("You won with a natural!")
		os.Exit(0)
	} else if totalDeck(dhand) == 21 {
		drawHands(uhand,dhand)
		fmt.Println("I win with a natural!")
		os.Exit(0)
	}

	dhand[1].hidden = true // hide dealer's second card
	drawHands(uhand,dhand)
	fmt.Println("No naturals. Let's play!")
	
	time.Sleep(4 * time.Second)

	for totalDeck(uhand) < 21 {
		drawHands(uhand,dhand)

		fmt.Printf("Your total is %d. (H)it or (s)tand?\n", totalDeck(uhand))
		var resp string = ""
		fmt.Scanf("%s", &resp)
		if resp == "S" || resp == "s" {
			break
		} else if resp == "H" || resp == "h" {
			uhand = append(uhand, deck[0])
			deck = deck[1:]
		}

	}
	drawHands(uhand,dhand)

	var utotal int = totalDeck(uhand)

	if utotal > 21 {
		fmt.Println("You've busted! I win.")
	} else {
		fmt.Printf("Your total is %d.\n", utotal)
		fmt.Println("My turn to go.")
		dhand[1].hidden = false // unhide second card.

		if totalDeck(dhand) > 17 {
			dhand = append(dhand, deck[0])
			deck = deck[1:]
			drawHands(uhand,dhand)
		} else {
			for totalDeck(dhand) <= 17 {
				dhand = append(dhand, deck[0])
				deck = deck[1:]
				drawHands(uhand,dhand)
			}
		}

		if dtotal := totalDeck(dhand); dtotal > 21 {
			fmt.Println("I've busted! You win!")
		} else if dtotal < utotal {
			fmt.Println("You win!")
		} else if dtotal == utotal {
			fmt.Println("It's a tie!")
		} else {
			fmt.Println("I win!")
		}
p	}
	os.Exit(0)
}
