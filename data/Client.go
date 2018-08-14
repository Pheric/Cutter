package data

import (
	"fmt"
	"strings"
	"sync"
	"log"
)

const (
	WEEKLY = iota
	BIWEEKLY
	ON_DEMAND
)

type Client struct {
	tableName struct{} `sql:"clients"`
	Uuid      string   `sql:",pk"`
	First     string
	Last      string
	Address   string
	Phone     string
	Quote     float32
	Ttc       int
	Period    int // See period constants WEEKLY, BIWEEKLY, and ON_DEMAND
	Balance   float32
	Payments  []Payment `sql:"-"`
	Cuts      []Cut     `sql:"-"`
}

func LoadClientWithUuid(uuid string) (Client, error) {
	if uuid == "" {
		return MakeClient(), nil
	}

	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	c := new(Client)
	err := conn.Model(c).Where("uuid = ?", uuid).Select()
	if err != nil {
		err = fmt.Errorf("error while loading client with uuid %v: %v\n", uuid, err)
		return MakeClient(), err
	}

	c.Payments, err = LoadPaymentsForTarget(uuid)
	if err != nil {
		err = fmt.Errorf("error while loading payments for client with uuid %v: %v\n", uuid, err)
		return MakeClient(), err
	}

	c.Cuts = LoadCutsForClient(uuid)

	return *c, nil
}

func LoadAllClients() ([]Client, error) {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	type cId struct {
		tableName string `sql:"clients"`
		Uuid string
	}

	var ids []cId
	err := conn.Model(&ids).Select()
	if err != nil {
		err = fmt.Errorf("error while loading all clients' uuids: %v\n", err)
		return nil, err
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("error fetching clients: no clients found")
	}

	var clients []Client
	var wg sync.WaitGroup
	var lock = sync.Mutex{}
	wg.Add(len(ids))
	for _, u := range ids {
		go func(u string) {
			defer wg.Done()
			c, err := LoadClientWithUuid(u)
			if err != nil {
				log.Printf("Error fetching client: %v\n", err)
				return
			}
			lock.Lock()
			clients = append(clients, c)
			lock.Unlock()
		}(u.Uuid)
	}
	wg.Wait()

	return clients, nil
}

func MakeClient() Client {
	return Client {
		First: "John",
		Last: "Doe",
		Address: "123 Wallaby Way",
		Phone: "+1 (000) 000-0000",
		Quote: 30,
		Ttc: 15,
		Period: WEEKLY,
	}
}

func (c Client) SaveShallow() error {
	conn := openConnection()
	defer func(){
		if err := conn.Close(); err != nil {
			log.Printf("Error closing db connection: %v; **connection leak**\n", err)
		}
	}()

	_, err := conn.Model(&c).Returning("uuid").Insert()
	if err != nil && strings.Contains(err.Error(), "#23505") { // If row exists
		_, e := conn.Model(&c).Where("uuid = ?", c.Uuid).Update()
		err = e
	}

	return err
}
