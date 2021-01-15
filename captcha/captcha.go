package captcha

import "fmt"
import "time"
import "math/rand"
import "github.com/google/uuid"
import "sync"
import "log"



type Captcha struct {
	pattern int
	leftOperand int
	operator int
	rightOperand int
}

/*func (c Captcha) New(p int, lo int, op int, ro int)  string{
	return "test"
}*/

func New(p int, lo int, op int, ro int)  Captcha{
	return Captcha{p,lo,op,ro}
	//return "test"
}

func (cc Captcha) String() string {
	var result string = ""
	if cc.pattern == 1{
		// IntToAlphabet(cc.leftOperand) + parseOperator(cc.operator)+ cc.rightOperand
		result = fmt.Sprintf("%d %s %s",cc.leftOperand, parseOperator(cc.operator), IntToAlphabet(cc.rightOperand))
		
	}else if cc.pattern == 2{
		result = fmt.Sprintf("%s %s %d",IntToAlphabet(cc.leftOperand), parseOperator(cc.operator), cc.rightOperand)
	}

	return result
}

var src = rand.NewSource(time.Now().UnixNano());
var rnd = rand.New(src);
var store = map[string]int{}
var mux sync.Mutex
// atomic

func Answer(key string, ans int) bool{
	mux.Lock();
	defer mux.Unlock();

	if v,ok:= store[key];ok {
		delete(store, key)
		return v == ans 
	}else{
		log.Printf("not found %s key in store\n",key)
		return false
	}
}

func KeyQuestion() (string, string) {
	pattern, leftOperand, operator, rightOperand := rnd.Intn(2)+1,rnd.Intn(9)+1,rnd.Intn(3)+1,rnd.Intn(9)+1
	answer := 0
	switch{
	case operator == 1:
		answer = leftOperand + rightOperand
	case operator == 2:
		answer = leftOperand - rightOperand
	case operator == 3:
		answer = leftOperand * rightOperand
		
	}
	
	cc := New(pattern, leftOperand, operator, rightOperand)
	key := uuid.New().String()
	store[key] = answer

	return key,cc.String() 
}

func IntToAlphabet(i int) string{
	switch{
	case i == 0:
		return "zero"
	case i == 1:
			return "one"
		case i == 2:
			return "two"
		case i == 3:
			return "three"
		case i == 4:
			return "four"
		case i == 5:
			return "five"
		case i == 6:
			return "six"
		case i == 7:
			return "seven"
		case i == 8:
			return "eight"
		case i == 9:
			return "nine"
		default:
			return ""
		}
	}	

	func parseOperator(operator int) string{
		switch{
			case operator == 1: 
			return "+"
			case operator == 2: 
			return "-"
			case operator == 3: 
			return "*"
		default:
			return ""
		}
	}

