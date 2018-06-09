package main

    import (
        "./parser"
        "./virtualmachine"
    	"./operators"
	"./lists"
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
                                    "{ 'Password: ' print > ! input > ! x < 'password' password <  {  'You may enter!' print > ! input > ! 0  }  password >  x >  equal > ! if > !  {  endl >  'Nope lmao' add > ! print > ! 1  }  password >  x >  notequal > ! if > ! } func < func >  while > !")
        tokens := parser.Parse()
        // PrintList(tokens)
        VM := virtualmachine.MakeVM(tokens)
        
        
        operators.InstallLibrary(&VM)
        

        lists.InstallLibrary(&VM)
        
        VM.Run()

    }