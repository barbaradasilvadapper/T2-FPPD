//GRUPO:
//Ana Carolina Poletto
//Anthony Antonelli
//Barbara Dapper
//Diego de Mello
//Louise Northfleet

//-------------------------------------------------------------------------
//PROBLEMA:
// Students in the dining hall invoke dineand then leave. After invoking dine
// and before invoking leave a student is considered “ready to leave”.
// The synchronization constraint that applies to students is that, in order to
// maintain the illusion of social suave, a student may never sit at a table alone. A
// student is considered to be sitting alone if everyone else who has invoked dine
// invokes leave before she has finished dine.
// Puzzle: write code that enforces this constraint.

// -------------------------------------------------------------------------
// CÓDIGO COM SEMAFOROS:
package main

import (
	"T2-FPPD/FPPDSemaforo"
	"fmt"
	"time"
)

// semáforos
var mutex = FPPDSemaforo.NewSemaphore(1)
var okToLeave = FPPDSemaforo.NewSemaphore(0)

// contadores globais
var eating int = 0
var readyToLeave int = 0

func getFood(id int) {
	fmt.Printf("Estudante %d pegou comida.\n", id)
	fmt.Printf("Estudante %d está comendo.\n", id)
}

func leave(id int) {
	fmt.Printf("Estudante %d saiu.\n", id)
}

func student(id int) {
	// pegar comida
	mutex.Wait()
	eating++
	if eating == 2 && readyToLeave == 1 {
		okToLeave.Signal()
		readyToLeave--
	}
	mutex.Signal()
	getFood(id)
	time.Sleep(time.Millisecond * 500) // simula o tempo de comer

	// sair
	mutex.Wait()
	eating--
	readyToLeave++

	if eating == 1 && readyToLeave == 1 {
		mutex.Signal()
		okToLeave.Wait()
	} else if eating == 0 && readyToLeave == 2 {
		okToLeave.Signal()
		readyToLeave -= 2
		mutex.Signal()
	} else {
		readyToLeave--
		mutex.Signal()
	}

	leave(id)
}

func main() {
	for i := 1; i <= 5; i++ {
		go student(i)
	}

	// espera todos terminarem
	time.Sleep(time.Second * 5)
	fmt.Println("Todos os estudantes saíram.")
}
