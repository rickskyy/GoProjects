package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

/*а)Задача о новобранцах. Строю новобранцев дается команда «налево» или «направо».
Все новобранцы стараются исполнить приказ, но проблема в том, что они не знают где право, а где лево.
Следовательно, каждый новобра- нец поворачивается либо направо, либо налево. Если новобранец повернулся и видит,
что его сосед стоит к нему спиной, он считает, что все сделал пра- вильно. Если же они сталкиваются лицом к лицу,
то оба считают, что ошиб- лись, и разворачиваются на 180 градусов.
Создать многопоточное приложе- ние, моделирующее поведение строя новобранцев,
пока он не придет к ста- ционарному состоянию. Количество новобранцев ≥ 100.
Отдельный поток отвечает за часть строя не менее 50 новобранцев.
*/

var seed = rand.NewSource(time.Now().UnixNano())
var random = rand.New(seed)

type ArraySoldiers struct {
	arrayList [] bool
	sync.WaitGroup
}

func NewArraySoldiers(n int) *ArraySoldiers{
	return &ArraySoldiers{
		arrayList: initializeArraySoldiers(n),
	}
}

func initializeArraySoldiers(n int) []bool{
	array := make([]bool, n)
	for i := 0; i < n; i++ {
		array[i] = true
	}
	elemToChange := random.Intn(n)
	array[elemToChange] = false
	elemToChange = random.Intn(n)
	array[elemToChange] = false
	fmt.Println("False: ", elemToChange)
	return array
}

func rotateSoldiers(array *ArraySoldiers, group *sync.WaitGroup, start int, arraySize int){
	for i := start; i < arraySize; i++ {
		if array.arrayList[i] != array.arrayList[i+1] {
			array.arrayList[i] =! array.arrayList[i]
			array.arrayList[i+1] =! array.arrayList[i+1]
			i++
		}
	}
	group.Done()
}

func main() {
	const (
		N = 100
		N_THREADS = 2
	)
	array := NewArraySoldiers(N)
	group :=new (sync.WaitGroup)

	stopFlag := false
	for !stopFlag {
		group.Add(N_THREADS)
		stopFlag = true

		if array.arrayList[0] != array.arrayList[1] {
			array.arrayList[1] =! array.arrayList[1]
		}

		for i := 1; i < N-1; i++ {
			if array.arrayList[i] != array.arrayList[i+1]{
				stopFlag = false
			}
			fmt.Print(array.arrayList[i], " ")
		}
		fmt.Println("/-------------------------------------------------------------")

		if array.arrayList[N/2-1] != array.arrayList[N/2] {
			array.arrayList[N/2-1] = !array.arrayList[N/2-1]
			array.arrayList[N/2] = !array.arrayList[N/2]
		}

		for i := 0; i < N_THREADS; i++ {
			go rotateSoldiers(array, group, i*N/2+1, N-1)
		}
		group.Wait()
	}
}