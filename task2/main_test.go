package main
import (
	"reflect"
	"testing"
)

func TeststrMap(t *testing.T) {
	testcases := []struct{
		cases string
		input string
		result map[string]int
	}{
		{"empty","",map[string]int{}},
		{"one word","one",map[string]int{"one":1},},
		{"multiple words","hi bye show d d hi man so d d",map[string]int{"hi":2,"bye":1,"show":1,"d":4,"man":1,"so":1}},
	}
	for _,test:=range testcases{
		output:=strMap(&test.input)
		if !reflect.DeepEqual(output,test.result){
			t.Errorf("%s: strMap(%q) = %v; want %v", test.cases, test.input, output, test.result)
		}
	}
}

func TestCheckPalindrome(t *testing.T) {
    testcases := []struct {
        cases string
        input   string
        result bool
    }{
        {"empty", "", true},
        {"single char", "x", true},
        {"even palindrome", "abba", true},
        {"odd palindrome", "radar", true},
        {"not palindrome", "hello", false},
        {"mixed case non-palindrome", "AbBa", false},
        {"with space non-palindrome", "a b a", true},
    }

    for _, test := range testcases {
        output := checkPalindrome(test.input)
        if output != test.result {
            t.Errorf("%s: checkPalindrome(%q) = %v; want %v", test.cases, test.input, output, test.result)
        }
    }
}