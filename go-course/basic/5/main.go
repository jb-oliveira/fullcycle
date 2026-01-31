package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     string
	Cidade     string
	Estado     string
}

type Pessoa interface {
	Desativar()
}

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
	Endereco
	Address Endereco
}

func (c *Cliente) Desativar() {
	c.Ativo = false
}

func Desativacao(p Pessoa) {
	p.Desativar()
}

func main() {
	joao := Cliente{
		Nome:  "Joao",
		Idade: 48,
		Ativo: true,
	}
	joao.Cidade = "São Paulo"
	joao.Endereco.Cidade = "Recife"
	// joao.Desativar()
	Desativacao(&joao)

	fmt.Printf("Nome: %s, Idade: %d, Ativo: %t\n", joao.Nome, joao.Idade, joao.Ativo)
	fmt.Printf("Cliente: %v\n", joao)

}
