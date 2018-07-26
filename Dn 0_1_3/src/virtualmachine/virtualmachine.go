package virtualmachine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"../parser"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

type Register struct {
	sync.Mutex
	Values map[string]interface{}
}

type VirtualMachine struct {
	stack              []interface{}
	Registers          Register
	Function_Registers Register
	instructions       []string
	Library            map[string]func(*VirtualMachine)
}

func (self *VirtualMachine) GetStack() []interface{} {
	return self.stack
}

func (self *VirtualMachine) Load(instructions []string) {
	self.instructions = instructions
}

func MakeVM(instructions []string) VirtualMachine {
	var empty_arr []interface{}
	return VirtualMachine{empty_arr, Register{Values: make(map[string]interface{})}, Register{Values: make(map[string]interface{})}, instructions, make(map[string]func(*VirtualMachine))}
}

func (self *VirtualMachine) Run() {
	for _, token := range self.instructions {
		if token == "+" {
			self.Op_add()
		} else if token == "-" {
			self.Op_sub()
		} else if token == "*" {
			self.Op_mul()
		} else if token == "/" {
			self.Op_div()
		} else if token == "|" {
			self.Op_print()
		} else if token == "<" {
			self.Op_store(false)
		} else if token == ">" {
			self.Op_load(false)
		} else if token == "$" {
			self.Op_getln()
		} else if token == "!" {
			self.Op_call(true)
		} else if token == "#" {
			self.Op_call(false)
		} else if token == "%" {
			self.Op_real_call()
		} else if token == "&" {
			self.Op_loop(false)
		} else if token == "@" {
			self.Op_create(false)
		} else if token == "." {
			self.Op_read(false)
		} else if token == "," {
			self.Op_write(false)
		} else if token == "^" {
			self.Op_top()
		} else if token == "=" {
			self.Op_eq()
		} else if token == ">>" {
			self.Op_greater()
		} else if token == "<<" {
			self.Op_less()
		} else if token == "\\>" {
			self.Push(">")
		} else if token == "\\<" {
			self.Push("<")
		} else if token == "\\{" {
			self.Push("{")
		} else if token == "\\}" {
			self.Push("}")
		} else if token == "\\|" {
			self.Push("|")
		} else if token == "\\." {
			self.Push(".")
		} else if token == "\\," {
			self.Push(",")
		} else if token == "\\$" {
			self.Push("$")
		} else if token == "\\%" {
			self.Push("%")
		} else if token == "\\^" {
			self.Push("^")
		} else if token == "\\&" {
			self.Push("&")
		} else if token == "\\+" {
			self.Push("+")
		} else if token == "\\-" {
			self.Push("-")
		} else if token == "\\*" {
			self.Push("*")
		} else if token == "\\/" {
			self.Push("/")
		} else if token == "\\!" {
			self.Push("!")
		} else if token == "\\'" {
			self.Push("'")
		} else {
			self.Push(token)
		}
		// self.Function_Registers = make(map[string]interface{})
	}
}

func (self *VirtualMachine) Local_Run() {
	for _, token := range self.instructions {
		if token == "+" {
			self.Op_add()
		} else if token == "-" {
			self.Op_sub()
		} else if token == "*" {
			self.Op_mul()
		} else if token == "/" {
			self.Op_div()
		} else if token == "|" {
			self.Op_print()
		} else if token == "<" {
			self.Op_store(true)
		} else if token == ">" {
			self.Op_load(true)
		} else if token == "$" {
			self.Op_getln()
		} else if token == "!" {
			self.Op_call(true)
		} else if token == "#" {
			self.Op_call(false)
		} else if token == "%" {
			self.Op_real_call()
		} else if token == "&" {
			self.Op_loop(true)
		} else if token == "@" {
			self.Op_create(false)
		} else if token == "." {
			self.Op_read(false)
		} else if token == "," {
			self.Op_write(false)
		} else if token == "^" {
			self.Op_top()
		} else if token == "=" {
			self.Op_eq()
		} else if token == ">>" {
			self.Op_greater()
		} else if token == "<<" {
			self.Op_less()
		} else if token == "\\>" {
			self.Push(">")
		} else if token == "\\<" {
			self.Push("<")
		} else if token == "\\{" {
			self.Push("{")
		} else if token == "\\}" {
			self.Push("}")
		} else if token == "\\|" {
			self.Push("|")
		} else if token == "\\." {
			self.Push(".")
		} else if token == "\\," {
			self.Push(",")
		} else if token == "\\$" {
			self.Push("$")
		} else if token == "\\%" {
			self.Push("%")
		} else if token == "\\^" {
			self.Push("^")
		} else if token == "\\&" {
			self.Push("&")
		} else if token == "\\+" {
			self.Push("+")
		} else if token == "\\-" {
			self.Push("-")
		} else if token == "\\*" {
			self.Push("*")
		} else if token == "\\/" {
			self.Push("/")
		} else if token == "\\!" {
			self.Push("!")
		} else if token == "\\'" {
			self.Push("'")
		} else {
			self.Push(token)
		}
	}
}

