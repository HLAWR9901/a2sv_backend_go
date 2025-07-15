package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Cyan      = "\033[36m"
	Yellow    = "\033[33m"
	Green     = "\033[32m"
	Magenta   = "\033[35m"
	LightBlue = "\033[94m"
	Red       = "\033[31m"
)

// General
func Home() {
	clearScreen()
	fmt.Println(Bold + LightBlue + "╔" + strings.Repeat("=", 60) + "╗" + Reset)
	fmt.Println(Bold + LightBlue + "║               " + Reset + Bold + Green + "LIBRARY MANAGEMENT SYSTEM" + Reset + Bold + LightBlue + strings.Repeat(" ", 20) + "║" + Reset)
	fmt.Println(Bold + LightBlue + "╚" + strings.Repeat("═", 60) + "╝" + Reset)
}
func Menu() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n" + Bold + "Menu" + Reset)
	fmt.Println("1. Add new book")
	fmt.Println("2. Remove a book")
	fmt.Println("3. Borrow a book")
	fmt.Println("4. Return a book")
	fmt.Println("5. Add member")
	fmt.Println("6. Remove member")
	fmt.Println("7. List all available books")
	fmt.Println("8. List all borrowed books")
	fmt.Println("9. Exit")
	fmt.Print("> ")
	inputData, _ := reader.ReadString('\n')
	inputData = strings.TrimSpace(inputData)
	input, err := strconv.Atoi(inputData)
	if err != nil || input < 1 || input > 9 {
		return 0, errors.New("invalid input. Please enter 1 - 9")
	}
	return input, nil

}

// Specific
func clearScreen() {
	// ESC -> cursor home -> [2J clears screen
	fmt.Print("\x1b[H\x1b[2J")
}
func validName(s string) string {
	reader := bufio.NewReader(os.Stdin)
	var name string
	for {
		fmt.Printf("%s: ", s)
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Printf(Yellow+"%s should be atleast 1 character long.\n"+Reset, s)
			continue
		}
		break
	}
	name = strings.ToLower(name)
	return name
}

