package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/YerzhanAkhmetov/task/domain"
	internal "github.com/YerzhanAkhmetov/task/internal/cns"
)

func Start() {

	start := time.Now()
	sum := 0
	data, err := createOrReadFile()
	if err != nil {
		fmt.Println("Ошибка при создании/чтении файла:", err)
		os.Exit(1)
	}

	if len(data) != internal.CountValueArray {
		fmt.Println("Неверное количество значений в файле")
		os.Exit(1)
	}

	// Создаем канал для передачи данных в горутины-обработчики
	jobsData := make(chan domain.Data, internal.CountValueArray)
	res := make(chan int, internal.CountValueArray)

	// Запускаем несколько горутин-обработчиков
	for i := 0; i < internal.NumWorkers; i++ {
		go worker(jobsData, res)
	}

	// Передаем данные в канал для обработки
	for _, d := range data {
		jobsData <- d
	}
	close(jobsData)

	// Вывод данных
	for i := 0; i < internal.CountValueArray; i++ {
		sum += <-res
	}
	fmt.Printf("Итоговая сумма : %d\n", sum)
	fmt.Printf("Время выполнения программы: %s\n", time.Since(start).String())
}

func worker(jobsData <-chan domain.Data, res chan<- int) {
	for d := range jobsData {
		sumDigit := d.A + d.B
		res <- sumDigit
	}
}

func createOrReadFile() ([]domain.Data, error) {
	// Проверяем на существование файла
	_, err := os.Stat(internal.FileName)
	if err == nil {
		fmt.Println("Файл " + internal.FileName + " уже существует")

		// Чтение данных из файла
		data, err := os.ReadFile(internal.FileName)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", err)
			return nil, err
		}

		var dataArray []domain.Data
		if err := json.Unmarshal(data, &dataArray); err != nil {
			fmt.Println("Ошибка при разборе данных из файла:", err)
			return nil, err
		}
		return dataArray, nil
	}

	//Создаем файл
	start := time.Now()

	// генератор случайных чисел
	dataArray, err := GenerateRandNumber(internal.CountValueArray)
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
	err = os.WriteFile(internal.FileName, jsonData, 0644) // 644 определяют права доступа
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return nil, err
	}

	fmt.Println("Данные успешно сохранены в файл " + internal.FileName)
	fmt.Printf("Время выполнения создания файла: %s\n", time.Since(start))

	return dataArray, nil
}

func GenerateRandNumber(CountValueArray int) ([]domain.Data, error) {
	if CountValueArray <= 0 {
		return nil, errors.New("CountValueArray must be greater than 0")
	}
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	// Создаем срез для хранения объектов Data
	dataArray := make([]domain.Data, CountValueArray)

	// Заполняем массив случайными значениями
	for i := 0; i < CountValueArray; i++ {
		dataArray[i] = domain.Data{
			A: randGen.Intn(21) - 10, // Генерируем случайное число в диапазоне [-10, 10]
			B: randGen.Intn(21) - 10, // Генерируем случайное число в диапазоне [-10, 10]
		}
	}
	return dataArray, nil
}
