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
	AddMember(member models.Member) error
	RemoveMember(memberId int) error
}

type Library struct{
	Books map[int]*models.Book
	Members map[int]*models.Member
	NextMember,NextBook int
}

func (lib *Library) AddBook(book models.Book) error{
    for _, b := range lib.Books {
        if b.Title == book.Title {
            return errors.New("book title already exists")
        }
    }
	if _,ok:=lib.Books[book.ID];ok{
		return errors.New("book already exists")
	}
	lib.Books[book.ID] = &book
	return nil
} 

func (lib *Library) RemoveBook(bookID int) error{
	if _,ok:=lib.Books[bookID]; !ok{
		return errors.New("book not found")
	}
	delete(lib.Books,bookID)
	return nil
}

func (lib *Library) BorrowBook(bookID,memberID int) error{
	b,bexists := lib.Books[bookID]
	if !bexists{
		return errors.New("book not found")
	}
	m,mexists:=lib.Members[memberID]
	if !mexists{
		return errors.New("member not found")
	}
	if b.Status!=models.Available{
		return errors.New("book is not available")
	}
	b.Status = models.Borrowed
	// If Borrowed Books are maps
	// m.BorrowedBooks[b.ID] = *b
	m.BorrowedBooks = append(m.BorrowedBooks, *b)
	return nil
}

func (lib *Library) ReturnBook(bookID,memberID int) error{
	b,bexists := lib.Books[bookID]
	if !bexists{
		return errors.New("book not found")
	}
	m,mexists:=lib.Members[memberID]
	if !mexists{
		return errors.New("member not found")
	}
	if b.Status!=models.Borrowed{
		return errors.New("book is not borrowed")
	}
	b.Status = models.Available

	// O(1)? - comment: If req didn't specify the borrowedBooks to be slice (maps would be better choice)
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
	avail := []models.Book{}
	for _,b:= range lib.Books{
		if b.Status==models.Available{
			avail = append(avail,*b)
		}
	}
	return avail
}

func (lib *Library) ListBorrowedBooks(memberId int)[]models.Book{
	return lib.Members[memberId].BorrowedBooks
}

func (lib *Library) AddMember(member models.Member) error {
	if _,exists := lib.Members[member.ID]; exists{
		return errors.New("member already exists")
	}
	lib.Members[member.ID] = &member
	return nil
}

func (lib *Library) RemoveMember(memberId int) error{
	m,exists:=lib.Members[memberId]
	if!exists{
		return errors.New("member not found")
	}
	for _,book := range m.BorrowedBooks{
		lib.Books[book.ID].Status = models.Available
	}
	delete(lib.Members,memberId)
	return nil
}