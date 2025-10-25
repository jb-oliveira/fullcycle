package main

import "fmt"

type Cliente struct {
	Nome  string
	Ativo bool
}

func (c Cliente) Desativar() Cliente {
	c.Ativo = false
	return c
}

func teste() (string, bool) {
	return "Teste", true
}

func main() {

	joao := Cliente{
		Nome:  "João",
		Ativo: true,
	}
	joao = joao.Desativar()

	fmt.Printf("Cliente: %v\n", joao)

	var x interface{} = 10
	var y interface{} = "Teste"
	show(x)
	show(y)

	val, ok := y.(int)
	println(teste())
	println(ok)
	println(val)

}

func show(t interface{}) {
	fmt.Printf("O tipo é: %T e o valor é: %v\n", t, t)
}
