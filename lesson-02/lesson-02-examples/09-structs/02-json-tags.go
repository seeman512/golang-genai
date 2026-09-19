// Розділ 7.3 · Теги структур для JSON
//
// Демонструє: як struct tags керують зіставленням полів Go з
// ключами JSON, і опцію omitempty.
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Product struct {
	SKU     string  `json:"sku"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock,omitempty"`
}

func main() {
	// Unmarshal: JSON -> struct
	data := []byte(`{"sku":"A1","name":"Mug","price":9.99}`)
	var p Product
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("розпарсено: %+v\n", p)

	// Marshal: struct -> JSON
	out, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("назад у JSON:", string(out))
	// зверніть увагу: "in_stock" відсутній у виводі — InStock == false
	// і застосувалась опція omitempty
}
