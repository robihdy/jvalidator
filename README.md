# jvalidator

jvalidator: JSON validator for Go (Golang)

[![GoDoc](https://godoc.org/github.com/robihid/jvalidator?status.svg)](https://godoc.org/github.com/robihid/jvalidator)

## Index

* [Getting Started](#getting-started)
* [Example](#example)

## Getting Started

#### Download

```shell
go get -u github.com/robihid/jvalidator
```

## Example
```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robihid/jvalidator"
)

type Person struct {
	Name  string
	Email string
	Age   int
}

var rawJSON = []byte(`{
	"Name": "Robi Hidayat",
	"Email": "invalid.mail.com",
	"Age": 20
}`)

func main() {
	http.HandleFunc("/", addPerson)
	http.ListenAndServe("localhost:8080", nil)
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	var p Person

	validator, _ := jvalidator.New(rawJSON, &p) // ignore error

	validator.Required("Name", "Email", "Age")
	validator.Email("Email")

	if !validator.IsValid() {
		invalid, _ := validator.ToJSON() // ignore error
		w.Header().Set("Content-type", "text/json")
		fmt.Fprint(w, string(invalid)) // {"Email":["Must contain a valid email address."]}
		return
	}

	b, _ := json.Marshal(p) // ignore error
	w.Header().Set("Content-type", "text/json")
	fmt.Fprint(w, string(b))
}
```