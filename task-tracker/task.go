package main

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Description string
	Done        bool
}
