package models

type State int
const (
	Available State = iota
	Borrowed 
)
type Book struct{
	ID int
	Title string
	Author string
	Status State
}
