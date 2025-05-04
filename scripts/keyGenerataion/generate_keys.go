package main

import (
	"fmt"
	"kondangin-backend/internal/utils"
)

func main() {
	// Panggil fungsi untuk generate keys
	err := utils.GenerateKeys()
	if err != nil {
		fmt.Println("Error generating keys:", err)
		return
	}
	fmt.Println("Keys generated successfully and saved to keys/ directory")
}
