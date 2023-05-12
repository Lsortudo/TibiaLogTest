/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Declarar variaveis
var filePath string
var unknownDamageOrigin int
var healthBlackKnight int

// Declarar mapas
var enemyDamages = make(map[string]int)

// Dataclass  { Use struct in go }

// ReadFileCmd represents the ReadFile command
var ReadFileCmd = &cobra.Command{
	Use:   "ReadFile",
	Short: "Descrição curta",
	Long:  `Descrição longa.`,
	Run: func(cmd *cobra.Command, args []string) {
		ReadServerLogFile()
	},
}

func init() {
	rootCmd.AddCommand(ReadFileCmd)
	ReadFileCmd.PersistentFlags().StringVarP(&filePath, "path", "p", "", "path to log file")
	ReadFileCmd.MarkPersistentFlagRequired("path")
}

type InterfaceMessageProcessor interface {
	Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int)
}

type PlayerLossMessageProcessor struct{}

func (p *PlayerLossMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	var enemyName string
	parts := strings.Split(message, " ")
	for i, word := range parts {
		if word == "lose" && i < len(parts)-1 {
			// Extrai o valor numérico da palavra seguinte
			damageStr := parts[i+1]
			damage, err := strconv.Atoi(damageStr)
			if err != nil {
				fmt.Println("Erro ao converter o valor do dano:", err)
				continue
			}
			enemyName = parts[len(parts)-2]
			enemyName = strings.TrimRight(enemyName, ".")
			if enemyName == "hitpoint" {
				unknownDamageOrigin += damage
				continue
			}
			*playerDamageTaken += damage
			enemyDamages[enemyName] += damage
		}
	}
}

type PlayerHealedMessageProcessor struct{}

func (h *PlayerHealedMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	// Divide a linha em palavras
	parts := strings.Split(message, " ")
	// Procura a palavra que indica o valor da cura
	for i, word := range parts {
		if word == "for" && i < len(parts)-1 {
			// Extrai o valor numérico da palavra seguinte
			healStr := parts[i+1]
			heal, err := strconv.Atoi(healStr)
			if err != nil {
				fmt.Println("Erro ao converter o valor da cura:", err)
				continue
			}
			*playerHealed += heal
		}
	}
}

type PlayerExperiencedMessageProcessor struct{}

func (e *PlayerExperiencedMessageProcessor) Process(message string, playerHealed *int, playerDamageTaken *int, playerExperience *int) {
	parts := strings.Split(message, " ")

	for i, word := range parts {
		if word == "gained" && i < len(parts)-1 {
			// Extrai o valor numérico da palavra seguinte
			experienceStr := parts[i+1]
			experience, err := strconv.Atoi(experienceStr)
			if err != nil {
				fmt.Println("Erro ao converter o valor de experiência:", err)
				continue
			}

			*playerExperience += experience
		}
	}

}

func creatureBlackKnightHealth(message string) {
	damageStr := ""
	parts := strings.Split(message, "loses")

	if len(parts) >= 2 {
		parts = strings.Split(parts[1], "hitpoints")
		damageStr = strings.TrimSpace(parts[0])
	}

	damage, err := strconv.Atoi(damageStr)
	if err != nil {
		fmt.Println("Erro ao converter o valor do dano:", err)
	}

	healthBlackKnight += damage

}

func ReadServerLogFile() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Criando um scanner para ler o arquivo de log do jogo.
	scanner := bufio.NewScanner(f)

	// Inicializando as variáveis que serão usadas no app
	var playerDamageTaken, playerHealed, playerExperience int

	var playerLossMessageProcessor InterfaceMessageProcessor = &PlayerLossMessageProcessor{}
	var playerHealedMessageProcessor InterfaceMessageProcessor = &PlayerHealedMessageProcessor{}
	var playerExperiencedMessageProcessor InterfaceMessageProcessor = &PlayerExperiencedMessageProcessor{}
	// Itera sobre cada linha do arquivo
	for scanner.Scan() {
		message := scanner.Text()
		if strings.Contains(message, "You lose") {
			playerLossMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		}
		if strings.Contains(message, "You healed") {
			playerHealedMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		}
		if strings.Contains(message, "You gained") {
			playerExperiencedMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		}
		if strings.Contains(message, "Black Knight") {
			creatureBlackKnightHealth(message)
		}
		//if strings.Contains(message, "Black Knight" ) && strings.Contains(message, "loses")
	}

	// Console messages

	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Dano total que você sofreu: %d\n", playerDamageTaken)
	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Total de cura: %d\n", playerHealed)
	fmt.Printf("----------------------------------------------------\n")
	for enemy, damage := range enemyDamages {
		fmt.Printf("O monstro %s lhe causou %d de dano total\n", strings.Title(enemy), damage)
	}
	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Total de dano desconhecido: %d\n", unknownDamageOrigin)
	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Total de experiência obtida: %d\n", playerExperience)
	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Vida de Black Knight: %d\n", healthBlackKnight)
}
