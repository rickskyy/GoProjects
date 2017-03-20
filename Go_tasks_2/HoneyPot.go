package main
import (
	"fmt"
	"time"
	"math/rand"
)

/*
а) Задача о Винни-Пухе или правильные пчелы.
В одном лесу живут n пчел и один медведь, которые используют один горшок меда, вместимостью Н глотков.
Сначала горшок пустой. Пока горшок не наполнится, медведь спит.
Как только горшок заполняется, медведь просыпается и съедает весь мед, после чего снова засыпает.
Каждая пчела многократно собирает по одному глотку меда и кладет его в горшок.
Пчела, которая приносит последнюю порцию меда, будит медведя.
Создать многопоточное приложение, моделирующее поведение пчел и медведя.
 */

func main() {
	numberOfBees := 10
	volumeOfPot := 100

	fromBeeToPot := make(chan byte, 10)
	terminate := make(chan bool)
	wake := make(chan byte)

	pot := Pot{volumeOfPot, 0, fromBeeToPot, terminate, wake}

	bees := make([]Bee, numberOfBees, numberOfBees)
	for i:=0; i<numberOfBees; i++{
		bees[i].name = "Bee #" + fmt.Sprintf("%v", i)
		bees[i].bring = fromBeeToPot
		bees[i].terminate = terminate
	}

	bear := Bear{wake}

	for i:=0; i<numberOfBees; i++{
		go bees[i].Start()
	}
	go pot.Start()
	bear.Start()
}


type Pot struct {
	volume int
	collected int
	receive chan byte
	terminate chan bool
	wake chan byte
}

func (pot *Pot) Start(){
	for i:=range pot.receive{
		pot.collected += int(i)
		fmt.Printf("Pot collected %v honey\n", pot.collected)
		if pot.collected == pot.volume {
			pot.terminate <- true
			close(pot.receive)
			close(pot.terminate)
			fmt.Println("Pot is full.")
			pot.wake <- 1
			pot.collected = 0
			close(pot.wake)
		}
	}
}


type Bee struct {
	name string
	bring chan byte
	terminate chan bool
}

func (bee *Bee) Start(){
	for {
		select {
		case <-bee.terminate:
			return
		default:
			fmt.Printf("Bee %s brought honey.\n", bee.name)
			bee.bring <- 1
			var t uint32 = rand.Uint32() % 500
			time.Sleep((100 + time.Duration(t)) * time.Millisecond)
		}
	}
}

type Bear struct {
	wake chan byte
}

func (bear *Bear) Start(){
	for range bear.wake{
		fmt.Println("Bear waked up!")
	}
}
