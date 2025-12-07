#!/bin/bash

# Olhar a documentação em https://github.com/swaggo/swag?tab=readme-ov-file#getting-started
# Para gerar a documentação do projeto

# Inicialização
# swag init --parseDependency --dir ./ --output ./docs # gerado pelo kiro
swag init -g cmd/server/main.go # curso