//GRUPO:
//Ana Carolina Poletto
//Anthony Antonelli
//Barbara Dapper
//Diego de Mello
//Louise Northfleet

//-------------------------------------------------------------------------
//PROBLEMA:
// The next classical problem, called the Reader-Writer Problem, pertains to any
// situation where a data structure, database, or file system is read and modified
// by concurrent threads. While the data structure is being written or modified
// it is often necessary to bar other threads from reading, in order to prevent a
// reader from interrupting a modification in progress and reading inconsistent or
// invalid data.
// As in the producer-consumer problem, the solution is asymmetric. Readers
// and writers execute diﬀerent code before entering the critical section. The
// synchronization constraints are:
// 1. Any number of readers can be in the critical section simultaneously.
// 2. Writers must have exclusive access to the critical section.
// In other words, a writer cannot enter the critical section while any other
// thread (reader or writer) is there, and while the writer is there, no other thread
// may enter.

// -------------------------------------------------------------------------
// CÓDIGO COM SEMAFOROS:

package main

import (
	"T2-FPPD/FPPDSemaforo"
	"fmt"
	"time"
)

var readers = 0 // número de leitores
var mutex = FPPDSemaforo.NewSemaphore(1) // protege acesso a readers
var roomEmpty = FPPDSemaforo.NewSemaphore(1) // controla acesso à sala

// escritores
func writers(id int) {
	// tenta entrar na sala
	roomEmpty.Wait()
		// escreve
		fmt.Printf("Escritor %d está escrevendo\n", id)
		time.Sleep(time.Second)
		fmt.Printf("Escritor %d terminou de escrever\n", id)
	// sai da sala
	roomEmpty.Signal()
}

// leitores
func readersFunc(id int) {
	// tenta entrar na sala
	mutex.Wait()
		// incrementa número de leitores
		readers += 1
		// se for o primeiro, bloqueia a sala para escritores
		if readers == 1 {
			roomEmpty.Wait()
		}
	mutex.Signal()

	// lê
	fmt.Printf("Leitor %d está lendo\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Leitor %d terminou de ler\n", id)
	
	// sai da sala
	mutex.Wait()
		// decrementa número de leitores
		readers -= 1
		if readers == 0 {
			roomEmpty.Signal()
		}
	mutex.Signal()
}

func main() {
	// inicia escritores
	for i := 1; i <= 2; i++ {
		go writers(i)
	}

	// inicia leitores
	for i := 1; i <= 3; i++ {
		go readersFunc(i)
	}
	
	time.Sleep(10 * time.Second) // tempo para simulação
}