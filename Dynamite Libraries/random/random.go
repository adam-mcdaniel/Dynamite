package random

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	. "../virtualmachine"
)

func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomVM(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call Randint on between numbers when length of stack is less than 2"))
		os.Exit(1)
	}
	low := vm.Pop().(float64)
	high := vm.Pop().(float64)
	vm.Push(float64(Random(int(low), int(high)+1)))
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["random.randomInteger"] = RandomVM
	vm.Push("{ random.randomInteger % }")
	vm.Push("random.randint")
	vm.Op_store(false)
}
