package data

import (
	"time"
	"log"
)

type Payment struct {
	tableName struct{} `sql:"payments"`
	Uuid      string   `sql:",pk"`
	Target    string   // uuid
	Amount    float32
	Date      time.Time
}

func LoadPaymentWithUuid(uuid string) (Payment, error) {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	payment := new(Payment)
	err := conn.Model(payment).Where("uuid = ?", uuid).Select()

	return *payment, err
}

func LoadPaymentsForTarget(uuid string) ([]Payment, error) {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	var payments []Payment
	err := conn.Model(&payments).Where("target = ?", uuid).Limit(30).Select()

	return payments, err
}

func (p Payment) SavePayment() error {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	return conn.Insert(&p)
}
