// Example to demonstrate how to get realitme inverter data
package main

import (
	"github.com/brutella/geizhals"

	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		id = flag.String("id", "", "Product ID")
	)
	flag.Parse()

	if len(*id) == 0 {
		flag.PrintDefaults()
		return
	}

	pr, err := geizhals.GetProduct(*id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pairs := []string{
		fmt.Sprintf("price=%.2f", pr.MinPrice),
	}
	fmt.Printf("geizhals,id=%s %s", *id, strings.Join(pairs, ","))
	fmt.Println()
}
