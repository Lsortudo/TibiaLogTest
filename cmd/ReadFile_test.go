/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "testing"

func TestCreatureBlackKnightHealth(t *testing.T) {
	healthBlackKnight = 0

	// Teste caso a mensagem seja válida
	message := `18:46 A Black Knight loses 200 hitpoints due to your attack.  `

	creatureBlackKnightHealth(message)
	expectedHealth := 200
	if healthBlackKnight != expectedHealth {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHealth, healthBlackKnight)
	}
	message = "18:46 A Black Knight loses 40 hitpoints due to your attack."
	creatureBlackKnightHealth(message)
	expectedHealth += 40
	if healthBlackKnight != expectedHealth {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHealth, healthBlackKnight)
	}

	// Teste caso a mensagem não tenha a palavra-chave "hitpoints"
	message = "18:46 A Black Knight loses 15 due to your attack."
	creatureBlackKnightHealth(message)
	if healthBlackKnight != expectedHealth {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHealth, healthBlackKnight)
	}

	// Teste caso o valor do dano não possa ser convertido para inteiro
	message = "18:46 A Black Knight loses abc hitpoints due to your attack."
	creatureBlackKnightHealth(message)
	if healthBlackKnight != expectedHealth {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHealth, healthBlackKnight)
	}

}

func TestCreatureLootTotal(t *testing.T) {
	// Caso de teste com uma mensagem válida contendo um item
	message := "15:43 Loot of a cyclops: a sword"
	lootMap = make(map[string]int)
	creatureLootTotal(message)
	expectedLoot := map[string]int{"sword": 1}
	if !equalLootMap(lootMap, expectedLoot) {
		t.Errorf("Resultado incorreto. Esperado: %v, Obtido: %v", expectedLoot, lootMap)
	}

	// Caso de teste com uma mensagem válida contendo múltiplos itens
	message = "15:43 Loot of a cyclops: 3 gold coins, an apple, a potion"
	lootMap = make(map[string]int)
	creatureLootTotal(message)
	expectedLoot = map[string]int{"gold coins": 3, "apple": 1, "potion": 1}
	if !equalLootMap(lootMap, expectedLoot) {
		t.Errorf("Resultado incorreto. Esperado: %v, Obtido: %v", expectedLoot, lootMap)
	}

	// Caso de teste com uma mensagem válida contendo um item com quantidade implícita
	message = "18:31 Loot of a wyvern: dragon ham, 21 gold coins."
	lootMap = make(map[string]int)
	creatureLootTotal(message)
	expectedLoot = map[string]int{"dragon ham": 1, "gold coins": 21}
	if !equalLootMap(lootMap, expectedLoot) {
		t.Errorf("Resultado incorreto. Esperado: %v, Obtido: %v", expectedLoot, lootMap)
	}

	// Caso de teste com uma mensagem válida contendo o item "nothing"
	message = "18:31 Loot of a wyvern: nothing."
	lootMap = make(map[string]int)
	creatureLootTotal(message)
	expectedLoot = map[string]int{}
	if !equalLootMap(lootMap, expectedLoot) {
		t.Errorf("Resultado incorreto. Esperado: %v, Obtido: %v", expectedLoot, lootMap)
	}
}

func TestPlayerExperiencedMessageProcessor_Process(t *testing.T) {
	// Caso de teste com uma mensagem válida contendo ganho de experiência
	message := "15:44 You gained 100 experience points.	"
	playerHealed := 0
	playerDamageTaken := 0
	playerExperience := 0
	processor := PlayerExperiencedMessageProcessor{}
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedExperience := 100
	if playerExperience != expectedExperience {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedExperience, playerExperience)
	}

	// Caso de teste com uma mensagem válida contendo múltiplos ganhos de experiência
	message = "15:44 You gained 50 experience points.15:44 You gained 25 experience points."
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedExperience = 75
	if playerExperience != expectedExperience {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedExperience, playerExperience)
	}

	// Caso de teste com uma mensagem inválida sem ganho de experiência
	message = "15:43 You lose 31 hitpoints due to an attack by a cyclops. "
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedExperience = 0
	if playerExperience != expectedExperience {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedExperience, playerExperience)
	}
}

func TestPlayerHealedMessageProcessor_Process(t *testing.T) {
	// Caso de teste com uma mensagem válida contendo valor de cura
	message := "15:42 You healed yourself for 50 hitpoints."
	playerHealed := 0
	playerDamageTaken := 0
	playerExperience := 0
	processor := PlayerHealedMessageProcessor{}
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedHeal := 50
	if playerHealed != expectedHeal {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHeal, playerHealed)
	}

	// Caso de teste com uma mensagem válida contendo múltiplos valores de cura
	message = "15:42 You healed yourself for 30 hitpoints. 15:42 You healed yourself for 20 hitpoints."
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedHeal = 50
	if playerHealed != expectedHeal {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHeal, playerHealed)
	}

	// Caso de teste com uma mensagem inválida sem valor de cura
	message = "15:42 You gained 100 experience points."
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedHeal = 0
	if playerHealed != expectedHeal {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedHeal, playerHealed)
	}
}

func TestPlayerLossMessageProcessor_Process(t *testing.T) {
	// Caso de teste com uma mensagem válida contendo dano sofrido
	message := "15:43 You lose 50 hitpoints due to an attack by a cyclops."
	playerHealed := 0
	playerDamageTaken := 0
	playerExperience := 0
	processor := PlayerLossMessageProcessor{}
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedDamage := 50
	if playerDamageTaken != expectedDamage {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedDamage, playerDamageTaken)
	}

	// Caso de teste com uma mensagem válida contendo múltiplos danos sofridos
	message = "15:43 You lose 17 hitpoints due to an attack by a cyclops. 15:43 You lose 31 hitpoints due to an attack by a cyclops."
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedDamage = 48
	if playerDamageTaken != expectedDamage {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedDamage, playerDamageTaken)
	}

	// Caso de teste com uma mensagem inválida sem dano sofrido
	message = "15:44 You gained 150 experience points."
	playerHealed = 0
	playerDamageTaken = 0
	playerExperience = 0
	processor.Process(message, &playerHealed, &playerDamageTaken, &playerExperience)
	expectedDamage = 0
	if playerDamageTaken != expectedDamage {
		t.Errorf("Resultado incorreto. Esperado: %d, Obtido: %d", expectedDamage, playerDamageTaken)
	}
}

func equalLootMap(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for key, val := range a {
		if bval, ok := b[key]; !ok || bval != val {
			return false
		}
	}
	return true
}
