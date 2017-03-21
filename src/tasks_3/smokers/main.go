package main

import (
	"math/rand"
	"fmt"
	"time"
)

const PAPER = 0;
const FIRE = 1;
const TABACO = 2;


func getCigaretteStuff() (int, int) {
	stuff1 := rand.Intn(3)
	stuff2 := rand.Intn(3);
	for stuff2 == stuff1 {
		stuff2 = rand.Intn(3);
	}

	return stuff1, stuff2;
}

type table_t [3]bool;


func smoker(smokersChannels *[]chan bool, cigaretteSmokeChannel chan bool,  stuffId int, table *table_t) {
	fmt.Printf("Smoker%d: start\n", stuffId);
	for true {
		if <-(*smokersChannels)[stuffId] {
			fmt.Printf("Smoker%d: wake\n", stuffId);
			expectedStuff := !(*table)[stuffId];

			fmt.Printf("Smoker%d: exptected stuff: %v\n", stuffId, expectedStuff);
			if expectedStuff {
				clearTable(table);
				fmt.Printf("Smoker%d: smoking\n", stuffId);
				time.Sleep(time.Millisecond * 1000);
			}
			fmt.Printf("Smoker%d: result: %v\n", stuffId, expectedStuff);
			cigaretteSmokeChannel <- expectedStuff;
		}

	}
}


func dealer(smokersChannels *[]chan bool, cigaretteSmokeChannel chan bool, table *table_t) {
	fmt.Println("Dealer: start")
	for true {
		a, b := getCigaretteStuff();
		fmt.Printf("Dealer get %d and %d\n", a, b);
		putItems(table, a, b);
		fmt.Println("Sending event to smokers");
		for i := range (*smokersChannels) {
			(*smokersChannels)[i] <- true;
			if <-cigaretteSmokeChannel {
				break;
			}
		}

	}

}

func putItems(table *table_t, a int, b int) {

	(*table)[a] = true;
	(*table)[b] = true;

}

func clearTable(table *table_t)  {
	for i := range *table {
		(*table)[i] = false;
	}

}

func main() {

	var smokersChannels = make([]chan bool, 3);
	var cigaretteSmokeChannel = make(chan bool);

	table := table_t {
		false,
		false,
		false,
	};

	for i := range smokersChannels {
		smokersChannels[i] = make(chan bool);
	}
	go dealer(&smokersChannels, cigaretteSmokeChannel, &table);
	go smoker(&smokersChannels, cigaretteSmokeChannel, PAPER, &table);
	go smoker(&smokersChannels, cigaretteSmokeChannel, FIRE, &table);
	go smoker(&smokersChannels, cigaretteSmokeChannel, TABACO, &table);

	fmt.Scanln();
}