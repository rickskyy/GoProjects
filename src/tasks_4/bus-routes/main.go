package main

import (
	"time"
	"sync"
	"math/rand"
	"fmt"
	"bytes"
)

/*
Задача про автобус.
Создать многопоточное приложение, работающее с общим графом.
Для защиты операций с графом использовать блокировки чтения-записи.
Граф описывает множество городов и множество рейсов автобусов от города
А к городу Б с указанием цены билета (по умолчанию, если рейс от А к Б, он
идет и от Б к А, с одинаковой ценой). В приложении должны работать следующие потоки:
1) поток, изменяющий цену билета;
2) поток, удаляющий и добавляющий рейсы между городами;
3) поток, удаляющий старые города и добавляющий новые;
4) потоки, определяющие есть ли путь от произвольного города А до произвольного города Б, и
какова цена такой поездки (если прямого пути нет, то найти любой путь из существующих)
*/

var seed = rand.NewSource(time.Now().UnixNano())
var random = rand.New(seed)

var name_gen = [...]string {"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}

type journey_t struct {
	from int
	to int
	price [] int
	sync.RWMutex
}

type journeys_t struct {
	list [] journey_t
	cities [] string
	sync.RWMutex
}

func dfs(j* journeys_t, from_city int, to_city int) ([] int, int) {
	return dfs_recur(j, from_city, to_city, map[int]bool{}, make([]int, 0), 0)
}

func dfs_recur(j* journeys_t, cur int, to_city int, v map[int]bool, path [] int, summary int) ([] int, int) {
	v[cur] = true

	if(cur == to_city) {
		return path, summary
	}

	for journey_id := range j.list {

		//fmt.Printf("Journey from %d: \n %+v \n\n", j.list[journey_id].from, j.list[journey_id])

		if(j.list[journey_id].from == cur && !v[j.list[journey_id].to]) {

			min := j.list[journey_id].price[0]
			for i := range j.list[journey_id].price {
				if(j.list[journey_id].price[i] < min) {
					min = j.list[journey_id].price[i]
				}
			}

			return dfs_recur(j, j.list[journey_id].to, to_city, v, append(path, j.list[journey_id].to), summary + min)
		}
	}

	return make([]int, 0), -1
}

func _randomize_city(j* journeys_t) string {
	return j.cities[random.Intn(len(j.cities))]
}

func _randomize_name() string {
	var name bytes.Buffer

	for i := 1; i <= 8; i++ {
		name.WriteString(name_gen[random.Intn(len(name_gen))])
	}

	return name.String();
}

func _get_city_by_name(j* journeys_t, name string) int {
	for id, city := range j.cities {
		if(city == name) {
			return id;
		}
	}

	return -1;
}

func _add_city(j* journeys_t, name string) {
	j.cities = append(j.cities, name)
}

func _remove_city(j* journeys_t, name string) {

	id := _get_city_by_name(j, name)

	if id < 0 {
		return
	}

	for journey_id := 0; journey_id < len(j.list); journey_id++ {
		if(id == j.list[journey_id].from || id == j.list[journey_id].to) {
			j.list = append(j.list[:journey_id], j.list[journey_id+1:]...)
			journey_id--;
		}
	}

	j.cities = append(j.cities[:id], j.cities[id+1:]...)
}

func _get_path(j* journeys_t, city_from int, city_to int) int {

	for id, journey := range j.list {
		if(city_from == journey.from && city_to == journey.to) {
			return id;
		}
	}

	return -1;
}

func _add_path(j* journeys_t, city_from int, city_to int, price int) {

	path_id := _get_path(j, city_from, city_to)

	if(path_id < 0) {
		j.list = append(j.list, journey_t {
			from: city_from,
			to: city_to,
			price: []int{price},
		})
	} else {
		j.list[path_id].price = append(j.list[path_id].price, price)

	}
}

func _remove_path(j* journeys_t, city_from int, city_to int, price int) {

	path_id := _get_path(j, city_from, city_to)

	if(path_id < 0) {
		return
	}

	if(len(j.list[path_id].price) < 2) {
		j.list = append(j.list[:path_id], j.list[path_id+1:]...)
		return
	}

	for i, price_item := range j.list[path_id].price {
		if(price_item == price) {
			j.list[path_id].price = append(j.list[path_id].price[:i], j.list[path_id].price[i+1:]...)
			return
		}
	}
}

func _find_path(j* journeys_t, from_city string, to_city string) ([] int, int) {
	return dfs(j, _get_city_by_name(j, from_city), _get_city_by_name(j, to_city))
}

func _get_prices(j* journeys_t, city_from int, city_to int) []int {

	journey_id := _get_path(j, city_from, city_to)
	if(journey_id > 0) {
		return j.list[journey_id].price
	}

	return []int{}
}

func _change_price(j* journeys_t, city_from int, city_to int, price int, end_price int) {
	journey_id := _get_path(j, city_from, city_to)
	if(journey_id > 0) {
		for i, price_item := range j.list[journey_id].price {
			if(price_item == price) {
				j.list[journey_id].price[i] = end_price
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////

func change_ticket_price(j* journeys_t) {
	for {
		j.Lock()

		from := _randomize_city(j)
		to := _randomize_city(j)

		if(from == to) {
			j.Unlock()
			continue
		}

		prices := _get_prices(j, _get_city_by_name(j, from), _get_city_by_name(j, to));
		if(len(prices) == 0) {
			j.Unlock()
			continue
		}
		price := prices[random.Intn(len(prices))]
		end_price := random.Intn(10)
		_change_price(j, _get_city_by_name(j, from), _get_city_by_name(j, to), price, end_price);

		fmt.Printf("Change price for %s to %s: \n price: %d to %d \n\n", from, to, price, end_price)

		j.Unlock()
		time.Sleep(400 * time.Millisecond)
	}
}

func change_journeys(j* journeys_t) {
	for {
		j.Lock()

		from := _randomize_city(j)
		time.Sleep(10 * time.Millisecond)
		to := _randomize_city(j)

		if(from == to) {
			j.Unlock()
			continue
		}

		if random.Intn(10) > 4 {
			//add
			price := random.Intn(9) + 1
			_add_path(j, _get_city_by_name(j, from), _get_city_by_name(j, to), price)
			fmt.Printf("Add Journey from %s to %s: \n price: %d \n\n", from, to, price)
		} else {
			//remove
			prices := _get_prices(j, _get_city_by_name(j, from), _get_city_by_name(j, to));
			if(len(prices) == 0) {
				j.Unlock()
				continue
			}
			price := prices[random.Intn(len(prices))]
			_remove_path(j, _get_city_by_name(j, from), _get_city_by_name(j, to), price);
			fmt.Printf("Remove Journey from %s to %s: \n price: %d \n\n", from, to, price)
		}
		j.Unlock()

		time.Sleep(500 * time.Millisecond)
	}
}

func change_cities(j* journeys_t) {
	for {
		j.Lock()

		old_city := _randomize_city(j)
		new_city := _randomize_name()

		_remove_city(j, old_city)
		_add_city(j, new_city)

		fmt.Printf("Removed %s and added %s \n\n", old_city, new_city)

		j.Unlock()

		time.Sleep(6000 * time.Millisecond)
	}
}

func check_journey(j* journeys_t) {
	for {
		j.Lock()

		from := _randomize_city(j)
		to := _randomize_city(j)
		if(from == to) {
			j.Unlock()
			continue
		}
		path, sum := _find_path(j, from, to);

		fmt.Printf("Check Journey from %s to %s: \n sum: %d \n path: %v \n\n", from, to, sum, path)

		j.Unlock()

		time.Sleep(2000 * time.Millisecond)
	}
}

func main() {
	j := &journeys_t {
		list: make([]journey_t, 0),
		cities: make([]string, 0),
	}


	_add_city(j, "A")
	_add_city(j, "B")
	_add_city(j, "C")
	_add_city(j, "D")

	_add_path(j, _get_city_by_name(j, "A"), _get_city_by_name(j, "B"), 10)
	_add_path(j, _get_city_by_name(j, "A"), _get_city_by_name(j, "B"), 5)
	_add_path(j, _get_city_by_name(j, "B"), _get_city_by_name(j, "C"), 10)
	_add_path(j, _get_city_by_name(j, "C"), _get_city_by_name(j, "D"), 10)
	_add_path(j, _get_city_by_name(j, "D"), _get_city_by_name(j, "A"), 5)
	_add_path(j, _get_city_by_name(j, "C"), _get_city_by_name(j, "A"), 8)

	fmt.Printf("j: %+v\n", j)

	go change_cities(j)
	go change_journeys(j)
	go change_ticket_price(j)
	go check_journey(j)

	fmt.Scanln()
}


