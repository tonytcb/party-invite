package domain

type Customer struct {
	ID       int
	Name     string
	Location *Coordinate
}

func NewCustomer(id int, name string, location *Coordinate) Customer {
	return Customer{ID: id, Name: name, Location: location}
}

type Customers []Customer
