//GRUPO:
//Ana Carolina Poletto
//Anthony Antonelli
//Barbara Dapper
//Diego de Mello
//Louise Northfleet

//-------------------------------------------------------------------------
//PROBLEMA:
// Stand Claus sleeps in his shop at the North Pole and can only be
// awakened by either (1) all nine reindeer being back from their vaca-
// tion in the South Pacific, or (2) some of the elves having diﬃculty
// making toys; to allow Santa to get some sleep, the elves can only
// wake him when three of them have problems. When three elves are
// having their problems solved, any other elves wishing to visit Santa
// must wait for those elves to return. If Santa wakes up to find three
// elves waiting at his shop’s door, along with the last reindeer having
// come back from the tropics, Santa has decided that the elves can
// wait until after Christmas, because it is more important to get his
// sleigh ready. (It is assumed that the reindeer do not want to leave
// the tropics, and therefore they stay there until the last possible mo-
// ment.) The last reindeer to arrive must get Santa while the others
// wait in a warming hut before being harnessed to the sleigh.
// Here are some addition specifications:
// After the ninth reindeer arrives, Santa must invoke prepareSleigh, and
// then all nine reindeer must invoke getHitched.
// After the third elf arrives, Santa must invoke helpElves. Concurrently,
// all three elves should invoke getHelp.
// All three elves must invoke getHelp before any additional elves enter
// (increment the elf counter).
// Santa should run in a loop so he can help many sets of elves. We can assume
// that there are exactly 9 reindeer, but there may be any number of elves.

// -------------------------------------------------------------------------
// CÓDIGO COM SEMAFOROS:

package main

import (
	"T2-FPPD/FPPDSemaforo"
	"fmt"
	"time"
)

var elves = 0 // número de elfos
var reindeer = 0 // número de renas
var mutex = FPPDSemaforo.NewSemaphore(1) // protege acesso a counters
var santaSem = FPPDSemaforo.NewSemaphore(0) // acorda o Papai Noel
var reindeerSem = FPPDSemaforo.NewSemaphore(0) // renas esperando para serem presas ao trenó
var elfTex = FPPDSemaforo.NewSemaphore(1) // limita número de elfos que podem pedir ajuda

// elfo chega
func elfArrives() {
	// tenta pedir ajuda
	elfTex.Wait()
	// pede ajuda
	mutex.Wait()
	elves++

	fmt.Println("Elfo: Pede Ajuda", elves)

	// se for o terceiro elfo, acorda o Papai Noel
	if elves == 3 {
		santaSem.Signal()

	} else {
		// se não for o terceiro, libera o semáforo para outro elfo
		elfTex.Signal()

	}
	mutex.Signal()

	// espera ajuda do Papai Noel
	getHelp()

	// elfo foi ajudado, sai
	mutex.Wait()
	elves--
	// se for o último elfo a sair, libera o semáforo para outros elfos
	if elves == 0 {
		elfTex.Signal()
	}
	mutex.Signal()
}

// rena chega
func reindeerArrives() {
	// tenta entrar
	mutex.Wait()
	reindeer++
	fmt.Println("Rena: Volta das férias", reindeer)

	// se for a nona rena, acorda o Papai Noel
	if reindeer == 9 {
		for i := 0; i < 9; i++ {
			santaSem.Signal()
		}
	}
	mutex.Signal()

	// espera ser presa ao trenó
	reindeerSem.Wait()
	getHitched()
}

// Papai Noel prepara o trenó
func prepareSleigh() {
	fmt.Println("Santa: Preparando o trenó")
}

// renas são presas ao trenó
func getHitched() {
	fmt.Println("Santa: Preparando as renas")
	reindeerSem.Signal()
}

// Papai Noel ajuda os elfos
func helpElves() {
	fmt.Println("Santa: Ajudando os elfos", elves)
	time.Sleep(1 * time.Second)
}

// elfos recebem ajuda
func getHelp() {
	fmt.Println("Elfos: Esperando ajuda", elves)
	time.Sleep(1 * time.Second)
}


func main() {
	// inicia goroutines para a chegada dos elfos
	for i := 1; i <= 1000; i++ {
		go elfArrives()
	}

	// inicia goroutines para a chegada das renas
	for i := 1; i <= 9; i++ {
		go reindeerArrives()
	}

	//loop principal onde o Papai Noel verifica se precisa ajudar os elfos ou preparar o trenó
	for {
		// espera ser acordado por elfos ou renas
		santaSem.Wait()
		mutex.Wait()

		// se todas as renas chegaram, prepara o trenó
		if reindeer == 9 {
			prepareSleigh()
			reindeerSem.Signal()
			fmt.Println("É Natal, entregando presentes!")
			return
		} else if elves == 3 {
			// se três elfos estão esperando, ajuda os elfos
			helpElves()
		}

		mutex.Signal()
	}
}