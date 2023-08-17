package main

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

// time.Parse layout constants
const (
	longForm  = "Jan 2, 2006 at 3:04pm (MST)"
	shortForm = "2006-Jan-02"
)

func TestGORM(t *testing.T) {
	t.Skip()
	var pets []*Pet
	var toys []Toy
	var team []User
	var languages []Language
	var friends []*User

	var user1 User
	var manager1 User = User{
		Name: "Faras",
	}

	bday, _ := time.Parse(shortForm, "1989-Jan-28")
	// companyID := 10
	company := Company{
		Name: "TechieWorld",
	}
	accUserID := sql.NullInt64{
		Int64: 10,
		Valid: true,
	}

	for i := 0; i < 2; i++ {
		toy := Toy{
			Name:      "Baseball" + fmt.Sprintf("%d", i),
			OwnerID:   "Rikhu" + fmt.Sprintf("%d", i),
			OwnerType: "Cow" + fmt.Sprintf("%d", i),
		}
		pet := Pet{
			UserID: &user1.ID,
			Name:   "Rikhu" + fmt.Sprintf("%d", i),
			Toy:    toy,
		}
		language := Language{
			Code: "FC" + fmt.Sprintf("%d", i),
			Name: "XX" + fmt.Sprintf("%d", i),
		}
		user1.Name = "User" + fmt.Sprint(i)
		pets = append(pets, &pet)
		toys = append(toys, toy)
		team = append(team, user1)
		friends = append(friends, &user1)
		languages = append(languages, language)
	}

	user1 = User{
		Name:     "Faras Mohan Dewal",
		Age:      34,
		Birthday: &bday,
		Account: Account{
			UserID: accUserID,
		},
		Pets:      pets,
		Toys:      toys,
		CompanyID: &company.ID,
		Company:   company,
		ManagerID: &manager1.ID,
		Manager:   &manager1,
		Team:      team,
		Languages: languages,
		Friends:   friends,
		Active:    true,
	}

	DB.Create(&user1)

	var userResult User
	if err := DB.First(&userResult, user1.ID).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	} else {
		fmt.Println(" - Created User : ",
			userResult.ID, userResult.Active, userResult.Age,
			userResult.Name, userResult.Pets, userResult.Toys)
	}

	if err := DB.Find(&userResult, user1.ID).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	} else {
		fmt.Println(" - Created User : ",
			userResult.ID, userResult.Active, userResult.Age,
			userResult.Name, userResult.Pets, userResult.Toys,
			userResult.Friends, userResult.Languages, userResult.Manager,
			&userResult.ManagerID, userResult.Team)
	}
}

func TestGormDocs(t *testing.T) {
	t.Skip()
	bday, _ := time.Parse(shortForm, "1989-Jul-19")
	user := User{Name: "Faras", Age: 34, Birthday: &bday}

	// Create single entry
	result := DB.Create(&user)
	fmt.Println("Create single entry: ", user.ID, result.Error, result.RowsAffected)

	// Create multiple entries
	users := []*User{
		{Name: "FarasM"},
		{Name: "FarasMD"},
	}
	result = DB.Create(users)
	fmt.Println("Create multiple entries: ", result.Error, result.RowsAffected)
	for _, user := range users {
		fmt.Print(user.ID, " : ")
	}

	// Create multiple entries
	usersVal := []User{
		{Name: "FarasM"},
		{Name: "FarasMD"},
	}
	result = DB.Create(&usersVal)
	fmt.Println("Create multiple entries: ", result.Error, result.RowsAffected)
	for _, user := range users {
		fmt.Print(user.ID, " : ")
	}

	// selected fields used for create
	user = User{Name: "Faras", Age: 34, Birthday: &bday}
	result = DB.Select("Name", "Age", "CreatedAt").Create(&user)
	fmt.Println("Create entry with selected fields: ", result.Error, result.RowsAffected)
	fmt.Print(user.ID, " : ")

	// selected fields omitted for create
	user = User{Name: "Faras", Age: 34, Birthday: &bday}
	result = DB.Omit("Name", "Age", "CreatedAt").Create(&user)
	fmt.Println("Create entries with omitted fields: ", result.Error, result.RowsAffected)
	fmt.Print(user.ID, " : ")
}

func TestCreateWithAssociations(t *testing.T) {
	t.Skip()
	// Create with associations
	// INSERT INTO `users` ...
	// INSERT INTO `credit_cards` ...
	result := DB.Create(&User{
		Name:       "FarasMD",
		CreditCard: CreditCard{Number: "411111111111"},
	})
	fmt.Println(result.Error, result.RowsAffected)

	// Skip saving credit card details
	// INSERT INTO `toys`
	// INSERT INTO `users`
	toy := Toy{
		Name:      "Baseball",
		OwnerID:   "Rikhu",
		OwnerType: "Cow",
	}
	user := User{
		Name:       "FarasMD",
		CreditCard: CreditCard{Number: "411111111111"},
		Toys:       []Toy{toy},
	}
	result = DB.Omit("CreditCard").Create(&user)
	fmt.Println(result.Error, result.RowsAffected)

	// skip all associations
	// INSERT INTO `users`
	user = User{
		Name:       "FarasMD",
		CreditCard: CreditCard{Number: "411111111111"},
		Toys:       []Toy{toy},
	}
	DB.Omit(clause.Associations).Create(&user)
}

