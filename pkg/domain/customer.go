package domain

type Customer struct {
	ID       int
	Name     string
	Location *Coordinate
}

func NewCustomer(id int, name string, location *Coordinate) Customer {
	return Customer{ID: id, Name: name, Location: location}
}

func (c Customer) WithLocation(location *Coordinate) Customer {
	return NewCustomer(c.ID, c.Name, location)
}

type Customers []Customer
