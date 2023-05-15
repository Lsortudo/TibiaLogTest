/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type DamageByCreature map[string]int

type LootMap map[string]int

type JSONData struct {
	HitpointsHealed int `json:"hitpointsHealed"`
	DamageTaken     struct {
		Total          int              `json:"total"`
		ByCreatureKind DamageByCreature `json:"byCreatureKind"`
	} `json:"damageTaken"`
	ExperienceGained    int     `json:"experienceGained"`
	Loot                LootMap `json:"loot"`
	HealthBlackKnight   int     `json:"healthBlackKnight"`
	UnknownOriginDamage int     `json:"unknownOriginDamage"`
}

// Declarar variaveis
var filePath string
var unknownDamageOrigin int
var healthBlackKnight int

// Declarar mapas
var enemyDamages = make(map[string]int)
var lootMap = make(map[string]int)

// Dataclass  { Use struct in go }
type ByDamageDesc []struct {
	Creature string
	Damage   int
}

func (a ByDamageDesc) Len() int           { return len(a) }
func (a ByDamageDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDamageDesc) Less(i, j int) bool { return a[i].Damage > a[j].Damage }

//var lootList []Loot

type Loot struct {
	item  string
	count int
}
type ByCount []Loot

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].count > a[j].count }

type PlayerLossMessageProcessor struct{}
type PlayerHealedMessageProcessor struct{}
type PlayerExperiencedMessageProcessor struct{}

// ReadFileCmd represents the ReadFile command
var ReadFileCmd = &cobra.Command{
	Use:   "ReadFile",
	Short: "A code to read logs from a game called Tibia",
	Long:  `This is a project where u can run the code go run main.go --path yourfilepath to read an txt file containing logs from a game called Tibia`,
	Run: func(cmd *cobra.Command, args []string) {
		ReadServerLogFile()
	},
}

func init() {
	rootCmd.AddCommand(ReadFileCmd)
	ReadFileCmd.PersistentFlags().StringVarP(&filePath, "path", "p", "", "path to log file")
	ReadFileCmd.MarkPersistentFlagRequired("path")
}

type InterfaceMessageProcessor interface { // Interface to implement
	Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int)
}

func (p *PlayerLossMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	var enemyName string
	parts := strings.Split(message, " ") // Split the line into words
	for i, word := range parts {         // Search for the word that indicates my damage value
		if word == "lose" && i < len(parts)-1 {
			damageStr := parts[i+1]                // I'm gonna get the value after i+1 where i = "lose"
			damage, err := strconv.Atoi(damageStr) // Convert to integer
			if err != nil {
				fmt.Println("Erro ao converter o valor do dano:", err) //
				continue
			}
			enemyName = parts[len(parts)-2]                          // Get the creature name by getting the total length of parts which is all the words in the line, and going backwards
			enemyName = strings.TrimRight(enemyName, ".")            // Remove the "." on the end of the creature name
			if enemyName == "hitpoint" || enemyName == "hitpoints" { // an if to get the unknown source damage and store it to show it later, since 'unknown' is not a creature
				unknownDamageOrigin += damage
				continue
			}
			*playerDamageTaken += damage      // Define the pointer with += damage
			enemyDamages[enemyName] += damage // Map giving each creature name + their damage
		}
	}
}
func (h *PlayerHealedMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	parts := strings.Split(message, " ") // Split the line into words
	for i, word := range parts {         // Search for the word that indicates my heal value
		if word == "for" && i < len(parts)-1 {
			healStr := parts[i+1]              // I'm gonna get the value after i+1 where i = "for"
			heal, err := strconv.Atoi(healStr) // Convert to integer
			if err != nil {
				fmt.Println("Erro ao converter o valor da cura:", err)
				continue
			}
			*playerHealed += heal // Define the pointer with += heal
		}
	}
}
func (e *PlayerExperiencedMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	parts := strings.Split(message, " ") // Split the line into words
	for i, word := range parts {         // Search for the word that indicates my experience gained value
		if word == "gained" && i < len(parts)-1 {
			experienceStr := parts[i+1]                    // I'm gonna get the value after i+1 where i = "gained"
			experience, err := strconv.Atoi(experienceStr) // Convert to integer
			if err != nil {
				fmt.Println("Erro ao converter o valor de experiência:", err)
				continue
			}

			*playerExperience += experience // Define the pointer with += experience
		}
	}

}
func creatureBlackKnightHealth(message string) {
	damageStr := ""                          // Initializes an empty string to store dmg from the message
	parts := strings.Split(message, "loses") // Splits the message using "loses" as separator

	if len(parts) >= 2 { // Check if theres two elements, cuz my message is like "loses 15 hitpoints" so in this if i'll get "15" and "hitpoints"
		parts = strings.Split(parts[1], "hitpoints")
		damageStr = strings.TrimSpace(parts[0])
	}

	damage, err := strconv.Atoi(damageStr) // Convert to integer
	if err != nil {
		fmt.Println("Erro ao converter o valor do dano:", err)
	}

	healthBlackKnight += damage // Creature life is the sum of all damage before dies

}

