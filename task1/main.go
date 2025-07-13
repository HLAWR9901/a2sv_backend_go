package main

import (
	"bufio" //must use it for the name First, middle, last
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ANSI color codes : obviously searched :)
const (
    Reset = "\033[0m"
    Bold = "\033[1m"
    Cyan = "\033[36m"
    Yellow = "\033[33m"
    Green = "\033[32m"
    Magenta = "\033[35m"
    LightBlue = "\033[94m"
	Red       = "\033[31m"
)

// student struct type
type Student struct{
	fname, mname, lname string 
	grades map[string]float64
	average float64
}

// Initializing methods for student struct
type Action interface{
	displayInfo() error
	inputData() error
	updateGrade(s string, val float64) error
	deleteGrade(s string) error
	addGrade(i int)
	getAverage() (float64,error)
}

func (stud *Student) displayInfo() error{
	n := len(stud.grades)
	full := stud.fname + " "+stud.mname+" "+stud.lname
	//border
	fmt.Println(Bold + LightBlue + "╔" + strings.Repeat("═", 50) + "╗" + Reset)
	//name
	fmt.Printf(Bold+LightBlue+"║ Student: "+Reset+"%s"+Bold+LightBlue+"%s║\n"+Reset, full, strings.Repeat(" ", 40-len(full)),)
	// number of courses
	fmt.Printf(Bold+LightBlue+"║ Courses: "+Reset+"%d"+Bold+LightBlue+"%s║\n"+Reset,n, strings.Repeat(" ", 40-len(fmt.Sprintf("%d", n))),)
	// border
    fmt.Println(Bold + LightBlue + "╠" + strings.Repeat("═", 50) + "╣" + Reset)
    // table header
	fmt.Printf(Bold+Magenta+"║ %2s │ %-20s │ %5s ║"+Reset+"\n","#", "Subject", "Grade",)
    fmt.Println(Bold + LightBlue + "╠" + strings.Repeat("═", 50) + "╣" + Reset)
	// rows
	i := 1
	for key,value:= range stud.grades{
		// Red: Fail, Yellow: Need Improvement, Green: Good
		color := Green
		if value < 80 && value > 49 {
			color = Yellow
		} else if value < 50{
			color = Red
		}
		fmt.Printf(Bold+Yellow+"║ %2d │ "+Reset+"%-20s"+Bold+" │ "+color+"%5.1f"+Reset+Bold+Yellow+" ║\n"+Reset, i, key, value,)
        i++
	}
    fmt.Println(Bold + LightBlue + "╠" + strings.Repeat("═", 50) + "╣" + Reset)
    fmt.Printf(Bold+Cyan+"║ Average"+Reset+"%38.2f"+Bold+Cyan+" ║\n"+Reset, stud.average)
    fmt.Println(Bold + LightBlue + "╚" + strings.Repeat("═", 50) + "╝" + Reset)
	return nil
}


func (stud *Student) inputData() error{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Your Full Name: <First Name> <Middle Name> <Last Name>\nFull Name: ")
	nameInput, _ := reader.ReadString('\n')
	// handle other unicode instead of letters
	for _, r := range nameInput {
        if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
            return errors.New("names must include only letters")
        }
    }
	nameInput = strings.TrimSpace(nameInput) // leading and trailing removal
	names := strings.Fields(nameInput) // split with white space
	if len(names) != 3 {
		return errors.New("full name must include first, middle, and last names")
	}
	stud.fname = names[0]
	stud.lname = names[2]
	stud.mname = names[1]

	var n int
	for {
		fmt.Print("Number of courses you have taken: ")
		inp, _ := reader.ReadString('\n')
		inp = strings.TrimSpace(inp)
		var err error
		n, err = strconv.Atoi(inp)
		if err != nil || n <= 0 {
			fmt.Println(Red + "Invalid input. Please enter a positive integer." + Reset)
			continue
		}
		break
	}
	fmt.Println("List the names and grades of each courses below.")
	for i:=1; i<=n;i++{
		stud.addGrade(i)
	}
	return nil
}


func (stud *Student) updateGrade(s string, val float64) error {
	key := strings.ToLower(s)
	if _,ok:= stud.grades[key]; !ok{
		return errors.New("Subject not found")
	}
	for val<0{
		fmt.Println(Red + "Invalid input. Please enter a positive integer." + Reset)
		fmt.Print("Enter grade: ")
		fmt.Scan(&val)
	}
	n := float64(len(stud.grades))
	total := stud.average*n + val - stud.grades[key]
	stud.average = total/n
	stud.grades[key]=val
	return nil
}

func (stud *Student) deleteGrade(s string) error{
    key := strings.ToLower(s)
    val, ok := stud.grades[key]
    if !ok {
        return errors.New("Subject not found")
    }
	n := float64(len(stud.grades))
	total := stud.average*n - val
	stud.average = total/(n-1)
	delete(stud.grades,key)
	return nil
}

