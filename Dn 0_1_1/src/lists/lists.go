package lists

import (
	"errors"
	"fmt"
	"os"

	. "../virtualmachine"
	// "strconv"
	"strings"
)

func MakeList(vm *VirtualMachine) {
	vm.Push("{@ self < a self > counter , self < self >}")
	vm.Op_call(true)
}

func AddToList(vm *VirtualMachine) {
	vm.Push("{ self < data < data > self > self > counter . , self < self > counter . 'a' + self > counter , }")
	vm.Op_call(true)
}

func Pop(vm *VirtualMachine) {
	self := vm.Pop().(map[string]interface{})
	vm.Push(self)
	Length(vm)
	length := vm.Pop().(int)
	_, ok := self[strings.Repeat("a", length)]
	if ok {
		delete(self, strings.Repeat("a", length))
		self["counter"] = strings.Repeat("a", length)
	}
	vm.Push(self)
}

func Length(vm *VirtualMachine) {
	vm.Push("{self < self > counter .}")
	vm.Op_call(true)
	len := len(vm.Pop().(string))
	vm.Push(len - 1)
}

func Items(vm *VirtualMachine) {
	self := vm.Pop().(map[string]interface{})
	vm.Push(self)
	Length(vm)
	limit := vm.Pop().(int)
	for i := 1; i < limit+1; i++ {
		vm.Push(self)
		vm.Push("{ self < self > " + strings.Repeat("a", i) + " . }")
		vm.Op_call(true)
	}
}

func Index(vm *VirtualMachine) {
	self := vm.Pop().(map[string]interface{})
	indice := vm.Pop().(float64) + 1
	vm.Push(self)
	Length(vm)
	limit := vm.Pop().(int)
	for i := 1; i < limit+1; i++ {
		if float64(i) == float64(indice) {
			vm.Push(self)
			vm.Push("{ self < self > " + strings.Repeat("a", i) + " . }")
			vm.Op_call(true)
		}
	}
}

func For(vm *VirtualMachine) {
	if vm.Length() < 3 {
		fmt.Println(errors.New("ERROR: Tried to do For loop when length of stack is less than 3"))
		os.Exit(1)
	}

	value := vm.Pop()
	list := vm.Pop()
	function := vm.Pop()
	for {
		vm.Push(list)
		Length(vm)
		length := vm.Pop().(int)
		if length > 0 {
			vm.Push(float64(length - 1))
			vm.Push(list)
			Index(vm)
			vm.Push(value)
			vm.Op_store(true)
			vm.Push(function)
			vm.Op_call(true)
			vm.Push(list)
			Pop(vm)
			list = vm.Pop()
		} else {
			break
		}
	}
}

func Range(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to create Range when length of stack is less than 2"))
		os.Exit(1)
	}
	low := vm.Pop().(float64)
	high := vm.Pop().(float64)
	MakeList(vm)
	list := vm.Pop()
	for i := high; i > low-1; i-- {
		vm.Push(float64(i))
		vm.Push(list)
		AddToList(vm)
		list = vm.Pop()
	}
	vm.Push(list)
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["operators.list"] = MakeList
	vm.Push("{ operators.list % }")
	vm.Push("lists.List")
	vm.Op_store(false)
	vm.Library["operators.append"] = AddToList
	vm.Push("{ operators.append % }")
	vm.Push("lists.Append")
	vm.Op_store(false)
	vm.Library["operators.len"] = Length
	vm.Push("{ operators.len % }")
	vm.Push("lists.Length")
	vm.Op_store(false)
	vm.Library["operators.items"] = Items
	vm.Push("{ operators.items % }")
	vm.Push("lists.Items")
	vm.Op_store(false)
	vm.Library["operators.pop"] = Pop
	vm.Push("{ operators.pop % }")
	vm.Push("lists.Pop")
	vm.Op_store(false)
	vm.Library["operators.index"] = Index
	vm.Push("{ operators.index % }")
	vm.Push("lists.Index")
	vm.Op_store(false)
	vm.Library["operators.for"] = For
	vm.Push("{ operators.for % }")
	vm.Push("for")
	vm.Op_store(false)
	vm.Library["operators.range"] = Range
	vm.Push("{ operators.range % }")
	vm.Push("range")
	vm.Op_store(false)
}
