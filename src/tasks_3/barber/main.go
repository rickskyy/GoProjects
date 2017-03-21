package tasks_3

/*
Задача о парикмахере. В тихом городке есть парикмахерская. Салон парикмахерской мал,
ходить там может только парикмахер и один посетитель. Парикмахер всю жизнь обслуживает
посетителей. Когда в салоне никого нет, он спит в кресле. Когда посетитель приходит и видит
спящего парикмахера, он будет его, садится в кресло и спит, пока парикмахер занят стрижкой.
Если посетитель приходит, а парикмахер занят, то он встает в очередь и засыпает. После
стрижки парикмахер сам провожает посетителя. Если есть ожидающие посетители, то
парикмахер будит одного из них и ждет пока тот сядет в кресло парикмахера и начинает
стрижку. Если никого нет, он снова садится в свое кресло и засыпает до прихода посетителя.
Создать многопоточное приложение, моделирующее рабочий день парикмахерской.
 */

import (
	"fmt"
	"time"
)

var queue = make(chan chan string)

func barber () {
	for {
		fmt.Println("Barber: sleeping...")
		current := <-queue

		fmt.Println("Barber: waking next visitor...")
		current <- "COME_OVER"

		fmt.Printf("Barber: working on visitor %s\n", <-current)
		time.Sleep(300 * time.Millisecond)
		fmt.Println("Barber: done.")

		current <- "LEAVE"
		<- current
	}
}

func visitor (name string) {
	fmt.Printf("%s: has came...\n", name)

	speech := make(chan string)
	queue <- speech

	<-speech;
	speech <- name

	fmt.Printf("%s: sleeping while barber works...\n", name)

	<-speech;
	speech <- "LEFT"
	fmt.Printf("%s: leaving...\n", name)
}

func main () {
	go barber()

	time.Sleep(300 * time.Millisecond)
	go visitor("John")
	go visitor("Derek")
	go visitor("Stalin")
	time.Sleep(300 * time.Millisecond)
	go visitor("Petya")
	go visitor("Grishka")

	fmt.Scanln()
}
