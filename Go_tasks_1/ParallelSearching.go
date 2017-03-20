package main

import (
	"fmt"
	"time"
)

/* Первая задача о Винни-Пухе, или неправильные пчелы. Неправильные пчелы, подсчитав
в конце месяца убытки от наличия в лесу Винни-Пуха, решили разыскать его и наказать в
назидание всем другим любителям сладкого. Для поисков медведя они поделили лес на
участки, каждый из которых прочесывает одна стая неправильных пчел. В случае нахождения
медведя на своем участке стая проводит показательное наказание и возвращается в улей.
Если участок прочесан, а Винни-Пух на нем не обнаружен, стая также возвращается в улей.
Требуется создать многопоточное приложение, моделирующее действия пчел. При решении
использовать парадигму портфеля задач.
*/

var Matrix [][]int

// Here’s the worker, of which we’ll run several concurrent instances.
// These workers will receive work on the jobs channel and send the corresponding results on results. We’ll sleep a second per job to simulate an expensive task.
func worker(id int, jobs <-chan []int, results chan<- []int, terminateChannel <-chan bool) {
	for j := range jobs {
		select {
			case <-terminateChannel:
				return
			default:
				fmt.Println("bee group", id, "started job", j)
				time.Sleep(time.Second)
				var l []int
				for i:= range j {
					l = append(l, j[i])
				}
				results <- l
				fmt.Println("bee group", id, "finished job", j)
		}
	}
}

func fillingMatrix (n, m int) [][]int{
	matrix := make([][]int, n, m)
	for i := range matrix {
		matrix[i] = make([]int, m)
		for j := range matrix[i] {
			if i == 4 && j == 7 {
				matrix[i][j] = 1
			} else {
				matrix[i][j] = 0
			}
		}
	}
	fmt.Println(matrix)

	return matrix
}

func main() {
	var (
		n = 10
		m = 10
	)

	Matrix = fillingMatrix(n, m)
	terminateChannel := make(chan bool)
	jobs := make(chan []int, 100)
	results := make(chan []int, 100)

	//Workers initially blocked because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results, terminateChannel)
	}

	//Here we send jobs and then close that channel to indicate that’s all the work we have.
	for j := 0; j < n; j++ {
		jobs <- Matrix[j]
	}
	close(jobs)

	var terminateFlag bool = false
	//Finally we collect all the results of the work.
	for a := 0; a < n; a++ {
		var result []int = <-results
		fmt.Println(result)
		for i := range result {
			if result[i] == 1 {
				fmt.Println("Bear found")
				terminateChannel <- true
				terminateFlag = true
				goto terminate
			}
		}
		terminate:
			if terminateFlag {
				break
			}

	}
}
