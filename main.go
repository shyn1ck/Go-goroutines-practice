package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	HelloChannel("Vasya")
	NumberGenerator(10)
	totalSum := calculateSum(100)
	fmt.Println("Сумма чисел от 1 до", 100, "равна:", totalSum)
	ParallelMessages("Vasya", "Sadam")

}

func HelloChannel(name string) {
	messageChannel := make(chan string)

	go func() {
		messageChannel <- "Hello " + name
	}()

	go func() {
		msg := <-messageChannel
		fmt.Println("Message :", msg)
	}()

	time.Sleep(1 * time.Second)
}

func NumberGenerator(n int) {
	numberChannel := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= n; i++ {
			numberChannel <- i
		}
		close(numberChannel)
	}()

	go func() {
		defer wg.Done()
		for num := range numberChannel {
			fmt.Println("Получено число:", num)
		}
	}()
	wg.Wait()
}

func calculateSum(n int) int {
	numGoroutines := 4                       // Количество горутин (4 в данном случае)
	results := make(chan int, numGoroutines) // Канал для сбора результатов
	var wg sync.WaitGroup                    // Используем WaitGroup для ожидания завершения всех горутин

	// Определяем шаг диапазона для каждой горутины
	step := n / numGoroutines

	// Запускаем горутины для вычисления суммы в каждом диапазоне
	for i := 0; i < numGoroutines; i++ {
		start := i*step + 1
		end := (i + 1) * step

		// Последняя горутина должна обработать оставшиеся числа
		if i == numGoroutines-1 {
			end = n
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении работы горутины
			sum := 0
			for j := start; j <= end; j++ {
				sum += j
			}
			results <- sum // Отправляем результат в канал
		}(start, end)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Подсчитываем итоговую сумму, собирая результаты из канала
	totalSum := 0
	for result := range results {
		totalSum += result
	}

	return totalSum
}

func ParallelMessages(msg1, msg2 string) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go func() {
		channel1 <- msg1
	}()

	go func() {
		channel2 <- msg2
	}()
	firstMessage := <-channel1
	secondMessage := <-channel2
	fmt.Println("Получено:", firstMessage)
	fmt.Println("Получено:", secondMessage)
}
