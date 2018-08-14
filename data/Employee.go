package data

import (
	"strings"
	"log"
)

type Employee struct {
	tableName struct{} `sql:"Employees"`
	Uuid      string   `sql:",pk"`
	First     string
	Last      string
	Phone     string
	Owed      float32
	Payments  []Payment `sql:"-"`
}

func LoadEmployeeWithUuid(uuid string) (Employee, error) {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	e := new(Employee)
	err := conn.Model(e).Where("uuid = ?", uuid).Select()
	if err != nil {
		return Employee{}, err
	}

	payments, err := LoadPaymentsForTarget(uuid)
	if err != nil {
		return *e, err
	}
	e.Payments = payments

	return *e, nil
}

// Selects only the first employee with this name
// Could later be updated/extended by returning a []*Employee
func LoadEmployeeByName(last, first string) (*Employee, error) {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	e := new(Employee)
	err := conn.Model(e).Where("first = ?", first).Where("last = ?", last).Limit(1).Select()

	return e, err
}

func LoadAllEmployees() ([]Employee, error) {
	conn := openConnection()

	var employees []Employee
	err := conn.Model(&employees).Select()

	return employees, err
}

func (e Employee) Save() error {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	_, err := conn.Model(e).Insert()
	if err != nil && strings.Contains(err.Error(), "#23505") { // If row exists
		_, e := conn.Model(e).Where("uuid = ?", e.Uuid).Update()
		err = e
	}

	return err
}
