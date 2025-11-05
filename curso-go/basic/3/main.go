package main

import "fmt"

func main() {

	salaryMake := make(map[string]float64)
	salaryMake["Joao"] = 34.5
	salaryMake["Bosco"] = 4.5
	salaryMake["Oliveira"] = 78.5
	for k, v := range salaryMake {
		fmt.Printf("O salario de %s é de %f\n", k, v)
	}
	for _, v := range salaryMake {
		fmt.Printf("O salario é de %v\n", v)
	}
	for k := range salaryMake {
		fmt.Printf("O nome é de %v\n", k)
	}

	salary := map[string]float64{
		"João":  12.3,
		"Bosco": 34.6}

	fmt.Println(salary["João"])
	fmt.Println(salary["João22"])
	delete(salary, "João")
	fmt.Println(salary["João"])
	salary["Joao"] = 56.7
	fmt.Println(salary["Joao"])

}
