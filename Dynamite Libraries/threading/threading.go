package threading

import (
	. "../virtualmachine"
)

func MakeThread(vm *VirtualMachine) {
	go vm.Op_call(true)
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["threading.makethread"] = MakeThread
	vm.Push("{ threading.makethread % }")
	vm.Push("threading.Thread")
	vm.Op_store(false)
}
