package main

import (
	"fmt"
	"testing"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	var user1 User
	var pet1 Pet
	var toy1 Toy

	toy1 = Toy{
		Name:      "Baseball",
		OwnerID:   "Rikhu",
		OwnerType: "Dog",
	}

	pet1 = Pet{
		UserID: &user1.ID,
		Name:   "Rikhu",
		Toy:    toy1,
	}
	pets := []*Pet{&pet1}
	toys := []Toy{toy1}
	//bday, _ := time.Parse("1989-Jan-28", "1989-Jan-28")
	//fmt.Println(bday)
	user1 = User{
		Name: "Faras Mohan Dewal",
		Age:  34,
		//	Birthday: &bday,
		Active: true,
		Pets:   pets,
		Toys:   toys,
	}

	DB.Create(&user1)

	var result User
	if err := DB.First(&result, user1.ID).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	} else {
		fmt.Println(" - Created User : ",
			result.ID, result.Active, result.Age, result.Name, result.Pets, result.Toys)
	}
}
