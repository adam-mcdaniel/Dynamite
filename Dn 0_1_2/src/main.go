package main

    import (
        "./parser"
        "./virtualmachine"
    	"./operators"
	"./lists"
	"./random"
	"./threading"
        "fmt"
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
                                    "{__unpack_func__ > #} unpack < " +
                                    "{x < ':' x > + ':' + |} print_lit < " +
                                    "{func < 1 func > { } &} while < " +
                                    "{condition < function < {function > ! 0 condition <} {condition >} &} if < " +
                                    "{ @ self < lists.List > ! self >  list , self <  { self <  a <  a >  self >  list . lists.Append > ! self >  list , self <  self > } append < append >  self >  append , self <  { self <  self >  list . lists.Pop > ! self >  list , self <  self > } pop < pop >  self >  pop , self <  { self <  self >  list . lists.Length > ! } len < len >  self >  len , self <  { self <  self >  list . lists.Items > ! } items < items >  self >  items , self <  { self <  n <  n >  self >  list . lists.Index > ! } index < index >  self >  index , self <  self > } list < " +
                                    "{equal > ! not > !} notequal < " +
                                    "{ s <  s >  print > ! endl >  print > ! } println < " +
                                    "{ case <  if_then <  else_then <   {  if_then > !  }  case >  if > !  {  else_then > !  }  case >  not > ! if > ! } ifelse < " +
                                    " {  { input > ! s < 'Sending: ' print > ! s >  println > ! s >  threading.PushToChannel > ! send > ! } send < send > !  }  threading.Thread > !  {  {  {  'Receiving: ' print > ! threading.GetBottomChannel > ! println > ! threading.DeleteFromChannel > !  }  0 threading.GetBottomChannel > ! notequal > ! if > ! receive > ! } receive < receive > !  }  threading.Thread > !  {  1  }  while > !")
        tokens := parser.Parse()
        // PrintList(tokens)
        VM := virtualmachine.MakeVM(tokens)
        
        
        operators.InstallLibrary(&VM)
        

        lists.InstallLibrary(&VM)
        

        random.InstallLibrary(&VM)
        

        threading.InstallLibrary(&VM)
        
        VM.Run()

    }