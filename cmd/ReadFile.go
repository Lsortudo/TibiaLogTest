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
			*playerDamageTaken += damage
		}
	}
	/*parts := strings.Split(message, " lose ")
	n, err := fmt.Sscanf(parts[1], "%d hitpoints ", &damage)
	if n == 1 && err == nil {
		// Imprimindo a quantidade de dano sofrido pelo jogador e adicionando-a à variável playerDamageTaken.
		fmt.Printf("Você sofreu %d de dano\n", damage)
		*playerDamageTaken += damage
	}*/
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

	// Itera sobre cada linha do arquivo
	for scanner.Scan() {
		message := scanner.Text()
		if strings.Contains(message, "You lose") {
			playerLossMessageProcessor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
		}
	}

	// Console messages

	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("Dano total que você sofreu: %d\n", playerDamageTaken)

}
