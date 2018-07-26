package main

import (
	"fmt"

	"./lists"
	"./operators"
	"./parser"
	"./threading"
	"./udp"
	"./virtualmachine"
)

func PrintList(tokens []string) {
	fmt.Print("[")
	for i, token := range tokens {
		fmt.Print("'")
		fmt.Print(token)
		fmt.Print("'")
		if i < len(tokens)-1 {
			fmt.Println(", ")
		}
	}
	fmt.Println("]")
}

func main() {
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
		"{>} pointer_read < " +
		"{<} pointer_write < " +
		"{__unpack_func__ > #} unpack < " +
		"{x < ':' x > + ':' + |} print_lit < " +
		"{func < 1 func > { } &} while < " +
		"{condition < function < {function > ! 0 condition <} {condition >} &} if < " +
		"{ @ self < lists.List > ! self >  list , self <  { self <  a <  a >  self >  list . lists.Append > ! self >  list , self <  self > } append < append >  self >  append , self <  { self <  self >  list . lists.Pop > ! self >  list , self <  self > } pop < pop >  self >  pop , self <  { self <  self >  list . lists.Length > ! } len < len >  self >  len , self <  { self <  self >  list . lists.Items > ! } items < items >  self >  items , self <  { self <  n <  n >  self >  list . lists.Index > ! } index < index >  self >  index , self <  self > } list < " +
		"{equal > ! not > !} notequal < " +
		"{ s <  s >  print > ! endl >  print > ! } println < " +
		"{ case <  if_then <  else_then <   {  if_then > !  }  case >  if > !  {  else_then > !  }  case >  not > ! if > ! } ifelse < " +
		"{ l <  elem <  0 r <  {   {  1 r <  }  elem >  x >  equal > ! if > !  }  l >  'x' for > ! r > } in < { list_append <  e <  e >  list_append >  list_append >  append . ! list_append < list_append > } append < { list_pop <  list_pop >  list_pop >  append . ! list_pop < list_pop > } pop < { list_len <  list_len >  list_len >  len . ! r < r > } len < { v <  list_arg <  f <   {  v >  pointer_read > ! list_arg >  list_arg >  index . ! v >  pointer_write > ! f > !  }  1 list_arg >  list_arg >  len . ! sub > ! 0 range > ! v >  for_reserved > ! } for < " +
		"{ '\\8080' udp.Listen > ! 'Server started.' print > ! endl >  print > ! list > ! clients < { udp.ServerReceive > ! 'text' unpack > ! 'addr' unpack > !  {  addr >  clients >  append > ! clients < clients >  len > ! println > !  }  addr >  clients >  in > ! not > ! if > !  {   {  c >  text >  udp.ServerSend > !  }  addr >  c >  notequal > ! if > !  }  clients >  'c' for > ! loop > ! } loop < loop > ! } main < main > !")
	tokens := parser.Parse()
	// PrintList(tokens)
	VM := virtualmachine.MakeVM(tokens)

	operators.InstallLibrary(&VM)

	threading.InstallLibrary(&VM)

	udp.InstallLibrary(&VM)

	lists.InstallLibrary(&VM)

	VM.Run()

}
