package main

import (
	"fmt"
	"strings"
	"time"
)

/*
б) Создать приложение с четырьмя потоками. Каждый поток работает
с собственной строкой. Строки могут содержать только символы А, B, C, D.
Поток может поменять символ А на С или С на А или В на D или D на В. По-
токи останавливаются когда общее количество символов А и В становится
равным хотя бы для трех строк.
 */

func increase(str string, diff int) string {

	cs := strings.Count(str, "C");

	if(cs > 0) {
		str = strings.Replace(str, "C", "A", 1);
	} else {
		str = strings.Replace(str, "C", "A", 1);
	}

	return str;
}

func decrease(str string, diff int) string {
	as := strings.Count(str, "A");

	if(as > 1) {
		str = strings.Replace(str, "A", "C", 1);
	} else {
		str = strings.Replace(str, "A", "C", 1);
	}

	return str;
}

func str_diff(str1 *string, str2 *string) int {
	var diff int;

	for i := range (*str1) {
		if((*str1)[i] == 'A' || (*str1)[i] == 'B') {
			diff++
		}
	}

	for i := range (*str2) {
		if((*str2)[i] == 'A' || (*str2)[i] == 'B') {
			diff--
		}
	}

	return diff;
}

func check_equality(strings *[4]string, diff *[4][4]int) bool {

	for i:=0; i < 4; i++ {
		for j:=0; j < 4; j++ {
			(*diff)[i][j] = str_diff(&(*strings)[i], &(*strings)[j]);
			(*diff)[j][i] = -(*diff)[i][j];
		}
	}

	for i:=0; i < 4; i++ {
		s := 0
		for j:=0; j < 4; j++ {
			if(i!=j) {
				if(diff[i][j] == 0) {
					s++;
				}
			}
		}
		if(s >= 2) {
			return false;
		}
	}

	return true;
}

func worker(str *string, diff *[4]int, done chan int) {
	sum_diff := 0;

	for i:=0; i<4; i++ {
		sum_diff += (*diff)[i];
	}

	if(sum_diff < 0) {
		(*str) = increase((*str), (-sum_diff) / 3);
	} else {
		(*str) = decrease((*str), sum_diff / 3);
	}

	fmt.Printf("%s: %v\n", *str, *diff);

	done <- 1;
}

func main() {

	var diff [4][4]int;
	var strs [4]string;
	done := make(chan int)

	strs = [4]string {
		"DACDDD",
		"ACABAB",
		"AAAAAA",
		"CDCDCD",
	};

	for check_equality(&strs, &diff) {

		fmt.Printf("%v\n", diff);

		for i := 0; i < 4; i++ {
			go worker(&strs[i], &diff[i], done);
		}

		for i := 0; i < 4; i++ {
			<- done;
		}

		time.Sleep(time.Millisecond * 200);

		fmt.Print("\n");
	}

	fmt.Printf("%v\n", diff);

}