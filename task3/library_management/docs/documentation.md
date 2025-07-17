# 📚 Library Management System

A terminal-based **Library Management System** built in Go that enables librarians to manage books and members effectively.

---

## 📁 Project Structure

```
library_management/
├── main.go                # Entry point of the application
├── controllers/           # Handles user interface and input/output
│   └── library_controller.go
├── models/                # Contains data structures for Book and Member
│   ├── book.go
│   └── member.go
├── services/              # Core logic for managing books and members
│   └── library_service.go
```

---

## 🚀 Features

* ✅ Add and remove books
* ✅ Borrow and return books
* ✅ Add and remove members
* ✅ List all available books
* ✅ List all borrowed books per member
* 🖥️ CLI interface with color-coded messages and ASCII art

---

## 🧱 Components

### 📦 `models`

#### `Book`

```go
type Book struct {
    ID     int
    Title  string
    Author string
    Status State
}
```

* `Status` can be:

  * `Available`
  * `Borrowed`

#### `Member`

```go
type Member struct {
    ID            int
    Name          string
    BorrowedBooks []Book
}
```

---

### 🛠️ `services`

#### Interface: `LibraryManager`

```go
type LibraryManager interface {
    AddBook(book models.Book) error
    RemoveBook(bookID int) error
    BorrowBook(bookID, memberID int) error
    ReturnBook(bookID, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberId int) []models.Book
    AddMember(member models.Member) error
    RemoveMember(memberId int) error
}
```

#### Struct: `Library`

* Stores:

  * `Books`: map of all books
  * `Members`: map of all registered members
  * `NextBook` / `NextMember`: auto-increment counters

---

### 👤 `controllers`

#### Key Function: `Dashboard()`

Main function that:

* Initializes library state
* Presents CLI menu options
* Handles all core actions (Add/Remove/Borrow/Return/List)

---

## 📌 Usage

To run the application:

```sh
go run main.go
```

Follow the menu options in the terminal UI to perform different library operations.

---

## 📒 Sample Menu

```
Menu
1. Add new book
2. Remove a book
3. Borrow a book
4. Return a book
5. Add member
6. Remove member
7. List all available books
8. List all borrowed books
9. Exit
```

---

## 🧪 Sample Actions

* Add a Book → enter title and author
* Borrow a Book → provide valid book and member IDs
* Return a Book → requires previous borrow history
* List Books → shows formatted table of available or borrowed books

---

## 🎨 Extras

* Terminal UI uses ANSI escape codes for color
* Custom ASCII art is printed on exit

---