func creatureLootTotal(message string) {
	lootText := strings.Split(message, ": ")[1]     // Split message string using ":" as the separator and selects the second part (index[1])
	lootText = strings.TrimRight(lootText, ".,\"'") // Removing ".,'" at the end of the string
	items := strings.Split(lootText, ", ")          // Splits loot text into multiple items using "," as the separator
	for _, item := range items {                    // Use '_' so i don't have to specify index value on the loop
		itemParts := strings.Fields(item) // Splits the item string into fields(words)
		count := 1                        // Default value for items such as: a gold coin, diamond sword, plate legs
		if len(itemParts) > 1 {           // Check if theres more than one element, meaning there's already a quantity like "43 gold coins" or there's a specified name like "a gold coin"
			if itemParts[0] == "a" || itemParts[0] == "an" { // Check if an item has a/an before the name such as: a gold coin
				itemParts = itemParts[1:] // Remove the a/an from it, soo it can store only the item name
			}
			quantity, err := strconv.Atoi(itemParts[0]) // Convert to integer
			if err == nil {
				count = quantity          // Get the quantity item so i can send to my map itemParts and count as key/value
				itemParts = itemParts[1:] // If successfull remove the quantity from item parts cuz i'll use itemParts as the name and quantity is already on my count
			}
		} else {
			count = 1 // If there's just one element, set it to 1, // I can remove this else later on
		}
		itemName := strings.Join(itemParts, " ") // itemName will be equal as the name from itemParts without "a/an" or quantity
		switch itemName {                        // Switch case to check if there's a nothing item, if yes, it won't add to the map to be showed, removing this u can also track how much 'nothing' has been dropped
		case "nothing":
		default:
			lootMap[itemName] += count
		}
	}
}

func getSingularItem(item string) string {
	if strings.HasSuffix(item, "s") { // Verify if 'item' ends with the letter "s"
		return item[:len(item)-1] // If yes uses slicing to remove the "s" and return the substring of 'item' excludind the last character
	}
	return item // Return the original item
}