func TestDefaultValue(t *testing.T) {
	t.Skip()
	// selected fields omitted for create
	bday, _ := time.Parse(shortForm, "1989-Jul-19")
	user := User{Name: "Faras", Age: 34, Birthday: &bday}
	result := DB.Omit("Name", "Age", "CreatedAt").Create(&user)
	fmt.Println("Create entries with omitted fields: ", result.Error, result.RowsAffected)
	fmt.Println(user.ID, " : ", user.Name)
	// user.Name still has the populated entry
	// read from database reveals the default name (as name was omitted during creation)
	result = DB.Find(&user, user.ID)
	fmt.Println("Find user record: ", user.ID, " : ", user.Name)
	fmt.Println(result.Error, result.RowsAffected)
}

func TestUpsert(t *testing.T) {
	t.Skip()
	user := User{
		Name:  "FarasMD",
		Age:   34,
		Role:  "Owner",
		Count: 4,
	}
	// 2023/08/17 10:36:57 /home/dewal/workspace/playground/main_test.go:218
	// [136.769ms] [rows:1] INSERT INTO `users`
	// (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,
	// `birthday`,`company_id`,`manager_id`,`active`,`role`)
	// VALUES ('2023-08-17 10:36:57.096','2023-08-17 10:36:57.096',NULL,'FarasMD',34,
	// NULL,NULL,NULL,false,'Owner')
	DB.Create(&user)

	// Error 1062 (23000): Duplicate entry '1' for key 'users.PRIMARY'
	DB.Create(&user)

	// 2023/08/17 10:36:57 /home/dewal/workspace/playground/main_test.go:230
	// [2.401ms] [rows:0] INSERT INTO `users`
	// (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,
	// `birthday`,`company_id`,`manager_id`,`active`,`role`,`id`)
	// VALUES ('2023-08-17 10:36:57.096','2023-08-17 10:36:57.096',NULL,'FarasMD',34,
	// NULL,NULL,NULL,false,'Owner',1) ON DUPLICATE KEY UPDATE `id`=`id`
	// Do nothing on conflict
	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

	// Update columns to default value on `id` conflict
	// 2023/08/17 10:36:57 /home/dewal/workspace/playground/main_test.go:245
	// [85.946ms] [rows:2] INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,
	// `birthday`,`company_id`,`manager_id`,`active`,`role`,`id`)
	// VALUES ('2023-08-17 10:36:57.096','2023-08-17 10:36:57.096',NULL,'FarasMD',34,
	// NULL,NULL,NULL,false,'Owner',1) ON DUPLICATE KEY UPDATE `role`='user'
	// mysql> select * from users;
	// +----+-------------------------+-------------------------+------------+---------+------+----------+------------+------------+--------+------+
	// | id | created_at              | updated_at              | deleted_at | name    | age  | birthday | company_id | manager_id | active | role |
	// +----+-------------------------+-------------------------+------------+---------+------+----------+------------+------------+--------+------+
	// |  1 | 2023-08-17 10:36:57.096 | 2023-08-17 10:36:57.096 | NULL       | FarasMD |   34 | NULL     |       NULL |       NULL |      0 | user |
	// +----+-------------------------+-------------------------+------------+---------+------+----------+------------+------------+--------+------+
	// 1 row in set (0.00 sec)
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
	}).Create(&user)

	// Use SQL expression
	user.Count = 1
	// Since new struct value of Count (=1) less than db entry for Count (=4), count is not updated
	DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(
			map[string]interface{}{
				"count": gorm.Expr("GREATEST(count, VALUES(count))"),
			}),
	}).Create(&user)

	// Update columns to new value on `id` conflict
	// Only name and age updated upon conflict
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
	}).Create(&user)

	// Update all columns to new value on conflict except primary keys and those columns having default values from sql func
	// All fields updated to struct value. Role changed back to value of struct, i.e. 'Owner'
	// Count will be updated to the latest in struct
	DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)
}

