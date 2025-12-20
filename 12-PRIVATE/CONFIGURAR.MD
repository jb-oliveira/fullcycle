# Como  o go mod acessar o github

1. git config --global url."git@github.com:".insteadOf "https://github.com/"
2. go env -w GOPRIVATE="github.com*"



