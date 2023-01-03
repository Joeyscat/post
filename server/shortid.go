package main

import "github.com/teris-io/shortid"

var sid *shortid.Shortid

func init() {
	var err error
	sid, err = shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		panic(err)
	}
}

func shortId() (string, error) {
	return sid.Generate()
}
