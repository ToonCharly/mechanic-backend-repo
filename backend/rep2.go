package main

import (
	"os"
	"strings"
)

func main() {
	b, _ := os.ReadFile("internal/application/usecases/service_usecase.go")
	str := string(b)
	str = strings.Replace(str, "Recibo r\xe1pido", "Recibo rapido", -1)
	str = strings.Replace(str, "rápido", "rapido", -1)
	os.WriteFile("internal/application/usecases/service_usecase.go", []byte(str), 0644)
}
