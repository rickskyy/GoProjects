package main

import (
	"fmt"
	"time"
	"math/rand"
)

/*
б)Задача о парикмахере. В тихом городке есть парикмахерская. Салон парикмахерской мал, ходить там может только парикмахер и один посетитель. 
Парикмахер всю жизнь обслуживает посетителей. Когда в салоне никого нет, он спит в кресле. 
Когда посетитель приходит и видит спящего парикмахера, он будет его, садится в кресло и спит, пока парикмахер занят стрижкой. 
Если посетитель приходит, а парикмахер занят, то он встает в очередь и засыпает. После стрижки парикмахер сам провожает посетителя. 
Если есть ожидающие посетители, то парикмахер будит одного из них и ждет пока тот сядет в кресло парикмахера и начинает стрижку. 
Если никого нет, он снова садится в свое кресло и засыпает до прихода посетителя. 
Создать многопоточное приложение, моделирующее рабочий день парикмахерской.
 */

var seed = rand.NewSource(time.Now().UnixNano())

var random = rand.New(seed)

func send(client []Client, toBarber chan *Client){
	for i:=0; i<len(client); i++{
		toBarber <- &client[i]
		time.Sleep((200 + time.Duration(random.Intn(1000))) * time.Millisecond)
	}
	close(toBarber)
}

func main() {
	numberOfClients := 10
	toBarber := make(chan *Client, 10)
	fromBarber := make(chan *Client, 10)

	barber := Barber{toBarber, fromBarber}
	go barber.Start()

	client := make([]Client, numberOfClients, numberOfClients)
	for i:=0; i<numberOfClients; i++{
		client[i].name = "Client #" + fmt.Sprintf("%v", i)
	}

	go send(client, toBarber)

	for in := range fromBarber{
		fmt.Printf("Barber finished with client %s.\n", in.name)
	}
}

type Barber struct {
	get chan *Client
	release chan *Client
}

func (Barber *Barber) Start(){
	for client := range Barber.get{
		fmt.Printf("Barber working with client %s.\n", client.name)
		fmt.Println(" ")
		time.Sleep(400)
		Barber.release <- client
	}
	close(Barber.release)
}

type Client struct {
	name string
}
