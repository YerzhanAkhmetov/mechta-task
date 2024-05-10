package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const FileName = "data.json"
const CountValueArray = 1000000
const NumWorkers = 3 // Количество горутин-обработчиков

type Data struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	start := time.Now() // Записываем текущее время
	sum := 0
	data, err := createOrReadFile()
	if err != nil {
		fmt.Println("Ошибка при создании/чтении файла:", err)
		os.Exit(1)
	}

	// Создаем канал для передачи данных в горутины-обработчики
	dataCh := make(chan Data, CountValueArray)

	// Создаем мьютекс для синхронизации доступа к sum
	var mu sync.Mutex

	// Запускаем несколько горутин-обработчиков
	var wg sync.WaitGroup
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		go worker(dataCh, &sum, &mu, &wg)
	}

	// Передаем данные в канал для обработки
	for _, d := range data {
		dataCh <- d
	}
	close(dataCh) // Закрываем канал после передачи всех данных

	// Ожидаем завершения всех горутин-обработчиков
	wg.Wait()

	// Вывод данных
	fmt.Printf("Итоговая сумма : %d\n", sum)
	fmt.Printf("Время выполнения программы: %s\n", time.Since(start).String())
}

func worker(dataCh <-chan Data, sum *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range dataCh {
		sumDigit := d.A + d.B

		// Блокируем мьютекс для синхронизации доступа к sum
		mu.Lock()
		*sum += sumDigit
		mu.Unlock()
	}
}

func createOrReadFile() ([]Data, error) {
	// Проверяем на существование файла
	_, err := os.Stat(FileName)
	if err == nil {
		fmt.Println("Файл " + FileName + " уже существует")

		// Чтение данных из файла
		data, err := os.ReadFile(FileName)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", err)
			return nil, err
		}

		var dataArray []Data
		if err := json.Unmarshal(data, &dataArray); err != nil {
			fmt.Println("Ошибка при разборе данных из файла:", err)
			return nil, err
		}
		return dataArray, nil
	}

	start := time.Now() // Записываем текущее время

	// генератор случайных чисел
	dataArray, err := GenerateRandNumber(CountValueArray)
	if err != nil {
		fmt.Println("Ошибка при генерации случайных чисел:", err)
		return nil, err
	}

	// Преобразование данных в JSON
	jsonData, err := json.MarshalIndent(dataArray, "", "    ")
	if err != nil {
		fmt.Println("Ошибка преобразования данных в JSON:", err)
		return nil, err
	}

	// Запись JSON в файл
	err = os.WriteFile(FileName, jsonData, 0644) // 644 определяют права доступа
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return nil, err
	}

	fmt.Println("Данные успешно сохранены в файл " + FileName)

	// Вывод данных
	elapsed := time.Since(start) // Вычисляем прошедшее время
	fmt.Printf("Время выполнения создания файла: %s\n", elapsed)
	return dataArray, nil
}

func GenerateRandNumber(CountValueArray int) ([]Data, error) {
	if CountValueArray <= 0 {
		return nil, errors.New("CountValueArray must be greater than 0")
	}
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	// Создаем срез для хранения объектов Data
	dataArray := make([]Data, CountValueArray)

	// Заполняем массив случайными значениями
	for i := 0; i < CountValueArray; i++ {
		dataArray[i] = Data{
			A: randGen.Intn(21) - 10, // Генерируем случайное число в диапазоне [-10, 10]
			B: randGen.Intn(21) - 10, // Генерируем случайное число в диапазоне [-10, 10]
		}
	}
	return dataArray, nil
}
