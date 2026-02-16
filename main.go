package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var input string

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input = scanner.Text()
		cInput := cleanInput(input)
		fmt.Printf("Your command was: %s\n", cInput[0])
		input = ""
	}
}
