package net

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"

	. "../virtualmachine"
)

var conn net.Conn
var conns map[float64]net.Conn

func ServerSend(newmessage string, connection float64) {
	conns[connection].Write([]byte(newmessage + " \n"))
}

func ServerSendWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call ServerSend when length of stack is less than 1"))
		os.Exit(1)
	}
	message := vm.Pop()
	connection := vm.Pop()

	switch message.(type) {
	case string:
		switch connection.(type) {
		case float64:
			ServerSend(message.(string), connection.(float64))
		default:
			fmt.Println(errors.New("ERROR: Connection number must be Number when using ServerSend"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot send Non-string data from server"))
		os.Exit(1)
	}
}

func ClientSend(message string) {
	fmt.Fprintf(conn, message+" \n")
}

func ClientSendWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call ClientSend when length of stack is less than 1"))
		os.Exit(1)
	}
	message := vm.Pop()

	switch message.(type) {
	case string:
		ClientSend(message.(string))
	default:
		fmt.Println(errors.New("ERROR: Cannot send Non-string data from client"))
		os.Exit(1)
	}
}

func ClientReceive() string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	return message
}

func ClientReceiveWrapper(vm *VirtualMachine) {
	vm.Push(ClientReceive())
}

func ServerReceive(connection float64) string {
	message, _ := bufio.NewReader(conns[connection]).ReadString('\n')
	return message
}

func ServerReceiveWrapper(vm *VirtualMachine) {
	connection := vm.Pop()

	switch connection.(type) {
	case float64:
		vm.Push(ServerReceive(connection.(float64)))
	default:
		fmt.Println(errors.New("ERROR: Connecion number must be a Number when using ServerReceive"))
		os.Exit(1)
	}
}

func Connect(ip, port string) {
	conn, _ = net.Dial("tcp", ip+":"+port)
}

func ConnectWrapper(vm *VirtualMachine) {
	if vm.Length() < 2 {
		fmt.Println(errors.New("ERROR: Tried to call Connect when length of stack is less than 2"))
		os.Exit(1)
	}
	ip := vm.Pop()
	port := vm.Pop()

	switch ip.(type) {
	case string:
		switch port.(type) {
		case string:
			Connect(ip.(string), port.(string))
		default:
			fmt.Println(errors.New("ERROR: The connection Port must be a string"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: The IP must be a string"))
		os.Exit(1)
	}
}

func Listen(port string, connection float64) {
	ln, _ := net.Listen("tcp", ":"+port)
	conns[connection], _ = ln.Accept()
}

func ListenWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call Listen when length of stack is less than 1"))
		os.Exit(1)
	}

	port := vm.Pop()
	connection := vm.Pop()

	switch port.(type) {
	case string:
		switch connection.(type) {
		case float64:
			Listen(port.(string), connection.(float64))
		default:
			fmt.Println(errors.New("ERROR: The connection number must be a Number"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: The listening Port must be a string"))
		os.Exit(1)
	}
}

// func main() {
// 	Connect("127.0.0.1", "8080")
// 	for {
// 		text := "Hello world!"
// 		ClientSend(text)
// 		message := Receive()
// 		fmt.Print("Message from server: " + message)
// 	}
// }

// func main() {
// 	Listen("8080")
// 	for {
// 		message := Receive()
// 		fmt.Print("Message Received:", string(message))
// 		newmessage := strings.ToUpper(message)
// 		ServerSend(newmessage)
// 	}
// }

func InstallLibrary(vm *VirtualMachine) {
	conns = make(map[float64]net.Conn)
	vm.Library["net.listen"] = ListenWrapper
	vm.Push("{ net.listen % }")
	vm.Push("net.Listen")
	vm.Op_store(false)
	vm.Library["net.connect"] = ConnectWrapper
	vm.Push("{ net.connect % }")
	vm.Push("net.Connect")
	vm.Op_store(false)
	vm.Library["net.serversend"] = ServerSendWrapper
	vm.Push("{ net.serversend % }")
	vm.Push("net.ServerSend")
	vm.Op_store(false)
	vm.Library["net.clientsend"] = ClientSendWrapper
	vm.Push("{ net.clientsend % }")
	vm.Push("net.ClientSend")
	vm.Op_store(false)
	vm.Library["net.serverreceive"] = ServerReceiveWrapper
	vm.Push("{ net.serverreceive % }")
	vm.Push("net.ServerReceive")
	vm.Op_store(false)
	vm.Library["net.clientreceive"] = ClientReceiveWrapper
	vm.Push("{ net.clientreceive % }")
	vm.Push("net.ClientReceive")
	vm.Op_store(false)
}
