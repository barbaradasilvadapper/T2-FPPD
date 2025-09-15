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

var (
	elves     = 0
	reindeer  = 0
	santaSem  = FPPDSemaforo.NewSemaphore(0)
	reinSem   = FPPDSemaforo.NewSemaphore(0)
	elfTex    = FPPDSemaforo.NewSemaphore(1)
	mutex     = FPPDSemaforo.NewSemaphore(1)
	totalRein = 9
	elfSem    = FPPDSemaforo.NewSemaphore(0)
	santaBusy = FPPDSemaforo.NewSemaphore(1)
	elvesDone = FPPDSemaforo.NewSemaphore(0)
	reinDone = FPPDSemaforo.NewSemaphore(0)

)

func prepareSleigh() {
	fmt.Println("Papai Noel: Preparando o trenó!")
}

func helpElves() {
	fmt.Println("Papai Noel: Ajudando 3 elfos...")
}

func getHitched(id int) {
	fmt.Printf("Rena %d foi atrelada ao trenó!\n", id)
}

func getHelp(id int) {
	fmt.Printf("Elfo %d recebeu ajuda!\n", id)
}

// Santa loop
func santa() {
	for {
		santaSem.Wait()
		santaBusy.Wait()
		mutex.Wait()
		if reindeer == totalRein {
			prepareSleigh()
			for i := 0; i < totalRein; i++ {
				reinSem.Signal()
			}
			for i := 0; i < totalRein; i++ {
				reinDone.Wait()
			}
			reindeer = 0
			mutex.Signal()
		} else if elves == 3 {
			helpElves()
			// Libera os 3 elfos após ajudar
			for i := 0; i < 3; i++ {
				elfSem.Signal()
			}
			// espera os 3 elfos terminarem
			for i := 0; i < 3; i++ {
				elvesDone.Wait()
			}
			mutex.Signal()
		} else {
			mutex.Signal()
		}
		santaBusy.Signal()
	}
}

// Reindeer process
func reindeerFunc(id int) {
	mutex.Wait()
	reindeer++
	if reindeer == totalRein {
		santaSem.Signal()
	}
	mutex.Signal()

	reinSem.Wait()
	getHitched(id)
	reinDone.Signal()

}

// Elf process
func elfFunc(id int) {
	elfTex.Wait()
	mutex.Wait()
	elves++
	if elves == 3 {
		santaSem.Signal()
	} else {
		elfTex.Signal()
	}
	mutex.Signal()

	elfSem.Wait()
	getHelp(id)
	elvesDone.Signal()


	mutex.Wait()
	elves--
	if elves == 0 {
		elfTex.Signal()
	}
	mutex.Signal()
}

func main() {
	go santa()

	// Renas
	for i := 1; i <= totalRein; i++ {
		go reindeerFunc(i)
	}

	// Elfos (pode ter muitos)
	for i := 1; i <= 10; i++ {
		go elfFunc(i)
	}

	time.Sleep(5 * time.Second) // tempo para simulação
}
