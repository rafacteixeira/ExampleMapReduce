package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	f, err := os.Open("fifa data editado.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	lists := make(chan []PlayerInfo)
	finalValue := make(chan []PlayerInfo)
	var wg sync.WaitGroup

	wg.Add(len(lines))

	for _, line := range lines {
		go func(player []string) {
			defer wg.Done()
			lists <- Map(player)
		}(line)
	}

	go Reducer(lists, finalValue)
	wg.Wait()
	close(lists)

	fmt.Println(<- finalValue)
}

func Map(player []string) []PlayerInfo {
	var list []PlayerInfo
	age, _ := strconv.Atoi(player[2])
	list = append(list, PlayerInfo{
		ID:   player[0],
		Name: player[1],
		Age:  age,
	})
	return list
}

func Reducer(mapList chan []PlayerInfo, sendFinalValue chan []PlayerInfo) {
	var final []PlayerInfo
	for list := range mapList {
		for _, value := range list {
			if value.Age <= 16 {
				final = append(final, value)
			}
		}
	}
	sendFinalValue <- final
}

type PlayerInfo struct {
	ID   string
	Name string
	Age  int
}
