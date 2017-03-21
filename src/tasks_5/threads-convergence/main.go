package main

import (
	"time"
	"fmt"
	"math/rand"
)

/*
с) Создать приложение с тремя потоками. Каждый поток работает со
своим массивом, потоки проверяют сумму элементов своего массива с суммами элементов других потоков и останавливаются,
когда все три суммы равны между собой. Если суммы не равны, каждый поток прибавляет единицу к одному элементу массива
или отнимает единицу от одного элемента массива, затем снова проверяет условие равенства сумм. На момент останов-
ки всех трех потоков, суммы элементов массивов должны быть одинаковы.
*/


const SIZE = 500;
const MAX = 100;

var seed = rand.NewSource(time.Now().UnixNano());
var random = rand.New(seed);

type custom_array struct {
	data       * []int;
	sum_change chan int;
	sum        chan int;
	terminate  chan bool;
}

func handler(arrays* []custom_array) {
	for true {
		var array_sums = make([]int, 3);
		result := true;
		prevSum := 0;
		for i := range *arrays {
			array_sums[i] = <-(*arrays)[i].sum;
			if (i > 0 && array_sums[i] != prevSum) {
				result = false;
			}
			prevSum = array_sums[i];
		}

		if (result) {
			for i := range *arrays {
				(*arrays)[i].terminate <- true;

			}
		} else {
			for i := range *arrays {
				(*arrays)[i].sum_change <- getChangeValue(i, &array_sums);
			}
			fmt.Println("SUCCESS");
		}
	}
}


func generateArray() custom_array {
	result_array := make([]int, SIZE)
	for i := range result_array {
		result_array[i] = random.Intn(MAX);
	}
	result_channel := make(chan int);
	terminate := make(chan bool);
	sum := make(chan int);
	return custom_array{
		data: &result_array,
		sum_change: result_channel,
		terminate: terminate,
		sum: sum,
	}
}

func changeValue(array* []int, value int) {

	index := random.Intn(SIZE);
	(*array)[index] += value;
	fmt.Printf("Index %d changed on value %d\n", index, value);
}

func getChangeValue(index int, sums* []int) int {
	dif := 0;
	for i := range *sums {
		if (i == index) {
			continue;
		}

		dif += (*sums)[index] - (*sums)[i];
	}
	fmt.Printf("Diff = %d\n", dif);
	if (dif > 0) {
		return -1;
	}
	return 1;
}

func worker(array* custom_array) {
	for true {
		sum := 0;
		for i := range *array.data {
			sum += (*array.data)[i];
		}
		fmt.Printf("NOW SUM: %d\n", sum);
		//fmt.Scanln();
		array.sum <-sum;
		select {
		case <-array.terminate:
			fmt.Printf("Terminate with sum %d", sum);
			return;
		default:
			changeValue(array.data, <-array.sum_change);
		}
	}
}

func main() {
	var arrays = make([]custom_array, 3, SIZE);
	for i := range arrays {
		arrays[i] = generateArray();
		go worker(&arrays[i]);
		fmt.Println(arrays[i].data)
	}

	go handler(&arrays);
	fmt.Scanln();

}
