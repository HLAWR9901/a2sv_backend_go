package services

import(
	"library_management/models"
	"errors"
)

type LibraryManager interface{
	AddBook(book models.Book) error
	RemoveBook(bookID int) error
	BorrowBook(bookID, memberID int) error
	ReturnBook(bookID, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberId int)[]models.Book
}

type Library struct{
	books map[int]*models.Book
	members map[int]*models.Member
}

func (lib *Library) AddBook(book models.Book) error{
	if _,ok:=lib.books[book.ID];ok{
		return errors.New("book already exists")
	}
	lib.books[book.ID] = &book
	return nil
} 

func (lib *Library) RemoveBook(bookID int) error{
	if _,ok:=lib.books[bookID]; !ok{
		return errors.New("book not found")
	}
	delete(lib.books,bookID)
	return nil
}

func (lib *Library) BorrowBook(bookID,memberID int) error{
	b,bexists := lib.books[bookID]
	if !bexists{
		return errors.New("book not found")
	}
	m,mexists:=lib.members[memberID]
	if !mexists{
		return errors.New("member not found")
	}
	if b.Status!=models.Available{
		return errors.New("book is not available")
	}
	b.Status = models.Borrowed
	// If Borrowed books are maps
	// m.BorrowedBooks[b.ID] = *b
	m.BorrowedBooks = append(m.BorrowedBooks, *b)
	return nil
}

func (lib *Library) ReturnBook(bookID,memberID int) error{
	b,bexists := lib.books[bookID]
	if !bexists{
		return errors.New("book not found")
	}
	m,mexists:=lib.members[memberID]
	if !mexists{
		return errors.New("member not found")
	}
	if b.Status!=models.Borrowed{
		return errors.New("book is not borrowed")
	}
	b.Status = models.Available

	// O(1)? - comment: If req didn't specify the borrowedbooks to be slice (maps would be better choice)
	// delete(m.BorrowedBooks,bookID)

	idx := -1
	for i,book:= range m.BorrowedBooks{
		if book.ID == bookID{
			idx = i
			break
		}
	}
    if idx == -1 {
        return errors.New("borrow record not found")
    }
	m.BorrowedBooks = append(m.BorrowedBooks[:idx],m.BorrowedBooks[idx+1:]...)
	return nil
}

func (lib *Library) ListAvailableBooks() []models.Book{
	avail := make([]models.Book,10)
	for _,b:= range lib.books{
		if b.Status==models.Available{
			avail = append(avail,*b)
		}
	}
	return avail
}

func (lib *Library) ListBorrowedBooks(memberId int)[]models.Book{
	return lib.members[memberId].BorrowedBooks
}

