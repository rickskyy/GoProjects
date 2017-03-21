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

type sum_array_t struct {
	data       [10]int;
	sum        int;
}

func check_sum(arrays* [3]sum_array_t, diffs *[3][3]int) bool {

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			(*diffs)[i][j] = (*arrays)[i].sum - (*arrays)[j].sum;
		}
	}

	for i := 1; i < 3; i++ {
		if ((*diffs)[0][i] != 0) {
			return true;
		}
	}

	return false;
}

func change_value(array* [10]int, diffs* [3]int) {

	sum := 0;

	for i:=0; i < 3; i++ {
		sum += (*diffs)[i];
	}

	index := random.Intn(10);

	if(sum > 0) {
		(*array)[index] -= 1;
	} else {
		(*array)[index] += 1;
	}
}

func worker(array *sum_array_t, diffs* [3]int, done chan int) {

	change_value(&array.data, diffs);

	(*array).sum = 0;
	for i := 0; i < 3; i++ {
		(*array).sum += (*array).data[i];
	}

	done <- 1;
}

func main() {
	var arrays [3]sum_array_t;
	var diffs [3][3]int;
	var done = make(chan int);

	for i := 0; i < len(arrays); i++ {
		sum := 0;

		var result_array [10]int;

		for j := 0; j < 10; j++ {
			result_array[j] = random.Intn(MAX);
			sum += result_array[j]
		}

		arrays[i] = sum_array_t{
			data: result_array,
			sum: sum,
		}
	}

	fmt.Printf("%v\n\n", arrays);

	for check_sum(&arrays, &diffs) {

		fmt.Printf("%v\n\n", diffs);

		for i := 0; i < 3; i++ {
			go worker(&arrays[i], &diffs[i], done);
			fmt.Println(arrays[i].data)
		}

		for i := 0; i < 3; i++ {
			<- done;
		}
	}

	fmt.Printf("%v\n\n", diffs);
}