func (self *VirtualMachine) Push(item interface{}) {
	switch item.(type) {
	case string:
		if len(item.(string)) > 0 {
			if IsNumeric(item.(string)) && !(item.(string)[0] == '\\') {
				f, _ := strconv.ParseFloat(item.(string), 64)
				self.stack = append(self.stack, f)
			} else {
				if item.(string)[0] == '\\' {
					self.stack = append(self.stack, item.(string)[1:])
				} else {
					self.stack = append(self.stack, item)
				}
			}
		} else {
			self.stack = append(self.stack, "")
		}
	case int:
		self.stack = append(self.stack, float64(item.(int)))
	default:
		self.stack = append(self.stack, item)
	}
}

func (self *VirtualMachine) Pop() interface{} {
	var Pop interface{}

	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to Pop from stack when length of stack is less than 1"))
		os.Exit(1)
	}

	Pop, self.stack = self.stack[len(self.stack)-1], self.stack[:len(self.stack)-1]
	return Pop
}

func (self *VirtualMachine) Op_add() {

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to add when length of stack is less than 2"))
		os.Exit(1)
	}

	operand1 := self.Pop()
	operand2 := self.Pop()
	switch operand1.(type) {
	case int:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(int) + operand2.(int))
		case float64:
			self.Push(float64(operand1.(int)) + operand2.(float64))
		default:
			fmt.Println(errors.New("ERROR: Cannot add Non-Number to Number"))
			os.Exit(1)
		}
	case float64:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(float64) + float64(operand2.(int)))
		case float64:
			self.Push(operand1.(float64) + operand2.(float64))
		default:
			fmt.Println(errors.New("ERROR: Cannot add Non-Number to Number"))
			os.Exit(1)
		}
	case string:
		switch operand2.(type) {
		case string:
			self.Push(operand1.(string) + operand2.(string))
		default:
			fmt.Println(errors.New("ERROR: Cannot add Non-String to String"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot add Non-Number or Non-Strings"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_sub() {

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to subtract when length of stack is less than 2"))
		os.Exit(1)
	}

	operand1 := self.Pop()
	operand2 := self.Pop()
	switch operand1.(type) {
	case int:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(int) - operand2.(int))
		case float64:
			self.Push(float64(operand1.(int)) - operand2.(float64))
		default:
			fmt.Println(errors.New("ERROR: Cannot subtract Non-Number from Number"))
			os.Exit(1)
		}
	case float64:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(float64) - float64(operand2.(int)))
		case float64:
			self.Push(operand1.(float64) - operand2.(float64))
		default:
			fmt.Println(errors.New("ERROR: Cannot subtract Non-Number from Number"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot subtract Non-Number"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_mul() {

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to multiply when length of stack is less than 2"))
		os.Exit(1)
	}

	operand1 := self.Pop()
	operand2 := self.Pop()
	switch operand1.(type) {
	case int:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(int) * operand2.(int))
		case string:
			self.Push(strings.Repeat(operand2.(string), operand1.(int)))
		default:
			fmt.Println(errors.New("ERROR: Cannot multiply Number by Non-Number or Non-String"))
			os.Exit(1)
		}
	case float64:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(float64) * float64(operand2.(int)))
		case float64:
			self.Push(operand1.(float64) * operand2.(float64))
		case string:
			self.Push(strings.Repeat(operand2.(string), int(operand1.(float64))))
		default:
			fmt.Println(errors.New("ERROR: Cannot multiply Number by Non-Number or Non-String"))
			os.Exit(1)
		}
	case string:
		switch operand2.(type) {
		case int:
			self.Push(strings.Repeat(operand1.(string), operand2.(int)))
		case float64:
			self.Push(strings.Repeat(operand1.(string), int(operand2.(float64))))
		default:
			fmt.Println(errors.New("ERROR: Cannot multiply String by Non-Number"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot multiply Non-Number or Non-String"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_div() {

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to divide when length of stack is less than 2"))
		os.Exit(1)
	}

	operand1 := self.Pop()
	operand2 := self.Pop()
	switch operand1.(type) {
	case int:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(int) / operand2.(int))
		default:
			fmt.Println(errors.New("ERROR: Cannot divide Number by Non-Number"))
			os.Exit(1)
		}
	case float64:
		switch operand2.(type) {
		case int:
			self.Push(operand1.(float64) / float64(operand2.(int)))
		case float64:
			self.Push(operand1.(float64) / operand2.(float64))
		default:
			fmt.Println(errors.New("ERROR: Cannot divide Number by Non-Number"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot divide Non-Numbers"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_store(local bool) {

	self.Registers.Lock()
	self.Function_Registers.Lock()
	defer self.Registers.Unlock()
	defer self.Function_Registers.Unlock()

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to store to Identifier when length of stack is less than 2"))
		os.Exit(1)
	}

	operand1 := self.Pop()
	operand2 := self.Pop()

	switch operand1.(type) {
	case string:
		if local {
			self.Function_Registers.Values[operand1.(string)] = operand2
		} else {
			self.Registers.Values[operand1.(string)] = operand2
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot store to Non-Identifier"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_load(local bool) {

	self.Registers.Lock()
	defer self.Registers.Unlock()
	self.Function_Registers.Lock()
	defer self.Function_Registers.Unlock()

	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to load from nothing"))
		os.Exit(1)
	}

	operand := self.Pop()
	// fmt.Print("Loading: ")
	// fmt.Println(operand)

	switch operand.(type) {
	case string:
		if local {
			if _, ok := self.Function_Registers.Values[operand.(string)]; ok {
				self.Push(self.Function_Registers.Values[operand.(string)])
			} else {
				self.Push(self.Registers.Values[operand.(string)])
			}
		} else {
			self.Push(self.Registers.Values[operand.(string)])
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot load from Non-Identifier"))
		os.Exit(1)
	}

	// if local {
	// 	if _, ok := self.Function_Registers[operand.(string)]; ok {
	// 		self.Push(self.Function_Registers[operand.(string)])
	// 	} else {
	// 		self.Push(self.Registers[operand.(string)])
	// 	}
	// } else {
	// 	self.Push(self.Registers[operand.(string)])
	// }
}

func (self *VirtualMachine) Op_call(local bool) {

	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to call a function from nothing"))
		os.Exit(1)
	}

	function := self.Pop()

	switch function.(type) {
	case string:
		parser := parser.MakeParser(function.(string)[1 : len(function.(string))-1])
		self.Load(parser.Parse())
		if local {
			self.Local_Run()
		} else {
			self.Run()
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot call a Non-Function"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_top() {

	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to duplicate nothing"))
		os.Exit(1)
	}

	operand := self.Pop()
	self.Push(operand)
	self.Push(operand)
}

func (self *VirtualMachine) Op_ws() {
	self.Push(" ")
}

func (self *VirtualMachine) Op_loop(local bool) {
	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to loop when length of stack is less than 2"))
		os.Exit(1)
	}

	condition := self.Pop()
	function := self.Pop()
	var value interface{}
	for {
		self.Push(condition)
		self.Op_call(local)
		value = self.Pop()

		if value.(float64) > 0 {
			self.Push(function)
			self.Op_call(local)
		} else {
			break
		}
	}
}

func (self *VirtualMachine) Op_read(local bool) {

	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to store to Identifier when length of stack is less than 2"))
		os.Exit(1)
	}

	var index interface{} = self.Pop()
	var this interface{} = self.Pop()

	switch index.(type) {
	case string:
		_, ok := this.(map[string]interface{})[index.(string)]
		if ok {
			self.Push(this.(map[string]interface{})[index.(string)])
		}
	default:
		fmt.Println(errors.New("ERROR: Tried to access member that was not in object"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_write(local bool) {

	if len(self.stack) < 3 {
		fmt.Println(errors.New("ERROR: Tried to write to object when length of stack is less than 3"))
		os.Exit(1)
	}

	var index interface{} = self.Pop()
	var this interface{} = self.Pop()
	var value interface{} = self.Pop()

	switch index.(type) {
	case string:
		switch this.(type) {
		case map[string]interface{}:
			this.(map[string]interface{})[index.(string)] = value
		default:
			fmt.Println(errors.New("ERROR: Tried to write to member of an object from nothing"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Tried to write to member that was not in object"))
		os.Exit(1)
	}
	self.Push(this)
}

func (self *VirtualMachine) Op_create(local bool) {
	self.Push(make(map[string]interface{}))
}

// def create(self, f=False):
//     self.Push({})

func (self *VirtualMachine) Op_getln() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	self.Push(text[:len(text)-2])
}

func (self *VirtualMachine) Op_greater() {
	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to compare two numbers with greater when length of stack is less than 2"))
		os.Exit(1)
	}

	var operand1 interface{} = self.Pop()
	var operand2 interface{} = self.Pop()

	switch operand1.(type) {
	case float64:
		switch operand2.(type) {
		case float64:
			if operand1.(float64) > operand2.(float64) {
				self.Push("1")
			} else {
				self.Push("0")
			}
		default:
			fmt.Println(errors.New("ERROR: Tried to compare a Number and a Non-Number with greater"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Tried to compare a Number and a Non-Number with greater"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_less() {
	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to compare two numbers with less when length of stack is less than 2"))
		os.Exit(1)
	}

	var operand1 interface{} = self.Pop()
	var operand2 interface{} = self.Pop()

	switch operand1.(type) {
	case float64:
		switch operand2.(type) {
		case float64:
			if operand1.(float64) < operand2.(float64) {
				self.Push("1")
			} else {
				self.Push("0")
			}
		default:
			fmt.Println(errors.New("ERROR: Tried to compare a Number and a Non-Number with less"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Tried to compare a Number and a Non-Number with less"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Op_eq() {
	if len(self.stack) < 2 {
		fmt.Println(errors.New("ERROR: Tried to compare two numbers with equals when length of stack is less than 2"))
		os.Exit(1)
	}

	var operand1 interface{} = self.Pop()
	var operand2 interface{} = self.Pop()

	switch operand1.(type) {
	case map[string]interface{}:
		eq := reflect.DeepEqual(operand1, operand2)
		if eq {
			self.Push("1")
		} else {
			self.Push("0")
		}
	default:
		if operand1 == operand2 {
			self.Push("1")
		} else {
			self.Push("0")
		}
	}
}

func (self *VirtualMachine) Op_print() {
	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to print when length of stack is less than 1"))
		os.Exit(1)
	}
	// fmt.Print("Printing...")
	operand := self.Pop()
	fmt.Print(operand)
}

func (self *VirtualMachine) Op_real_call() {
	if len(self.stack) < 1 {
		fmt.Println(errors.New("ERROR: Tried to call low level function when length of stack is less than 1"))
		os.Exit(1)
	}
	f := self.Pop()

	switch f.(type) {
	case string:
		self.Library[f.(string)](self)
	default:
		fmt.Println(errors.New("ERROR: Tried to call low level function from nothing"))
		os.Exit(1)
	}
}

func (self *VirtualMachine) Length() int {
	return len(self.stack)
}
