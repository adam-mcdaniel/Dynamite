package udp

import (
	"errors"
	"fmt"
	"net"
	"os"

	. "../virtualmachine"
)

var addrs map[string]*net.UDPAddr
var serverAddr *net.UDPAddr
var serverConn *net.UDPConn
var buffer []byte

func Listen(port string) {
	serverAddr, _ = net.ResolveUDPAddr("udp", ":"+port)
	serverConn, _ = net.ListenUDP("udp", serverAddr)
}

func ListenWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call Listen when length of stack is less than 1"))
		os.Exit(1)
	}

	port := vm.Pop()

	switch port.(type) {
	case string:
		Listen(port.(string))
	default:
		fmt.Println(errors.New("ERROR: The listening Port must be a string"))
		os.Exit(1)
	}
}

func ServerRead() (string, string) {
	n, addr, _ := serverConn.ReadFromUDP(buffer)
	str_addr := fmt.Sprintf("%s", addr)
	addrs[str_addr] = addr
	return string(buffer[0:n]), str_addr
}

func ServerReadWrapper(vm *VirtualMachine) {
	message, addr := ServerRead()
	vm.Push(addr)
	vm.Push(message)
}

func ServerSend(message, addr string) {
	b := []byte(message)
	serverConn.WriteTo(b, addrs[addr])
}

func ServerSendWrapper(vm *VirtualMachine) {
	if vm.Length() < 1 {
		fmt.Println(errors.New("ERROR: Tried to call ServerSend when length of stack is less than 1"))
		os.Exit(1)
	}
	message := vm.Pop()
	addr := vm.Pop()

	switch message.(type) {
	case string:
		switch addr.(type) {
		case string:
			ServerSend(message.(string), addr.(string))
		default:
			fmt.Println(errors.New("ERROR: Address must be String when using ServerSend"))
			os.Exit(1)
		}
	default:
		fmt.Println(errors.New("ERROR: Cannot send Non-string data from server"))
		os.Exit(1)
	}
}

var conn *net.UDPConn

func Connect(ip, port string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", ip+":"+port)
	LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	conn, _ = net.DialUDP("udp", LocalAddr, ServerAddr)
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

func ClientSend(message string) {
	buffer := []byte(message)
	conn.Write(buffer)
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
func ClientRead() string {
	b := make([]byte, 1024)
	n, _ := conn.Read(b)
	return string(b[0:n])
}

func ClientReadWrapper(vm *VirtualMachine) {
	vm.Push(ClientRead())
}

// client
//
// func main() {
// 	Connect("127.0.0.1", "8080")
// 	for {
// 		time.Sleep(time.Second)
// 		Send("Hey there")
// 		fmt.Println("got: ", Read())
// 	}
// }

// server
//
// func main() {
// 	buffer = make([]byte, 1024)
// 	addrs = make(map[string]*net.UDPAddr)

// 	Listen("8080")

// 	for {
// 		text, addr := Read()
// 		fmt.Println("Received ", text, " from ", addr)
// 		Send("Hello from server!", addr)
// 	}
// }

func InstallLibrary(vm *VirtualMachine) {
	buffer = make([]byte, 4096)
	addrs = make(map[string]*net.UDPAddr)
	vm.Library["udp.listen"] = ListenWrapper
	vm.Push("{ udp.listen % }")
	vm.Push("udp.Listen")
	vm.Op_store(false)
	vm.Library["udp.serversend"] = ServerSendWrapper
	vm.Push("{ udp.serversend % }")
	vm.Push("udp.ServerSend")
	vm.Op_store(false)
	vm.Library["udp.serverread"] = ServerReadWrapper
	vm.Push("{ udp.serverread % }")
	vm.Push("udp.ServerReceive")
	vm.Op_store(false)
	vm.Library["udp.connect"] = ConnectWrapper
	vm.Push("{ udp.connect % }")
	vm.Push("udp.Connect")
	vm.Op_store(false)
	vm.Library["udp.clientsend"] = ClientSendWrapper
	vm.Push("{ udp.clientsend % }")
	vm.Push("udp.ClientSend")
	vm.Op_store(false)
	vm.Library["udp.clientread"] = ClientReadWrapper
	vm.Push("{ udp.clientread % }")
	vm.Push("udp.ClientReceive")
	vm.Op_store(false)
}
