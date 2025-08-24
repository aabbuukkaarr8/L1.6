package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {

	fmt.Println("1. Выход по условию")
	stopByCondition()

	fmt.Println("2. Остановка через канал уведомления")
	stopByChannel()

	fmt.Println("3. Остановка через контекст")
	stopByContext()

	fmt.Println("4. Остановка через runtime.Goexit()")
	stopByGoexit()

	fmt.Println("5. Остановка через sync.WaitGroup")
	stopByWaitGroup()

	fmt.Println("6. Остановка через select с таймером")
	stopByTimer()

	fmt.Println("7. Остановка через panic/recover")
	stopByPanic()

	fmt.Println("8. Остановка через os.Exit")
}
func stopByCondition() {
	done := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			select {
			case <-done:
				fmt.Println("   Горутина остановлена по условию")
				return
			default:
				fmt.Printf("   Работаю... %d\n", i)
				time.Sleep(100 * time.Millisecond)
			}
		}
		fmt.Println("   Горутина завершилась естественным образом")
	}()

	time.Sleep(300 * time.Millisecond)
	close(done)
	time.Sleep(100 * time.Millisecond)
}

func stopByChannel() {
	stopCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-stopCh:
				fmt.Println("   Горутина остановлена через канал")
				return
			default:
				fmt.Println("   Работаю...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)
	stopCh <- struct{}{}
	time.Sleep(100 * time.Millisecond)
}
func stopByContext() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("   Горутина остановлена через контекст")
				return
			default:
				fmt.Println("   Работаю...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func stopByGoexit() {
	done := make(chan bool)

	go func() {
		defer func() {
			fmt.Println("   Горутина завершается через runtime.Goexit()")
			done <- true
		}()

		for i := 0; i < 3; i++ {
			fmt.Printf("   Работаю... %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}

		runtime.Goexit()
	}()

	<-done
}

func stopByWaitGroup() {
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-stopCh:
				fmt.Println("   Горутина остановлена через WaitGroup")
				return
			default:
				fmt.Println("   Работаю...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)
	stopCh <- struct{}{}
	wg.Wait()
}

func stopByTimer() {
	timer := time.NewTimer(300 * time.Millisecond)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("   Горутина остановлена по таймеру")
				return
			default:
				fmt.Println("   Работаю...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(400 * time.Millisecond)
}

// 7. Остановка через panic/recover
func stopByPanic() {
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("   Горутина остановлена через panic: %v\n", r)
			}
			done <- true
		}()

		for i := 0; i < 3; i++ {
			fmt.Printf("   Работаю... %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}

		panic("искусственная остановка")
	}()

	<-done
}