func Dashboard() {
	lib := services.Library{
		Books:      make(map[int]*models.Book),
		Members:    make(map[int]*models.Member),
		NextBook:   1,
		NextMember: 1,
	}
	var librarian services.LibraryManager = &lib
	for {
		Home()
		for {
			input, err := Menu()
			if err != nil {
				fmt.Println(Red + err.Error() + Reset)
				continue
			}
			switch input {
			case 1:
				for {
					title := validName("Book Title")
					author := validName("Book Author")
					b := models.Book{ID: lib.NextBook, Title: title, Author: author, Status: models.Available}
					lib.NextBook++
					err := librarian.AddBook(b)
					if err != nil {
						fmt.Printf("%s%s%s\n", Red, err.Error(), Reset)
						continue
					}
					fmt.Printf("%s%s%s\n", Green, "Successfully Added", Reset)
					break
				}
			case 2:
				for {
					n := len(lib.Books)
					idData := validName("Book ID")
					id, err := strconv.Atoi(idData)
					if err != nil || id < 1 || id > n {
						fmt.Printf("%sInvalid Book ID.%s\n", Yellow, Reset)
						continue
					}
					e := librarian.RemoveBook(id)
					if e != nil {
						fmt.Printf("%s%s%s\n", Red, e.Error(), Reset)
						continue
					}
					fmt.Printf("%s%s%s\n", Green, "Successfully Removed", Reset)
					break
				}
			case 3:
				var bookId, memberId int
				// book
				for {
					bookIdData := validName("Book ID")
					id, err := strconv.Atoi(bookIdData)
					if err != nil || id < 1 || id >= lib.NextBook {
						fmt.Printf("%sInvalid Book ID.%s\n", Yellow, Reset)
						continue
					}
					bookId = id
					break
				}
				// member
				for {
					memberIdData := validName("Member ID")
					id, err := strconv.Atoi(memberIdData)
					if err != nil || id < 1 || id > len(lib.Members) {
						fmt.Printf("%sInvalid Member ID.%s\n", Yellow, Reset)
						continue
					}
					memberId = id
					break
				}
				// borrow
				if e := librarian.BorrowBook(bookId, memberId); e != nil {
					fmt.Printf("%s%s%s\n", Red, e.Error(), Reset)
				} else {
					fmt.Printf("%sSuccessfully Borrowed%s\n", Green, Reset)
				}
			case 4:
				var bookId, memberId int
				// book
				for {
					bookIdData := validName("Book ID")
					id, err := strconv.Atoi(bookIdData)
					if err != nil || id < 1 || id > len(lib.Books) {
						fmt.Printf("%sInvalid Book ID.%s\n", Yellow, Reset)
						continue
					}
					bookId = id
					break
				}
				// member
				for {
					memberIdData := validName("Member ID")
					id, err := strconv.Atoi(memberIdData)
					if err != nil || id < 1 || id > len(lib.Members) {
						fmt.Printf("%sInvalid Member ID.%s\n", Yellow, Reset)
						continue
					}
					memberId = id
					break
				}
				// return
				if e := librarian.ReturnBook(bookId, memberId); e != nil {
					fmt.Printf("%s%s%s\n", Red, e.Error(), Reset)
				} else {
					fmt.Printf("%sSuccessfully Returned%s\n", Green, Reset)
				}
			case 5:
				for {
					name := validName("Member Name")
					m := models.Member{
						ID: lib.NextMember,
						Name: name,
						BorrowedBooks: []models.Book{},
					}
					if err := librarian.AddMember(m); err != nil {
						fmt.Printf("%s%s%s\n", Red, err.Error(), Reset)
						continue
					}
					lib.NextMember++
					fmt.Printf("%sSuccessfully Added Member [ID: %d]%s\n", Green, m.ID, Reset)
					break
				}

			case 6:
				for {
					n := len(lib.Members)
					idData := validName("Member ID")
					id, err := strconv.Atoi(idData)
					if err != nil || id < 1 || id > n {
						fmt.Printf("%sInvalid Member ID.%s\n", Yellow, Reset)
						continue
					}
					e := librarian.RemoveMember(id)
					if e != nil {
						fmt.Printf("%s%s%s\n", Red, e.Error(), Reset)
						continue
					}
					fmt.Printf("%s%s%s\n", Green, "Successfully Removed", Reset)
					break
				}
			case 7:
				data := lib.ListAvailableBooks()
				if len(data) == 0 {
					fmt.Println(Yellow + "There are no available books in the Library LOL" + Reset)
				} else {
					width := 60
					fmt.Println(Bold + LightBlue + "╔" + strings.Repeat("═", width) + "╗" + Reset)
					fmt.Printf(Bold+LightBlue+"║ %-4s║ %-22s║ %-32s║\n", "ID", "Title", "Author"+Reset)
					fmt.Println(Bold + LightBlue + "╠" + strings.Repeat("═", width) + "╣" + Reset)

					for _, book := range data {
						fmt.Printf(Bold+LightBlue+"║ "+Reset+Magenta+"%-4d"+Reset+Bold+LightBlue+"║ "+Reset+Green+"%-22s"+Reset+Bold+LightBlue+"║ "+Reset+"%-28s"+Bold+LightBlue+"║\n",
							book.ID, book.Title, book.Author)
					}
					fmt.Println(Bold + LightBlue + "╚" + strings.Repeat("═", width) + "╝" + Reset)
				}
			case 8:
				// member
				var memberId int
				for {
					memberIdData := validName("Member ID")
					id, err := strconv.Atoi(memberIdData)
					if err != nil || id < 1 || id > len(lib.Members) {
						fmt.Printf("%sInvalid Member ID.%s\n", Yellow, Reset)
						continue
					}
					memberId = id
					break
				}
				data := lib.ListBorrowedBooks(memberId)
				if len(data) == 0 {
					fmt.Println(Yellow + "There are no Borrowed books" + Reset)
				} else {
					width := 60
					fmt.Println(Bold + LightBlue + "╔" + strings.Repeat("═", width) + "╗" + Reset)
					fmt.Printf(Bold+LightBlue+"║ %-4s║ %-22s║ %-32s║\n", "ID", "Title", "Author"+Reset)
					fmt.Println(Bold + LightBlue + "╠" + strings.Repeat("═", width) + "╣" + Reset)

					for _, book := range data {
						fmt.Printf(Bold+LightBlue+"║ "+Reset+Magenta+"%-4d"+Reset+Bold+LightBlue+"║ "+Reset+Green+"%-22s"+Reset+Bold+LightBlue+"║ "+Reset+"%-28s"+Bold+LightBlue+"║\n",
							book.ID, book.Title, book.Author)
					}
					fmt.Println(Bold + LightBlue + "╚" + strings.Repeat("═", width) + "╝" + Reset)
				}
			case 9:
				fmt.Print("Exiting program ")
				for i := 0; i < 3; i++ {
					fmt.Print(" . ")
					time.Sleep(1 * time.Second)
				}
				printAwesome()
				os.Exit(0)
			}
		}
	}
}

func printAwesome() {
	fmt.Println(`
           ,,                               .-.
          || |                               ) )
          || |   ,                          '-'
          || |  | |
          || '--' |
    ,,    || .----'
   || |   || |
   |  '---'| |
   '------.| |                                  _____
   ((_))  || |      (  _                       / /|\ \
   (o o)  || |      ))("),                    | | | | |
____\_/___||_|_____((__^_))____________________\_\|/_/__
	
	`)
}