func (stud *Student) addGrade(i int){
	reader := bufio.NewReader(os.Stdin)
	// Subject input
	var subject string
	for {
		if i > 0 {
			fmt.Printf("Subject %d: ", i)
		} else {
			fmt.Print("New Subject: ")
		}
		subject, _ = reader.ReadString('\n')
		subject = strings.TrimSpace(subject)
		if subject == "" {
			fmt.Println(Red + "Subject name cannot be empty" + Reset)
			continue
		}
		ok:=false
		for _, r := range subject {
        	if !unicode.IsLetter(r) && !unicode.IsSpace(r){
            	fmt.Println("Subject names must include only letters")
				ok = true
				break
			}	
        }
		if ok{
			continue
		}
		break
	}
	var mark float64
	for {
		if i>0{
			fmt.Printf("%s Grade: ",subject)
		} else {
			fmt.Printf("Grade for the new subject: ")
		}
		grade,_:=reader.ReadString('\n')
		grade=strings.TrimSpace(grade)
		var err error
		mark,err = strconv.ParseFloat(grade,64)
		if err!=nil{
			fmt.Println(Red + "Invalid grade. Please enter a number." + Reset)
			continue
		}
		if mark < 0 || mark >100 {
			fmt.Println(Red + "Grade must be between 0 - 100 inclusive." + Reset)
			continue
		}
		break
	}
	key := strings.ToLower(subject)
	n := float64(len(stud.grades))
	stud.average = (stud.average*n + mark) / (n + 1)
	stud.grades[key] = mark
}

func (stud *Student) getAverage() (float64,error){
	if len(stud.grades)==0{
		return 0,errors.New("no courses registered")
	}
	return stud.average,nil
}

func main(){
	s := Student{ grades: make(map[string]float64) }
	var I Action = &s
	reader := bufio.NewReader(os.Stdin)
outer:
	for {
		fmt.Println("\n" + Bold + "Grade Calculator" + Reset)
		fmt.Println("1. Register")
		fmt.Println("2. Exit")
		fmt.Print(">")
		inp, _ := reader.ReadString('\n')
		inp = strings.TrimSpace(inp)
		input, err := strconv.Atoi(inp)
		if err != nil {
			fmt.Println(Red + "Invalid input. Please enter 1 or 2." + Reset)
			continue
		}
		switch input{
		case 1:
			if err:=I.inputData(); err!=nil{
				fmt.Println(Red+"Error:", err, Reset)
				continue outer
			}
			I.displayInfo()

			for {
				fmt.Println("\n" + Bold + "Menu" + Reset)
				fmt.Println("1. Calculate average")
				fmt.Println("2. Add Grade")
				fmt.Println("3. Update Grade")
				fmt.Println("4. Delete Grade")
				fmt.Println("5. Print Report")
				fmt.Println("6. Exit")
				fmt.Print("Choose option: ")
				chInp, _ := reader.ReadString('\n')
				chInp = strings.TrimSpace(chInp)
				choice, err := strconv.Atoi(chInp)
				if err!=nil || !(input>0 && input<7) {
					fmt.Println(Red + "Invalid input. Please enter 1 or 2." + Reset)
					continue
				}
				switch choice{
				case 1:
					val,err := I.getAverage()
					if err!=nil{
						fmt.Println(Red+"Error:", err, Reset)
					} else{
						fmt.Printf("Average: %.2f\n",val)
					}
				case 2:
					I.addGrade(0)
				case 3:
					fmt.Print("Subject to update: ")
					subject, _ := reader.ReadString('\n')
					subject = strings.TrimSpace(subject)
					for {
						fmt.Print("New grade:")
						Mark, _ := reader.ReadString('\n')
						Mark = strings.TrimSpace(Mark)
						newMark, err := strconv.ParseFloat(Mark, 64)
						if err!=nil{
							fmt.Println(Red + "Invalid grade. Enter a number." + Reset)
							continue
						} else if newMark<0 || newMark>100{
							fmt.Println(Red + "Grade must be between 0 - 100 inclusive" + Reset)
							continue
						}
						if err := I.updateGrade(subject, newMark); err != nil {
							fmt.Println(Red+"Error:", err, Reset)
						} else {
							fmt.Println(Green + "Grade updated successfully!" + Reset)
						}
						break
					}
				case 4:
					fmt.Print("Subject to delete: ")
					subject, _ := reader.ReadString('\n')
					subject = strings.TrimSpace(subject)
					if err := I.deleteGrade(subject); err!=nil{
						fmt.Printf("ERROR: %v\n",err)
					} else{
						fmt.Println(Green+"Successful Deletion!"+Reset)
					}
				case 5:
					I.displayInfo()
				case 6:
					fmt.Print("Thank you for using this program. \nExiting program")
					for i := 0; i < 3; i++ {
						fmt.Print(".")
						time.Sleep(500 * time.Millisecond)
					}
					os.Exit(0)
				default:
					fmt.Println(Red + "Invalid option. Choose 1-6" + Reset)
				}
			}
		case 2:
			fmt.Print("Thank you for using this program. \nExiting program")
			for i := 0; i < 3; i++ {
				fmt.Print(".")
				time.Sleep(500 * time.Millisecond)
			}
			os.Exit(0)
		default:
			fmt.Println(Red + "Invalid option. Choose 1 or 2" + Reset)
		}
	}
}