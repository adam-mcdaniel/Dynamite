package parser

// "fmt"

type Parser struct {
	line   string
	tokens []string
}

func MakeParser(line string) Parser {
	// fmt.Println("Making Parser...")
	var empty []string
	return Parser{line, empty}
}

func (self *Parser) Load(line string) {
	// fmt.Print("Loading line: ")
	// fmt.Print(line)
	// fmt.Println("...")
	self.line = line
}

func (self *Parser) Parse() []string {
	// fmt.Println("Parsing...")
	var depth int = 0
	var is_string bool = false
	var token string = ""

	for i := 0; i < len(self.line); i++ {
		if self.line[i] == '\'' && depth < 1 {
			if is_string {
				// fmt.Println(token)
				is_string = false
				token += string('\'')
				if token == "''" {
					self.tokens = append(self.tokens, token)
				} else {
					self.tokens = append(self.tokens, token[1:len(token)-1])
				}
				token = ""
				continue
			} else {
				is_string = true
				self.tokens = append(self.tokens, token)
				token = ""
			}
		}
		if is_string {
			token += string(self.line[i])
			if i == len(self.line)-1 {
				self.tokens = append(self.tokens, token)
			}
		} else {
			if self.line[i] == ' ' {
				if depth < 1 {
					self.tokens = append(self.tokens, token)
					token = ""
				} else {
					token += string(' ')
				}
			} else if self.line[i] == '{' {
				if depth < 1 {
					self.tokens = append(self.tokens, token)
					token = "{"
				} else {
					token += string('{')
				}
				depth += 1
			} else if self.line[i] == '}' {
				depth -= 1
				token += string('}')
				if depth < 1 {
					self.tokens = append(self.tokens, token)
					token = ""
				}
			} else {
				if string(self.line[i]) != "\n" {
					token += string(self.line[i])
				} else {
					self.tokens = append(self.tokens, token)
					token = ""
				}
				if i == len(self.line)-1 {
					self.tokens = append(self.tokens, token)
				}
			}
		}
	}

	var returned_tokens []string
	for _, elem := range self.tokens {
		if elem == "" {

		} else if elem == "''" {
			returned_tokens = append(returned_tokens, "")
		} else {
			returned_tokens = append(returned_tokens, elem)
		}
	}

	self.tokens = returned_tokens
	return self.tokens
}
