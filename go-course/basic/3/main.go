package main

import "fmt"

func main() {

	salaryMake := make(map[string]float64)
	salaryMake["Joao"] = 34.5
	salaryMake["Bosco"] = 4.5
	salaryMake["Oliveira"] = 78.5
	for k, v := range salaryMake {
		fmt.Printf("The salary of %s is %f\n", k, v)
	}
	for _, v := range salaryMake {
		fmt.Printf("The salary is %v\n", v)
	}
	for k := range salaryMake {
		fmt.Printf("The name is %v\n", k)
	}

	salary := map[string]float64{
		"John":  12.3,
		"Bosco": 34.6}

	fmt.Println(salary["John"])
	fmt.Println(salary["John22"])
	delete(salary, "John")
	fmt.Println(salary["John"])
	salary["Joao"] = 56.7
	fmt.Println(salary["Joao"])

}
