
# Technical test

## Introduction
Develop an application to process and extract data from Tibia server log.
In this log, a new line is created whenever the player executes some action, or is involved in some world action.

## Tasks
 ### Log Parser
 The application must:
**Read the log file
Extract the following informations from the log:**
- Total damage healed by the player;
- Total damage taken by the player;
- Total damage taken by the player per creature kind;
- Total experience received;
- Amount of items dropped by each creature;


 **_Note: Total damage must consider damage taken by unknown origins;_**
 **_Note: Unknown origins should not appear as a creature kind on damage taken._**
 _Extra: What is the total health of Black Knight creature?_
 _Extra: What is the total damage taken by unknown origins?_



## Response

This is the answer i got when reading the TXT that was given in the Technical Test email.
(Expected answer may change slightly because I am sorting descending based on quantity)
Expected response that i got:
```

Total healed: 8048
Total damage suffered: 7683
The creature Dragon dealt 2427 total damage 
The creature Lord dealt 1770 total damage   
The creature Scorpion dealt 470 total damage
The creature Knight dealt 449 total damage  
The creature Bonelord dealt 234 total damage
The creature Cyclops dealt 101 total damage 
The creature Wyvern dealt 47 total damage   
The creature Smith dealt 29 total damage    
The creature Soldier dealt 17 total damage  
The creature Ghoul dealt 2 total damage     
Experience gained: 31363
1177 gold coin
26 dragon ham
7 bolt
6 meat
4 steel shield
3 scorpion tail
3 dragon's tail
3 plate leg
2 crossbow
2 white mushroom
2 burst arrow
2 small diamond
2 longsword
2 steel helmet
2 hatchet
2 soldier helmet
1 bone
1 ham
1 green dragon leather
1 cyclops trophy
1 leather leg
1 cyclops toe
1 royal spear
1 spellbook
1 pick
1 letter
1 two handed sword
1 copper shield
1 strong health potion
------------------------ EXTRAS ------------------------
Total damage from unknown sources: 2137
Black Knight: 1800

```

## How to run
I used Cobra (CLI) sooo to run the command file called `ReadFile` u can use these 2 options:

    go run main.go ReadFile --path yourPath\ServerLog.txt
    
or

    ./main ReadFile -p yourPath\ServerLog.txt

Use `./main` to see available commands
Use `./main ReadFile` to see the flags

Use `go test -v ./cmd/` to execute some basic unit tests, there u can see errors trying to pass 'abc' as values

# Cobra

**Cobra** allows you to define commands and subcommands, as well as custom flags and arguments, to create a hierarchy of commands that can be executed by the user; link to [Github](https://github.com/spf13/cobra/).



