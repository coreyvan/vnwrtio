package main

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Page represents a page
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := fmt.Sprintf("pages/%s.txt", p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := fmt.Sprintf("pages/%s.txt", title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}
	return &Page{Title: title, Body: body}, nil
}
