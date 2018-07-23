package data

import "time"

type Payment struct {
	tableName struct{} `sql:"payments"`
	Uuid      string   `sql:",pk"`
	Target    string   // uuid
	Amount    float32
	Date      time.Time
}

func LoadPaymentWithUuid(uuid string) (Payment, error) {
	conn := openConnection()

	payment := new(Payment)
	err := conn.Model(payment).Where("uuid = ?", uuid).Select()

	return *payment, err
}

func LoadPaymentsForTarget(uuid string) ([]Payment, error) {
	conn := openConnection()

	var payments []Payment
	err := conn.Model(&payments).Where("target = ?", uuid).Limit(30).Select()

	return payments, err
}

func (p Payment) SavePayment() error {
	conn := openConnection()

	return conn.Insert(p)
}