func TestRead(t *testing.T) {
	{
		var user User
		// Get the first record ordered by primary key
		result := DB.First(&user)
		// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
		fmt.Println(result.RowsAffected, result.Error)
		fmt.Println(user.ID, user.Name, user.Birthday, user.Pets)
		// check error ErrRecordNotFound
		fmt.Println(errors.Is(result.Error, gorm.ErrRecordNotFound))
	}
	{
		var user User
		// Get one record, no specified order
		result := DB.Take(&user)
		// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 1
		fmt.Println(result.RowsAffected, result.Error)
		fmt.Println(user.ID, user.Name, user.Birthday, user.Pets)
		// check error ErrRecordNotFound
		fmt.Println(errors.Is(result.Error, gorm.ErrRecordNotFound))
	}
	{
		var user User
		// Get last record, ordered by primary key desc
		result := DB.Last(&user)
		// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
		fmt.Println(result.RowsAffected, result.Error)
		fmt.Println(user.ID, user.Name, user.Birthday, user.Pets)
		// check error ErrRecordNotFound
		fmt.Println(errors.Is(result.Error, gorm.ErrRecordNotFound))
	}
	{
		// doesn't work
		user := map[string]interface{}{}
		result := DB.Table("users").First(&user)
		// SELECT * FROM `users` ORDER BY `users`. LIMIT 1
		fmt.Println(result.Error, result.RowsAffected) // model value required 0
		fmt.Println(user)
	}
	{
		// works with Take
		user := map[string]interface{}{}
		result := DB.Table("users").Take(&user)
		// SELECT * FROM `users` LIMIT 1
		fmt.Println(result.Error, result.RowsAffected)
		fmt.Println(user)
	}
	{
		var users []User
		result := DB.Find(&users, []int{1, 2, 3})
		// SELECT * FROM `users` WHERE `users`.`id` IN (1,2,3) AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
		fmt.Println(users)
	}
	{
		// Get all records
		var users []User
		result := DB.Find(&users)
		// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL

		fmt.Println(result.Error, result.RowsAffected)
		for index, user := range users {
			fmt.Println(index, user.ID, user.Name)
		}
	}

	{
		var user User // Get first matched record
		result := DB.Where("name = ?", "hidden").First(&user)
		// SELECT * FROM `users` WHERE name = 'hidden'
		// AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // Get all matched records
		result := DB.Where("name <> ?", "Faras").Find(&users)
		// SELECT * FROM `users` WHERE name <> 'Faras' AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // IN
		result := DB.Where("name IN ?", []string{"Faras", "FarasMD"}).Find(&users)
		// SELECT * FROM `users` WHERE name IN ('Faras','FarasMD') AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // LIKE
		result := DB.Where("name LIKE ?", "%ra%").Find(&users)
		// SELECT * FROM `users` WHERE name LIKE '%ra%' AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // AND
		result := DB.Where("name = ? AND age >= ?", "Faras", "22").Find(&users)
		// SELECT * FROM `users` WHERE (name = 'Faras' AND age >= '22') AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // Time
		earlierTime, _ := time.Parse(shortForm, "2023-Aug-17")
		result := DB.Where("updated_at > ?", earlierTime).Find(&users)
		// SELECT * FROM `users` WHERE updated_at > '2023-08-17 05:30:00' AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var users []User // BETWEEN
		earlierTime, _ := time.Parse(shortForm, "2023-Aug-17")
		result := DB.Where("created_at BETWEEN ? AND ?", earlierTime, time.Now()).Find(&users)
		// SELECT * FROM `users` WHERE (created_at BETWEEN '2023-08-17 05:30:00'
		// AND '2023-08-17 15:27:10.044') AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected)
	}
	{
		var user User
		user.ID = 10
		result := DB.Where("id = ?", 20).First(&user)
		// SELECT * FROM `users` WHERE id = 20 AND `users`.`deleted_at` IS NULL
		//               AND `users`.`id` = 10 ORDER BY `users`.`id` LIMIT 1
		fmt.Println(result.Error, result.RowsAffected) // record not found 0
		fmt.Println(user.ID)                           // 10

	}
	{
		var user User //Struct
		result := DB.Where(&User{Name: "Faras", Age: 20}).First(&user)
		// SELECT * FROM users WHERE name = "Faras" AND age = 20 ORDER BY id LIMIT 1;
		fmt.Println(result.Error, result.RowsAffected) // record not found 0
		fmt.Println(user.ID)                           // 10
	}
	{
		var users []User // Map
		result := DB.Where(map[string]interface{}{"name": "Faras", "age": 20}).Find(&users)
		// SELECT * FROM `users` WHERE `age` = 20 AND `name` = 'Faras' AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected) // <nil> 0
		fmt.Println(users)                             // []
	}
	{
		var users []User // Slice of primary keys
		result := DB.Where([]int64{20, 21, 22}).Find(&users)
		// SELECT * FROM `users` WHERE `users`.`id` IN (20,21,22) AND `users`.`deleted_at` IS NULL
		fmt.Println(result.Error, result.RowsAffected) // <nil> 3

		for index, user := range users {
			fmt.Println(index, user.ID, user.Name)
		}
		/*
			0 20 User1
			1 21 Faras
			2 22 FarasM
		*/
	}
	// 	rows, err := DB.Table("orders").Select("date(created_at) as date, sum(amount) as total")
	// .Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
}
