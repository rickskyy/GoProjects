package tasks_3

/**
Задача о Винни-Пухе или правильные пчелы. В одном лесу живут n пчел и один медведь,
которые используют один горшок меда, вместимостью Н глотков. Сначала горшок пустой.
Пока горшок не наполнится, медведь спит. Как только горшок заполняется, медведь
просыпается и съедает весь мед, после чего снова засыпает. Каждая пчела многократно
собирает по одному глотку меда и кладет его в горшок. Пчела, которая приносит последнюю
порцию меда, будит медведя. Создать многопоточное приложение, моделирующее поведение
пчел и медведя.
 */

import (
	"fmt"
	"time"
)

const N = 15
const n = 10

var pot int = 0
var done = make(chan bool)

func bee(id int) {
	fmt.Printf("Bee %d sleep...\n", id)

	for true {
		status := <-done

		if(status) {
			done <- status
			continue;
		}

		if pot >= N {
			done <- true
		} else {
			fmt.Printf("Bee %d wake\n", id)
			pot++
			time.Sleep(100 * time.Millisecond)
			done <- false
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("Bee %d done in %d\n", id, pot)
		}

		fmt.Printf("Bee %d sleep...\n", id)
	}
}


func bear() {
	fmt.Println("Bear sleep...")

	for {
		status := <-done

		if(!status) {
			done <- status
			continue;
		}

		fmt.Println("Bear wake")

		time.Sleep(2000 * time.Millisecond)
		pot = 0



		done <- false
		fmt.Println("Bear sleep...")
	}
}

func main()  {


	go bear()

	for i := 0; i < n; i++ {
		go bee(i)
	}

	done <- false;

	fmt.Scanln();
}