//GRUPO:
//Ana Carolina Poletto
//Anthony Antonelli
//Barbara Dapper
//Diego de Mello
//Louise Northfleet

//-------------------------------------------------------------------------
//PROBLEMA:
// A tribe of savages eats communal dinners from a large pot that
// can hold M servings of stewed missionary1. When a savage wants to
// eat, he helps himself from the pot, unless it is empty. If the pot is
// empty, the savage wakes up the cook and then waits until the cook
// has refilled the pot.
// Any number of savage threads run the following code:
// Listing 5.1: Unsynchronized savage code
// 1 while True:
// 2 getServingFromPot()
// 3 eat()
// And one cook thread runs this code:
// Listing 5.2: Unsynchronized cook code
// 1 while True:
// 2 putServingsInPot(M)
// The synchronization constraints are:
// Savages cannot invoke getServingFromPot if the pot is empty.
// The cook can invoke putServingsInPot only if the pot is empty.
// Puzzle: Add code for the savages and the cook that satisfies the synchro-
// nization constraints.

// -------------------------------------------------------------------------
// CÓDIGO COM SEMAFOROS:

package main

import (
	"T2-FPPD/FPPDSemaforo"
	"fmt"
	"time"
)

const (
	M = 5  // capacidade do pote
	N = 10 // número de selvagens
)

var (
	servings = 0 // porções no pote
	mutex = FPPDSemaforo.NewSemaphore(1) // protege acesso a servings
	emptyPot = FPPDSemaforo.NewSemaphore(0) // sinaliza pote vazio -> cozinheiro
	fullPot = FPPDSemaforo.NewSemaphore(0) // sinaliza pote cheio -> selvagem
)

// encher pote
func putServingsInPot(m int) {
	fmt.Printf("Cozinheiro encheu o pote com %d porções!\n", m)
}

// pegar porção do pote
func getServingFromPot(id int, s int) {
	fmt.Printf("Selvagem %d pegou uma porção. Restam %d no pote.\n", id, s)
}

// comer
func eat(id int) {
	fmt.Printf("Selvagem %d está comendo.\n", id)
}

// cozinheiro 
func cook() {
	for {
		// espera algum selvagem avisar pote vazio
		emptyPot.Wait()
		// enche o pote
		putServingsInPot(M)
		// sinaliza que o pote está cheio
		fullPot.Signal()
	}
}

// selvagem
func savage(id int) {
	for {
		// tenta pegar porção
		mutex.Wait()
		// se vazio 
		if servings == 0 {
			fmt.Printf("\n")
			fmt.Printf("Selvagem %d está esperando.\n", id)
			// avisa cozinheiro que o pote está vazio
			emptyPot.Signal()
			// espera o cozinheiro encher o pote
			fullPot.Wait()
			// atualiza porções no pote
			servings = M
		}
		// pega porção
		servings--
		getServingFromPot(id, servings)
		// libera acesso ao pote
		mutex.Signal()
		// come
		eat(id)
	}
}

func main() {
	// inicia cozinheiro
	go cook()

	// inicia selvagens
	for i := 1; i <= N; i++ {
		go savage(i)
	}

	time.Sleep(10 * time.Second) // tempo de simulação
}