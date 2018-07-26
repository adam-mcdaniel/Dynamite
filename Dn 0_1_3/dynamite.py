import os
import sys
import re

def Symbol(a):
    if str(a) == "define":
        return "<"
    if str(a) == "function":
        return "function <"
    elif str(a) == "call":
        return "> !"
    elif str(a) == "lambda":
        return "!"
    elif str(a) == "return":
        return ">"    
    elif str(a) == "dot":
        return "."
    elif str(a) == '':
        return " "
    else:
        return str(a)


def commentRemover(text):
    def replacer(match):
        s = match.group(0)
        if s.startswith(';'):
            return " " # note: a space and not an empty string
        else:
            return s
    pattern = re.compile(
        r';.*?$|/\*.*?\*/|\'(?:\\.|[^\\\'])*\'|"(?:\\.|[^\\"])*"',
        re.DOTALL | re.MULTILINE
    )
    return re.sub(pattern, replacer, text)


def tokenize(script):
    script = script.replace("(", " ( ").replace(")", " ) ")
    script = commentRemover(script)
    tokens = []
    current_token = ""
    in_string = False
    for char in script:
        if in_string:
            if char == "\"":
                in_string = False
                tokens.append("`" + current_token + "`")
                current_token = ""
            else:
                current_token += char
        else:
            if char == "\"":
                tokens.append(current_token)                
                in_string = True
            elif char == " ":
                tokens.append(current_token)                
                current_token = ""
            else:
                current_token += char
    
    tokens = list(map(lambda a: a.replace("\n",""), tokens))
    tokens = list(filter(lambda a: a != '', tokens))

    # print(tokens)
    return tokens
    
def parse(program):
    return read_from_tokens(tokenize(program))

def read_from_tokens(tokens):
    if len(tokens) == 0:
        raise SyntaxError('unexpected EOF')
    token = tokens.pop(0)

    if token == '(':
        L = []
        while tokens[0] != ')':
            L.append(read_from_tokens(tokens))
        tokens.pop(0) # pop off ')'
        if len(L) > 1:
            if L[0] == "<":
                if len(L) == 4 and type(L[2]) == list and type(L[3]) == list:
                    # print(L)
                    L[2] = list(map(lambda s: s + " < ", L[2]))
                    L[3] = L[2] + L[3]
                    del L[2]
                    L = L[::-1]
                    L.insert(0, "{")
                    L.insert(2, "}")
                elif len(L) == 4 and type(L[1]) == list and type(L[3]) == str:
                    L = [(lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t)(L[3]), L[1], L[2], "~", L[1] if L[1] != "@" else '']
                elif len(L) == 4:
                    L = [(lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t)(L[3]), (lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t)(L[1]), L[2], "~", L[1] + " < " if L[1] != "@" else '']
                elif len(L) == 3:
                    L = L[::-1]
                    # if not "`" in L[0]
                    L = [(lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t)(L[0])] + L[1:]
                # print(L)
            elif L[0] == "> !":
                L = L[::-1]
                L = list(map((lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t), L[:-2])) + list(L[-2:])
            elif L[0] == "!":
                L = L[::-1]
                # L = list(map((lambda t: str(t) + list(L[-2:])
            elif L[0] == ">":
                L = L[::-1] if type(L[1]) != int and type(L[1]) != float and not "`" in L[1] and L[1][-1] != "> !" and "{" not in L[1][0] else [L[1]]

            elif L[0] == ".":
                L = [(lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) and not "@" in str(t) else t)(L[1]), L[2], '.']
            elif L[0] == "function <":
                # print(L)
                L = [" { ", L[1], " } "]
                # L = [(lambda t: str(t) + " > " if not str(t).replace(">","").replace(" ","").replace(".","").isdigit() and not "[" in str(t) and not "`" in str(t) else t)(L[1]), L[2], '.']
        
        return L
    elif token == ')':
        raise SyntaxError('unexpected )')
    else:
        return atom(token)

def atom(token):
    try: return int(token)
    except ValueError:
        try: return float(token)
        except ValueError:
            return Symbol(token)

def compile(script):
    byte_compile = lambda script: str(parse(script)).replace("[","").replace("]","").replace(",","").replace("'","").replace("\"","").replace("`","'").replace("~",",").replace(". >",".")
    with open(os.path.abspath("src/config.txt"), "r") as f:
        libs = list(map(lambda s: s.replace('\n',''), f.readlines()))
        f.close()

    with open(os.path.abspath("src/main.go"), "w") as f:
        f.write("""package main

    import (
        "./parser"
        "./virtualmachine"
    """ + '\n'.join(list(map(lambda l: "\t\"./" + l +"\"" , libs))) + """
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
        parser := parser.MakeParser("'\\n' endl < " +
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
                                    \"""" + byte_compile(script) + """\")
        tokens := parser.Parse()
        // PrintList(tokens)
        VM := virtualmachine.MakeVM(tokens)
        
        """ + '\n'.join(list(map(lambda l: """
        {}.InstallLibrary(&VM)
        """.format(l, l), libs))) + """
        VM.Run()

    }"""
        )
        f.close()

    os.system("go build  -ldflags=\"-s -w\" \"" + os.path.abspath("src/main.go") + "\"")

help_text = """
Usage:
dynamite path/to/script"""
try:
    with open(os.path.abspath(sys.argv[1]), "r") as f:
        try:
            os.remove(os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
        except:
            pass
        compile(f.read())
        # print(os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
        os.rename(os.path.abspath("main.exe"), os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
        f.close()
except Exception as e:
    print(e)
    print(help_text)
# with open(os.path.abspath(sys.argv[1]), "r") as f:
#     try:
#         os.remove(os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
#     except:
#         pass
#     compile(f.read())
#     print(os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
#     os.rename(os.path.abspath("main.exe"), os.path.dirname(os.path.abspath(sys.argv[1])) + "\\" + os.path.splitext(os.path.basename(os.path.abspath(sys.argv[1])))[0]+".exe")
#     f.close()