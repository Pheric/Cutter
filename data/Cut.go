package data

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"sync"
)

type Cut struct {
	tableName     struct{} `sql:"cuts"`
	Uuid          string   `sql:",pk"`
	Client string
	Date          time.Time
	Price         float32
	Employees     string     `sql:"employees"`
	EmployeeSlice []Employee `sql:"-"`
}

func LoadCutWithUuid(uuid string) (Cut, error) {
	conn := openConnection()

	var c Cut
	err := conn.Model(&c).Where("uuid = ?", uuid).Select()
	if err != nil {
		err = fmt.Errorf("error selecting cut via uuid: %v", err)
		return Cut{}, err
	}

	// Make regex for employee string cleanup
	reg := regexp.MustCompile(`[{} ]`)
	// Replace every match with nothing
	c.Employees = reg.ReplaceAllString(c.Employees, "")
	// Split employee string
	split := strings.Split(c.Employees, ",")
	// Load employees and add them to the Cut
	for _, e_uuid := range split {
		if e_uuid == "" {
			continue
		}
		employee, err := LoadEmployeeWithUuid(e_uuid)
		if err != nil {
			log.Printf("Error loading employee %v for cut %v: %v\n", e_uuid, c.Uuid, err)
			continue
		}
		c.EmployeeSlice = append(c.EmployeeSlice, employee)
	}

	return c, err
}

func LoadCutsForClient(uuid string) []Cut {
	conn := openConnection()

	// Load all uuids
	type cId struct {
		tableName struct{} `sql:"cuts"`
		Uuid string
	}
	var uuids []cId
	conn.Model(&uuids).Where("client = ?", uuid).Limit(30).Select()
	if len(uuids) == 0 {
		return nil
	}

	// Load each Cut
	var cuts []Cut
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(uuids))
	for _, uuid := range uuids {
		go func(uuid string) {
			cut, err := LoadCutWithUuid(uuid)
			if err != nil {
				log.Printf("Error while loading cut with uuid %s: %v", uuid, err)
				return
			}

			lock.Lock()
			cuts = append(cuts, cut)
			lock.Unlock()
			wg.Done()
		}(uuid.Uuid)
	}
	wg.Wait()
	log.Println(len(cuts))

	return cuts
}

func (c Cut) Save() error {
	conn := openConnection()

	// Hackery required because the ORM has issues with arrays apparently
	// Make it look like an array
	c.Employees = "{"
	for _, e := range c.EmployeeSlice[:len(c.EmployeeSlice)-1] {
		c.Employees += fmt.Sprintf(`"%s", `, e.Uuid)
	}
	c.Employees += fmt.Sprintf(`"%s"}`, c.EmployeeSlice[len(c.EmployeeSlice)-1].Uuid)

	_, err := conn.Model(c).Insert()

	return err
}
