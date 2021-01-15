package captcha_test

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"
import "github.com/pallat/todos/captcha"

func TestCaptchaPattern1(t *testing.T){
	operands := []string{"zero","one","two","three","four","five","six","seven","eight","nine"}
	operators := []string{"","+","-","*"}

	for givingOperand, want := range operands {
		for givingOperator, operatorValue := range operators {
		t.Run(fmt.Sprintf("operand %d",givingOperand), func(t *testing.T){
			givingPattern := 1
			lo := 1
			op := givingOperator
			ro := givingOperand
			want := fmt.Sprintf("1 %s %s",operatorValue,want)

			cc := captcha.New(givingPattern,lo,op,ro)

			get := cc.String()
			assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));		
		})
		}
	}
	/*
	givingPattern := 1
	lo := 1
	op := 1
	ro := 1
	want := "1 + one"

	cc := captcha.New(givingPattern,lo,op,ro)

	get := cc.String()
	assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));
	*/
}

func TestCaptchaPattern2(t *testing.T){
	
	operands := []string{"zero","one","two","three","four","five","six","seven","eight","nine"}
	
	for givingOperand, want := range operands {
		t.Run(fmt.Sprintf("operand %d",givingOperand), func(t *testing.T){
			givingPattern := 2
			lo := givingOperand
			op := 1
			ro := 1
			want := fmt.Sprintf("%s + 1",want)

			cc := captcha.New(givingPattern,lo,op,ro)

			get := cc.String()
			assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));		
		})

	}
	
	/*givingPattern := 2
	lo := 1
	op := 1
	ro := 1
	want := "one + 1"

	cc := captcha.New(givingPattern,lo,op,ro)

	get := cc.String()
	assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));
	*/
}

func TestCaptchaPatternDefault(t *testing.T){
	
	
	givingPattern := 1
	lo := 9
	op := 1
	ro := 10
	want := "9 + "

	cc := captcha.New(givingPattern,lo,op,ro)

	get := cc.String()
	assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));
	
}

func TestNumberOne(t * testing.T){
	givingPattern := 1
	lo := 1
	op := 1
	ro := 1
	want := "1 + one"

	cc := captcha.New(givingPattern,lo,op,ro)

	get := cc.String()
	assert.Equal(t,get,want,fmt.Sprintf("giving %d,%d,%d,%d want %s get %s",givingPattern,lo,op,ro,want,get));	
}