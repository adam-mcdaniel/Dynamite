package threading

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"../parser"
	. "../virtualmachine"
)

type Channel struct {
	sync.Mutex
	list []interface{}
}

var channel *Channel

func GetTopChannel(vm *VirtualMachine) {
	channel.Lock()
	defer channel.Unlock()
	if len(channel.list) > 0 {
		r := channel.list[len(channel.list)-1]
		if r != nil {
			vm.Push(r)
		} else {
			vm.Push(0)
		}
	} else {
		vm.Push(0)
	}
}

func GetBottomChannel(vm *VirtualMachine) {
	channel.Lock()
	defer channel.Unlock()
	if len(channel.list) > 0 {
		r := channel.list[0]
		if r != nil {
			vm.Push(r)
		} else {
			vm.Push(0)
		}
	} else {
		vm.Push(0)
	}
}

func PushToChannel(vm *VirtualMachine) {
	channel.Lock()
	defer channel.Unlock()
	if vm.Length() > 0 {
		channel.list = append(channel.list, vm.Pop())
	}
}

func DeleteFrom(vm *VirtualMachine) {
	channel.Lock()
	defer channel.Unlock()
	if len(channel.list) > 0 {
		channel.list = channel.list[1:]
	}
}

// func Stall(vm *VirtualMachine) {
// 	for {

// 	}
// }

func MakeThread(vm *VirtualMachine) {
	function := vm.Pop()
	switch function.(type) {
	case string:
		parser := parser.MakeParser("'\n' endl < " +
			"{*} mul < " +
			"{/} div < " +
			"{+} add < " +
			"{-} sub < " +
			"{>>} greater < " +
			"{<<} less < " +
			"{=} equal < " +
			"{$} input < " +
			"{|} print < " +
			"{<} __unpack_func__ < " +
			"{__unpack_func__ > #} unpack < " +
			"{x < ':' x > + ':' + |} print_lit < " +
			"{func < 1 func > { } &} while < " +
			"{condition < function < {function > ! 0 condition <} {condition >} &} if < " +
			"{ @ self < lists.List > ! self >  list , self <  { self <  a <  a >  self >  list . lists.Append > ! self >  list , self <  self > } append < append >  self >  append , self <  { self <  self >  list . lists.Pop > ! self >  list , self <  self > } pop < pop >  self >  pop , self <  { self <  self >  list . lists.Length > ! } len < len >  self >  len , self <  { self <  self >  list . lists.Items > ! } items < items >  self >  items , self <  { self <  n <  n >  self >  list . lists.Index > ! } index < index >  self >  index , self <  self > } list < " +
			"{equal > ! not > !} notequal < " +
			"{ s <  s >  print > ! endl >  print > ! } println < " +
			"{ case <  if_then <  else_then <   {  if_then > !  }  case >  if > !  {  else_then > !  }  case >  not > ! if > ! } ifelse < " +
			function.(string) + " ! ")

		tokens := parser.Parse()
		NewVM := MakeVM(tokens)

		for key, value := range vm.Library {
			NewVM.Library[key] = value
		}

		for key, value := range vm.Registers.Values {
			NewVM.Registers.Values[key] = value
		}

		// NewVM.Library = vm.Library
		// NewVM.Registers = vm.Registers

		go NewVM.Run()

	default:
		fmt.Println(errors.New("ERROR: MakeThread argument must be a function"))
		os.Exit(1)
	}
}

func Help(vm *VirtualMachine) {
	fmt.Println(`
function threading.Thread
{
	arg1: func function
	This function runs a function in a new thread
}

function threading.GetTopChannel
{
	This function returns the last thing added to the channel
}

function threading.GetBottomChannel
{
	This function returns the bottom element of the channel
}

function threading.DeleteFrom
{
	This function deletes the bottom element of the channel
}

function threading.PushToChannel
{
	This function adds to the top of the channel
}
`)
}

func InstallLibrary(vm *VirtualMachine) {
	channel = &Channel{}
	// locked = "idle"
	vm.Library["threading.makethread"] = MakeThread
	vm.Push("{ threading.makethread % }")
	vm.Push("threading.Thread")
	vm.Op_store(false)
	vm.Library["threading.get_top_from_channel"] = GetTopChannel
	vm.Push("{ threading.get_top_from_channel % }")
	vm.Push("threading.GetTopChannel")
	vm.Op_store(false)
	vm.Library["threading.get_bottom_from_channel"] = GetBottomChannel
	vm.Push("{ threading.get_bottom_from_channel % }")
	vm.Push("threading.GetBottomChannel")
	vm.Op_store(false)
	vm.Library["threading.push_to_channel"] = PushToChannel
	vm.Push("{ threading.push_to_channel % }")
	vm.Push("threading.PushToChannel")
	vm.Op_store(false)
	vm.Library["threading.delete_from_channel"] = DeleteFrom
	vm.Push("{ threading.delete_from_channel % }")
	vm.Push("threading.DeleteFromChannel")
	vm.Op_store(false)
	vm.Library["threading.help_function"] = Help
	vm.Push("{ threading.help_function % }")
	vm.Push("threading.help")
	vm.Op_store(false)
}
