package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"strings"
)

func createDB() {
	database, err := sql.Open("sqlite3", "./finance.db")
	checkErr(err)
	statements := "CREATE TABLE IF NOT EXISTS category(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT);CREATE TABLE IF NOT EXISTS payment(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT, price DOUBLE, timePeriod TEXT, typeOfTransaction TEXT, comments TEXT, categoryID INTEGER NOT NULL, FOREIGN KEY(categoryID) REFERENCES category(id));"
	_, err = database.Exec(statements)
	checkErr(err)
	database.Close()
}
func controlMenu() int {
	var option int
	fmt.Println("1. Управление категориями")
	fmt.Println("2. Управление платежами ")
	fmt.Println("3. Завершить ")
	fmt.Print("Выбор>")
	_, err := fmt.Scan(&option)
	checkErr(err)
	return option
}
func starter(option int) {
	var operation int
	if option == 1 {
		fmt.Println("1. Создать")
		fmt.Println("2. Прочитать")
		fmt.Println("3 Изменить")
		fmt.Println("4. Удалить")
		_, err := fmt.Scan(&operation)
		checkErr(err)
		categoryOP(operation)
	} else {
		fmt.Println("1. Создать")
		fmt.Println("2. Прочитать")
		fmt.Println("3 Изменить")
		fmt.Println("4. Удалить")
		_, err := fmt.Scan(&operation)
		checkErr(err)
		paymentOP(operation)
	}
}
func allCategory() {
	database, err := sql.Open("sqlite3", "./finance.db")
	checkErr(err)
	fmt.Println("Все категории:")
	rows, err := database.Query("SELECT id, name FROM category")
	checkErr(err)
	var id int
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		checkErr(err)
		fmt.Println(strconv.Itoa(id) + ". " + name)
	}
	database.Close()
}
func allPayments() {
	database, err := sql.Open("sqlite3", "./finance.db")
	checkErr(err)
	fmt.Println("Все платежи:")
	rows, err := database.Query("SELECT payment.id,payment.name, payment.price, payment.timePeriod, payment.typeOfTransaction,payment.comments, category.name FROM category JOIN payment ON category.id = payment.categoryID")
	checkErr(err)
	var id int
	var name string
	var price string
	var time string
	var typeOfPayment string
	var comments string
	var categoryName string
	for rows.Next() {
		err := rows.Scan(&id, &name, &price, &time, &typeOfPayment, &comments, &categoryName)
		checkErr(err)
		fmt.Println(strconv.Itoa(id) + ". " + name + " " + price + " " + time + " " + typeOfPayment + " " + comments + " " + categoryName)
	}
}
func categoryOP(operation int) {
	database, err := sql.Open("sqlite3", "./finance.db")
	checkErr(err)
	if operation == 1 {
		var name string
		fmt.Print("Название категории>")
		_, err = fmt.Scan(&name)
		checkErr(err)
		_, err := database.Exec("INSERT INTO category(name) VALUES(?) ", name)
		checkErr(err)

	} else if operation == 2 {
		allCategory()
	} else if operation == 3 {
		var id int
		var name string
		fmt.Println("Изменить название категории:")
		allCategory()
		fmt.Print("ID категории>")
		_, err = fmt.Scan(&id)
		checkErr(err)
		fmt.Print("Новое название категории>")
		_, err = fmt.Scan(&name)
		checkErr(err)
		_, err = database.Exec("UPDATE category SET name = ? WHERE id = ?", name, id)
		checkErr(err)

	} else if operation == 4 {
		var id int
		fmt.Println("Удалить категорию:")
		allCategory()
		fmt.Print("ID категории>")
		_, err = fmt.Scan(&id)
		checkErr(err)
		_, err = database.Exec("DELETE FROM category WHERE id = ?", id)
		checkErr(err)

	}
	database.Close()
	main()
}
func paymentOP(operation int) {
	database, err := sql.Open("sqlite3", "./finance.db")
	checkErr(err)
	if operation == 1 {
		var name string
		var price float32
		var time string
		var typeOfPayment string
		var comments string
		var categoryID int

		fmt.Println("Название платежа>")
		_, err = fmt.Scan(&name)
		checkErr(err)
		fmt.Println("Сумма платежа>")
		_, err = fmt.Scan(&price)
		checkErr(err)
		fmt.Println("Период платежа(yyyy-mm-dd):")
		_, err = fmt.Scan(&time)
		checkErr(err)
		fmt.Println("Тип платежа:")
		fmt.Println("1. Доходы")
		fmt.Println("2. Расходы")
		fmt.Print(">")
		_, err = fmt.Scan(&typeOfPayment)
		typeOfPayment = strings.TrimSpace(typeOfPayment)
		if typeOfPayment == "1" {
			typeOfPayment = "Доход"
		} else {
			typeOfPayment = "Расход"
		}
		fmt.Println("Комментарий платежа>")
		_, err = fmt.Scan(&comments)
		checkErr(err)
		fmt.Println("Категория платежа:")
		allCategory()
		fmt.Println("1. Существующий")
		fmt.Println("2. Новый")
		_, err = fmt.Scan(&categoryID)
		checkErr(err)
		if categoryID == 1 {
			fmt.Println("ID категории>")
			_, err = fmt.Scan(&categoryID)
			checkErr(err)
		} else {
			var name string
			fmt.Print("Название категории>")
			_, err = fmt.Scan(&name)
			checkErr(err)
			_, err := database.Exec("INSERT INTO category(name) VALUES(?) ", name)
			checkErr(err)
			row := database.QueryRow("SELECT id FROM category WHERE name = ?", name)
			err = row.Scan(&categoryID)
			checkErr(err)
		}
		_, err := database.Exec("INSERT INTO payment(name, price, timePeriod, typeOfTransaction, comments, categoryID) VALUES(?,?,?,?,?,?) ", name, price, time, typeOfPayment, comments, categoryID)
		checkErr(err)
	} else if operation == 2 {
		var option int
		fmt.Println("1. Все платежи\n2. Доходы\n3. Расходы")
		_, err = fmt.Scan(&option)
		checkErr(err)
		if option == 1 {
			allPayments()
		} else if option == 2 {
			var sum float64 = 0
			fmt.Println("Доходы:")
			rows, err := database.Query("SELECT payment.id,payment.name, payment.price, payment.timePeriod, payment.typeOfTransaction,payment.comments, category.name FROM category JOIN payment ON category.id = payment.categoryID WHERE payment.typeOfTransaction =?", "Доход")
			checkErr(err)
			var id int
			var name string
			var price string
			var time string
			var typeOfPayment string
			var comments string
			var categoryName string
			for rows.Next() {
				err := rows.Scan(&id, &name, &price, &time, &typeOfPayment, &comments, &categoryName)
				checkErr(err)
				num, err := strconv.ParseFloat(price, 64)
				sum += num
				fmt.Println(strconv.Itoa(id) + ". " + name + " " + price + " " + time + " " + typeOfPayment + " " + comments + " " + categoryName)
			}
			result := fmt.Sprintf("%f", sum)
			fmt.Println("Сумма доходов: " + result)
		} else {
			var sum float64 = 0
			fmt.Println("Расходы:")
			rows, err := database.Query("SELECT payment.id,payment.name, payment.price, payment.timePeriod, payment.typeOfTransaction,payment.comments, category.name FROM category JOIN payment ON category.id = payment.categoryID WHERE payment.typeOfTransaction =?", "Расход")
			checkErr(err)
			var id int
			var name string
			var price string
			var time string
			var typeOfPayment string
			var comments string
			var categoryName string
			for rows.Next() {
				err := rows.Scan(&id, &name, &price, &time, &typeOfPayment, &comments, &categoryName)
				checkErr(err)
				num, err := strconv.ParseFloat(price, 64)
				sum += num
				fmt.Println(strconv.Itoa(id) + ". " + name + " " + price + " " + time + " " + typeOfPayment + " " + comments + " " + categoryName)
			}
			result := fmt.Sprintf("%f", sum)
			fmt.Println("Сумма расходов: " + result)

		}

	} else if operation == 3 {
		var option int
		var id int
		allPayments()
		fmt.Println("ID платежа>")
		_, err = fmt.Scan(&id)
		checkErr(err)
		fmt.Println("1. Название\n2.Сумму \n3.Время \n4.Тип \n5.Комментарий \n6.категорию")
		_, err = fmt.Scan(&option)
		checkErr(err)
		if option == 1 {
			var name string
			fmt.Print("Название>")
			_, err = fmt.Scan(&name)
			checkErr(err)
			_, err = database.Exec("UPDATE payment SET name = ? WHERE id = ?", name, id)
			checkErr(err)
		} else if option == 2 {
			var price string
			fmt.Print("Сумма>")
			_, err = fmt.Scan(&price)
			_, err = database.Exec("UPDATE payment SET price = ? WHERE id = ?", price, id)
			checkErr(err)
		} else if option == 3 {
			var time string
			fmt.Print("Время>")
			_, err = fmt.Scan(&time)
			_, err = database.Exec("UPDATE payment SET timePeriod = ? WHERE id = ?", time, id)
			checkErr(err)
		} else if option == 4 {
			var typeOfPayment string
			fmt.Println("Тип платежа:")
			fmt.Println("1. Доходы")
			fmt.Println("2. Расходы")
			fmt.Print(">")
			_, err = fmt.Scan(&typeOfPayment)
			typeOfPayment = strings.TrimSpace(typeOfPayment)
			if typeOfPayment == "1" {
				typeOfPayment = "Доход"
			} else {
				typeOfPayment = "Расход"
			}
			_, err = database.Exec("UPDATE payment SET typeOfTransaction = ? WHERE id = ?", typeOfPayment, id)
			checkErr(err)
		} else if option == 5 {
			var comment string
			fmt.Print("Комментарий>")
			_, err = fmt.Scan(&comment)
			_, err = database.Exec("UPDATE payment SET comments = ? WHERE id = ?", comment, id)
			checkErr(err)
		} else if option == 6 {
			var categoryID int
			fmt.Println("Категория платежа:")
			allCategory()
			fmt.Println("1. Существующий")
			fmt.Println("2. Новый")
			_, err = fmt.Scan(&categoryID)
			checkErr(err)
			if categoryID == 1 {
				fmt.Println("ID категории>")
				_, err = fmt.Scan(&categoryID)
				checkErr(err)
			} else {
				var name string
				fmt.Print("Название категории>")
				_, err = fmt.Scan(&name)
				checkErr(err)
				_, err := database.Exec("INSERT INTO category(name) VALUES(?) ", name)
				checkErr(err)
				row := database.QueryRow("SELECT id FROM category WHERE name = ?", name)
				err = row.Scan(&categoryID)
				checkErr(err)
			}
			_, err = database.Exec("UPDATE payment SET categoryID = ? WHERE id = ?", categoryID, id)
			checkErr(err)
		}

	} else if operation == 4 {
		var id int
		allPayments()
		fmt.Print("ID платежа>")
		_, err = fmt.Scan(&id)
		checkErr(err)
		_, err = database.Exec("DELETE FROM payment WHERE id = ?", id)
		checkErr(err)
	}
	database.Close()
	main()
}
func main() {
	createDB()
	option := controlMenu()
	if option == 3 {
		os.Exit(3)
	}
	starter(option)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
