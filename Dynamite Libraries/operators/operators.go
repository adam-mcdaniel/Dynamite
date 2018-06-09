package operators

import (
	"errors"
	"fmt"
	"os"

	. "../virtualmachine"
)

func Not(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call Not on number when length of stack is less than 1"))
		os.Exit(1)
	}
	operand := vm.Pop().(float64)
	if operand >= 1 {
		vm.Push(float64(0))
	} else {
		vm.Push(float64(1))
	}
}

func And(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call And on two numbers when length of stack is less than 2"))
		os.Exit(1)
	}
	operand1 := vm.Pop().(float64)
	operand2 := vm.Pop().(float64)
	if operand1 == 1 && operand2 == 1 {
		vm.Push(float64(1))
	} else {
		vm.Push(float64(0))
	}
}

func Or(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call Or on two numbers when length of stack is less than 2"))
		os.Exit(1)
	}
	operand1 := vm.Pop().(float64)
	operand2 := vm.Pop().(float64)
	if operand1 == 1 || operand2 == 1 {
		vm.Push(float64(1))
	} else {
		vm.Push(float64(0))
	}
}

func Neg(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call Neg on a number when length of stack is less than 1"))
		os.Exit(1)
	}
	operand := vm.Pop().(float64)
	vm.Push(-float64(operand))
}

func Method(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call Method when length of stack is less than 2"))
		os.Exit(1)
	}
	object := vm.Pop()
	member := vm.Pop()
	vm.Push(object)
	vm.Push(member)
	vm.Op_read(true)
	vm.Op_call(true)
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["operators.not"] = Not
	vm.Push("{ operators.not % }")
	vm.Push("not")
	vm.Op_store(false)
	vm.Library["operators.and"] = And
	vm.Push("{ operators.and % }")
	vm.Push("and")
	vm.Op_store(false)
	vm.Library["operators.or"] = Or
	vm.Push("{ operators.or % }")
	vm.Push("or")
	vm.Op_store(false)
	vm.Library["operators.neg"] = Neg
	vm.Push("{ operators.neg % }")
	vm.Push("neg")
	vm.Op_store(false)
	vm.Library["operators.method"] = Method
	vm.Push("{ operators.method % }")
	vm.Push("method")
	vm.Op_store(false)
}
