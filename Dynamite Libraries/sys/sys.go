package sys

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "../virtualmachine"
)

func System(cmd_string, args string) {
	s := strings.Split(args, " ")
	cmd := exec.Command(cmd_string, s...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func SystemWrapper(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call System when length of stack is less than 2"))
		os.Exit(1)
	}

	cmd := vm.Pop()
	args := vm.Pop()

	switch cmd.(type) {
	case string:
		switch args.(type) {
		case string:
			System(cmd.(string), args.(string))
		default:
			fmt.Println(errors.New("ERROR: Args in system is a string"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Command in System is a string"))
		os.Exit(1)
	}
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["os.system_wrapper"] = SystemWrapper
	vm.Push("{ os.system_wrapper % }")
	vm.Push("os.System")
	vm.Op_store(false)
}