func ReadServerLogFile() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) // Creating a scanner to read the log file

	// Starting all vars thats gonna be used
	var playerDamageTaken, playerHealed, playerExperience int
	var playerLossMessageProcessor InterfaceMessageProcessor = &PlayerLossMessageProcessor{}
	var playerHealedMessageProcessor InterfaceMessageProcessor = &PlayerHealedMessageProcessor{}
	var playerExperiencedMessageProcessor InterfaceMessageProcessor = &PlayerExperiencedMessageProcessor{}

	for scanner.Scan() { // Starts a loop to read each line of input
		message := scanner.Text() // Message represents a string of each line on the file
		switch {
		case strings.Contains(message, "You lose"): // Check if contains You lose, if so it calls Process method passing Message as an argument as well as pointers and variables as arguments
			playerLossMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		case strings.Contains(message, "You healed"): // Check if contains You healed, if so it calls Process method passing Message as an argument as well as pointers and variables as arguments
			playerHealedMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		case strings.Contains(message, "You gained"): // Check if contains You gained, if so it calls Process method passing Message as an argument as well as pointers and variables as arguments
			playerExperiencedMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		case strings.Contains(message, "Black Knight"): // Check if contains Black Knight, if so it calls a function and pass Message as an argument
			creatureBlackKnightHealth(message)
		case strings.Contains(message, "Loot of"): // Check if contains Loot of, if so it calls a function and pass Message as an argument
			creatureLootTotal(message)
		}
	}
	combinedLootMap := make(map[string]int) // Starting a new map that's gonna be read
	for item, count := range lootMap {
		singularItem := getSingularItem(item)  // Calling a function to get singular items cuz there's items like: a gold coin | 15 gold coins
		combinedLootMap[singularItem] += count // I'll add to the combined map everything correctly
	}
	var combinedLootList []Loot // New variable as slice of Loot struct
	for item, count := range combinedLootMap {
		combinedLootList = append(combinedLootList, Loot{item, count}) // Appends a new loot struct to combinedLootList slice and i can access items from the map with loot.count and loot.item
	}
	sort.Sort(ByCount(combinedLootList)) // Sort in descending order

	// Console messages
	fmt.Printf("Total healed: %d\n", playerHealed)
	fmt.Printf("Total damage suffered: %d\n", playerDamageTaken+unknownDamageOrigin) // Adding the total damage from creatures + unknown origin as requested on additional notes
	// Sorting enemyDamages map
	sortedDamages := make([]struct {
		Creature string
		Damage   int
	}, 0, len(enemyDamages))

	for creature, damage := range enemyDamages {
		sortedDamages = append(sortedDamages, struct {
			Creature string
			Damage   int
		}{creature, damage})
	}

	sort.Sort(ByDamageDesc(sortedDamages))

	for _, entry := range sortedDamages {
		//fmt.Printf("Creature: %s, Damage: %d\n", strings.Title(entry.Creature), entry.Damage)
		fmt.Printf("The creature %s dealt %d total damage\n", strings.Title(entry.Creature), entry.Damage)
	}
	fmt.Printf("Experience gained: %d\n", playerExperience)

	for _, loot := range combinedLootList {
		fmt.Printf("%d %s\n", loot.count, loot.item)
	}
	fmt.Printf("------------------------ EXTRAS ------------------------\n")
	fmt.Printf("Total damage from unknown sources: %d\n", unknownDamageOrigin)
	fmt.Printf("Black Knight: %d\n", healthBlackKnight)
	fmt.Printf("------------------------ JSON ------------------------\n")
	// Convert data to JSON

	jsonData := JSONData{
		HitpointsHealed: playerHealed,
		DamageTaken: struct {
			Total          int              `json:"total"`
			ByCreatureKind DamageByCreature `json:"byCreatureKind"`
		}{
			Total:          playerDamageTaken + unknownDamageOrigin,
			ByCreatureKind: enemyDamages,
		},
		ExperienceGained:    playerExperience,
		Loot:                combinedLootMap,
		UnknownOriginDamage: unknownDamageOrigin,
		HealthBlackKnight:   healthBlackKnight,
	}
	// Printar JSON no console, utilizar JSON Beauty Website pra verificar se tá na estrutura solicitada
	/*jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal("Erro ao converter para JSON:", err)
	}

	fmt.Println(string(jsonBytes))*/
	jsonExternalOutput, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Fatal("Erro ao converter para JSON:", err)
	}
	filePathJson := "output.json"
	err = ioutil.WriteFile(filePathJson, jsonExternalOutput, 0644)
	if err != nil {
		log.Fatal("Erro ao gravar arquivo:", err)
	}
	fmt.Println("JSON gravado com sucesso no arquivo:", filePathJson)

}
