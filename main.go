package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CoffeeUnit struct {
	Amount int
	Price  int
}
type CoffeeShopData map[string]CoffeeUnit

func IsFileEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	return info.Size() == 0, err
}

func MakeCoffeeData(raw_data []byte, isEmpty bool) (CoffeeShopData, error) {
	coffeeMap := make(CoffeeShopData)
	if isEmpty {
		return coffeeMap, nil
	}
	err := json.Unmarshal(raw_data, &coffeeMap)
	return coffeeMap, err
}

func SaveCoffeDataToJSONFile(csd CoffeeShopData, path string) error {
	jsonString, err := json.Marshal(csd)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, json_err := file.Write(jsonString)
	if json_err != nil {
		return json_err
	}
	return nil
}
func main() {
	dataPath := "./coffee.json"
	fileContent, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	isEmpty, err := IsFileEmpty(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	coffeeData, err := MakeCoffeeData(fileContent, isEmpty)
	if err != nil {
		log.Fatal(err)
	}

	isRunning := true

	for {
		if !isRunning {
			break
		}
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		arguments := strings.Fields(line)
		command := arguments[0]
		switch command {
		case "exit":
			fmt.Println("Bye!")
			isRunning = false
		case "items":
			for k := range coffeeData {
				fmt.Println(k)
			}

		case "show":
			for k, contents := range coffeeData {
				fmt.Println(k)
				fmt.Printf("    Price: %d\n", contents.Price)
				fmt.Printf("    Amount: %d\n", contents.Amount)
			}

		case "set":
			if len(arguments) < 4 {
				fmt.Println("Not enough arguments: set <Coffee Type> <Field> <Amount>.")
				continue
			}
			name := arguments[1]
			field := arguments[2]
			value, err := strconv.Atoi(arguments[3])
			if err != nil {
				fmt.Println("<Value> must be a numerical value.")
				continue
			}

			if strings.ToLower(field) == "price" {
				if entry, ok := coffeeData[name]; ok {
					entry.Price = value
					coffeeData[name] = entry
				}
			}
			if strings.ToLower(field) == "amount" {
				if entry, ok := coffeeData[name]; ok {
					entry.Amount = value
					coffeeData[name] = entry
				}
			}

		case "add":
			if len(arguments) < 2 {
				fmt.Println("Not enough arguments: add <CoffeeType>")
				continue
			}
			name := arguments[1]
			coffeeNewEntry := CoffeeUnit{0, 0}
			_, ok := coffeeData[name]
			if ok {
				fmt.Printf("%s already added\n", name)
				continue
			}
			coffeeData[name] = coffeeNewEntry
		default:
			fmt.Printf("Command %s does not exist.\n", command)
		}
		err = SaveCoffeDataToJSONFile(coffeeData, dataPath)
		if err != nil {
			log.Fatal(err)
		}
	}

}
