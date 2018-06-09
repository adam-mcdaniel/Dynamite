package io

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	. "../virtualmachine"
)

var (
	result string
	err    error
	in     *bufio.Reader
)

func getInput(input chan string) {
	in = bufio.NewReader(os.Stdin)
	result, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input <- result[:len(result)-2]
}

func inputTimeout(seconds float64) string {
	for {
		input := make(chan string, 1)
		go getInput(input)

		select {
		case i := <-input:
			return i
		case <-time.After(time.Duration(seconds) * time.Second):
			in = bufio.NewReader(os.Stdin)
			return ""
		}
	}
}

func InputTimeoutWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call Listen when length of stack is less than 1"))
		os.Exit(1)
	}

	seconds := vm.Pop()

	switch seconds.(type) {
	case float64:
		vm.Push(inputTimeout(seconds.(float64)))
	default:
		fmt.Println(errors.New("ERROR: Seconds in input_timeout is a String"))
		os.Exit(1)
	}
}

func ReadFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}

	str := string(b) // convert content to a 'string'
	return str
}

func WriteFile(contents, path string) {
	file, _ := os.Create(path)
	fmt.Fprintf(file, contents)
	file.Close()
}

func ReadFileWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call ReadFile when length of stack is less than 1"))
		os.Exit(1)
	}

	path := vm.Pop()

	switch path.(type) {
	case string:
		vm.Push(ReadFile(path.(string)))
	default:
		fmt.Println(errors.New("ERROR: Path in ReadFile is a String"))
		os.Exit(1)
	}
}

func WriteFileWrapper(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call WriteFile when length of stack is less than 2"))
		os.Exit(1)
	}

	contents := vm.Pop()
	path := vm.Pop()

	switch path.(type) {
	case string:
		switch contents.(type) {
		case string:
			WriteFile(contents.(string), path.(string))
		default:
			fmt.Println(errors.New("ERROR: Contents in WriteFile is a String"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Path in WriteFile is a String"))
		os.Exit(1)
	}
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["io.input_timeout"] = InputTimeoutWrapper
	vm.Push("{ io.input_timeout % }")
	vm.Push("io.input_timeout")
	vm.Op_store(false)
	vm.Library["io.writefile"] = WriteFileWrapper
	vm.Push("{ io.writefile % }")
	vm.Push("io.WriteFile")
	vm.Op_store(false)
	vm.Library["io.readfile"] = ReadFileWrapper
	vm.Push("{ io.readfile % }")
	vm.Push("io.ReadFile")
	vm.Op_store(false)
}
