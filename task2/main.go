package main
import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"strconv"
)

//Task 1
func strmap(s *string) map[string]int {
	collection := strings.Fields(*s)
	res := make(map[string]int)
	for _,val := range collection{
		res[val]++
	}
	return res
}

//Task 2
func checkPalindrome(s string) bool{
	n := len(s)
	for i:=0; i<int(n/2);i++{
		if s[i]!=s[n-i-1]{
			return false
		}
	}
	return true
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input String: ")
	input,_ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	for{
		fmt.Println("Menu")
		fmt.Println("1. Word Frequency Count")
		fmt.Println("2. Palindrome Check")
		fmt.Print(">")
		choice,_:= reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		n,err := strconv.Atoi(choice)
		if err!=nil{
			fmt.Println("Invalid Input")
			continue
		}
		if n==1{
			fmt.Println(strmap(&input))
		} else if n==2{
			if checkPalindrome(input){
				fmt.Printf("%s is Palindrome.",input)
			}else{
				fmt.Printf("%s is not Palindrome.",input)
			}
		}else{
			fmt.Println("Invalid Input")
		}
		break
	}
}