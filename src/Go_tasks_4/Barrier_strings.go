package main

import (
	"sync"
	"fmt"
	"time"
	"math/rand"
	"bytes"
	"strings"
)

/*
б) Создать приложение с четырьмя потоками. Каждый поток работает с собственной строкой.
Строки могут содержать только символы А, B, C, D.
Поток может поменять символ А на С или С на А или В на D или D на В.
Потоки останавливаются когда общее количество символов А и В становится равным хотя бы для трех строк.
 */

var stringGenList = [...]string {"A", "B", "C", "D"}
var seed = rand.NewSource(time.Now().UnixNano())
var random = rand.New(seed)

type Strings struct {
	stringList []string
	sync.WaitGroup
}

func NewStrings(n, strSize int) *Strings{
	return &Strings{
		stringList: initializeList(n, strSize),
	}
}

func generateString(n int) string {
	var name bytes.Buffer
	for i := 1; i <= n; i++ {
		name.WriteString(stringGenList[random.Intn(len(stringGenList))])
	}
	return name.String()
}

func initializeList(n, strSize int) []string{
	stringList := make([]string, n)
	for i := range stringList {
		stringList[i] = generateString(strSize)
	}
	return stringList
}

func printStrings(list *Strings) {
	for _, i:= range list.stringList {
		fmt.Println(i)
	}
	fmt.Println("/--------------------------------------------")
}

func changeCharacter(list *Strings, group *sync.WaitGroup, index int) {
	var operation = rand.Intn(4)
	switch operation {
	case 0:  // A -> C
		list.stringList[index] = strings.Replace(list.stringList[index], "A", "C", 1)
	case 1:  // C -> A
		list.stringList[index] = strings.Replace(list.stringList[index], "C", "A", 1)
	case 2:  // B -> D
		list.stringList[index] = strings.Replace(list.stringList[index], "B", "D", 1)
	case 3:  // D -> B
		list.stringList[index] = strings.Replace(list.stringList[index], "D", "B", 1)
	}
	group.Done()
}

func checkStringRule(list *Strings, n int) bool{
	arrayA := make([]int, n)
	arrayB := make([]int, n)

	for j, i:= range list.stringList {
		arrayA[j] = strings.Count(i, "A")
		arrayB[j] = strings.Count(i, "B")
	}

	counter := 0
	for i:= 0; i < n; i++ {
		if arrayA[i] == arrayB[i] {
			counter++
		}
	}
	fmt.Printf("Counter: %d \n", counter)

	if counter >= 3 {
		return true
	} else {
		return false
	}
}

func StringSimulator(list *Strings, group *sync.WaitGroup, n int) {
	stopFlag := false
	for !stopFlag {
		group.Add(n)


		for i:=0; i < n; i++ {
			go changeCharacter(list, group, i)
		}

		if checkStringRule(list, n) {
			stopFlag = true
			fmt.Println("Strings matched the rule")
		}

		printStrings(list)
		time.Sleep(20*time.Millisecond)
		group.Wait()
	}

}

func main() {
	const (
		N = 4
		STR_SIZE = 10
	)

	list := NewStrings(N, STR_SIZE)
	group :=new (sync.WaitGroup)
	StringSimulator(list, group, N)
}






