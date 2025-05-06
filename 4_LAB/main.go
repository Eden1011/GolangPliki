package main

import (
	"errors"
	"fmt"
	"math/rand"
	"shop/structs"
	"sync"
	"time"
)

const (
	numWorkers    int     = 10
	maxRetries    int     = 2
	orderCount    int     = 20
	failureChance float64 = 0.25
)

func randTime(min, max int) time.Duration {
	randMs := rand.Intn(max-min+1) + min
	return time.Duration(randMs) * time.Millisecond
}

func processOrder(order structs.Order) structs.ProcessResult {
	randTime := randTime(1, 1000)
	time.Sleep(randTime)

	success := rand.Float64() > failureChance
	var err error = nil
	if !success {
		err = errors.New("Could not process order")
	}

	return structs.ProcessResult{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		Success:      success,
		ProcessTime:  randTime,
		Error:        err,
	}
}

func getOrderDetails(many int) ([]string, float64, string) {
	products := map[string]float64{
		"PC":     1500,
		"Laptop": 1100,
		"iPad":   325,
		"iPhone": 210,
		"Mouse":  15,
	}
	productNames := make([]string, 0, len(products))
	for name := range products {
		productNames = append(productNames, name)
	}
	items := make([]string, 0, many)
	var amount float64 = 0

	for range many {
		randIndex := rand.Intn(len(productNames))
		randName := productNames[randIndex]
		items = append(items, randName)
		amount += products[randName]
	}

	customerNames := [...]string{"Janusz", "Bartlomiej", "Ryba"}
	customerIndex := rand.Intn(len(customerNames))
	customerName := customerNames[customerIndex]

	return items, amount, customerName
}

func generateOrder(id int) structs.Order {
	items, amount, customerName := getOrderDetails(rand.Intn(5) + 1)
	return structs.Order{
		ID:           id,
		CustomerName: customerName,
		Items:        items,
		TotalAmount:  amount,
	}
}

func orderGenerator(orderCh chan<- structs.Order, count int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= count; i++ {
		order := generateOrder(i)
		orderCh <- order
		fmt.Printf("Generated order number %d for %s, with items: %v and sum of %.2f\n",
			i, order.CustomerName, len(order.Items), order.TotalAmount)
		time.Sleep(randTime(1, 255))
	}
	close(orderCh)
}

func worker(id int, orderCh <-chan structs.Order, resultCh chan<- structs.ProcessResult, retryCh chan<- structs.Order, wg *sync.WaitGroup) {
	defer wg.Done()

	for order := range orderCh {
		fmt.Printf("Worker number %d is processing order number %d for customer %s\n", id, order.ID, order.CustomerName)

		result := processOrder(order)

		if !result.Success {
			fmt.Printf("Worker number %d failed with order number %d for customer %s, retrying now...\n", id, order.ID, order.CustomerName)
			retryCh <- order
		} else {
			fmt.Printf("Worker number %d successfully processed order number %d for customer %s\n", id, order.ID, order.CustomerName)
		}
		resultCh <- result
	}
}

func retry(retryCh <-chan structs.Order, resultCh chan<- structs.ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()

	retryCount := make(map[int]int)

	for order := range retryCh {
		retryCount[order.ID]++

		if retryCount[order.ID] <= maxRetries {
			fmt.Printf("Retry attempt (%d/%d) for order #%d\n",
				retryCount[order.ID], maxRetries, order.ID)

			processingTime := randTime(200, 800)
			time.Sleep(processingTime)

			success := rand.Float64() > failureChance
			var err error

			if !success && retryCount[order.ID] < maxRetries {
				err = errors.New("retry failed")
				go func(o structs.Order) {
					time.Sleep(randTime(300, 600))
				}(order)
			} else if !success {
				err = errors.New("all retry attempts failed")
				resultCh <- structs.ProcessResult{
					OrderID:      order.ID,
					CustomerName: order.CustomerName,
					Success:      false,
					ProcessTime:  processingTime,
					Error:        err,
				}
			} else {
				fmt.Printf("Retry successful for order #%d\n", order.ID)
				resultCh <- structs.ProcessResult{
					OrderID:      order.ID,
					CustomerName: order.CustomerName,
					Success:      true,
					ProcessTime:  processingTime,
					Error:        nil,
				}
			}
		}
	}
}

func collectResults(resultCh <-chan structs.ProcessResult, done chan<- bool, orderCount int) {
	processed := make(map[int]bool)
	successful := 0
	failed := 0
	totalTime := time.Duration(0)

	for result := range resultCh {
		if _, exists := processed[result.OrderID]; exists {
			continue
		}

		processed[result.OrderID] = true

		if result.Success {
			successful++
			fmt.Printf("SUCCESS: Order #%d from %s processed in %v\n",
				result.OrderID, result.CustomerName, result.ProcessTime)
		} else {
			failed++
			fmt.Printf("FAILURE: Order #%d from %s failed: %v\n",
				result.OrderID, result.CustomerName, result.Error)
		}

		totalTime += result.ProcessTime

		if len(processed) >= orderCount {
			break
		}
	}

	total := successful + failed
	fmt.Println("\nORDER PROCESSING STATISTICS")
	fmt.Printf("Total orders: %d\n", total)
	fmt.Printf("Successful: %d (%.2f%%)\n", successful, float64(successful)/float64(total)*100)
	fmt.Printf("Failed: %d (%.2f%%)\n", failed, float64(failed)/float64(total)*100)

	if total > 0 {
		fmt.Printf("Average processing time: %v\n", totalTime/time.Duration(total))
	}

	done <- true
}

func main() {

	orderCh := make(chan structs.Order, orderCount)
	resultCh := make(chan structs.ProcessResult, orderCount*2)
	retryCh := make(chan structs.Order, orderCount)
	retryChSend := make(chan structs.Order, orderCount)

	done := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(1)
	go orderGenerator(orderCh, orderCount, &wg)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, orderCh, resultCh, retryChSend, &wg)
	}

	go func() {
		for order := range retryChSend {
			retryCh <- order
		}
	}()

	wg.Add(1)
	go retry(retryCh, resultCh, &wg)

	go collectResults(resultCh, done, orderCount)

	go func() {
		wg.Wait()
		close(retryChSend)
		close(retryCh)
		time.Sleep(time.Second * 2)
		close(resultCh)
	}()

	<-done
	fmt.Println("Processing complete.")
}
