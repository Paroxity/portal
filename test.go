package main

import "fmt"

type Entry struct {
	FirstField string
	SecondField int
}

type Test struct {
	entry []Entry
}

func main() {
	var entry Entry
	entries := make([]Entry, 2)
	for i := 0; i < 2; i++ {
		putValue(&entry, i)
		putValue2(&entry, "test")
		entries[i] = entry
	}

	fmt.Println(entries)
}

func putValue2(entry *Entry, value string) {
	entry.FirstField = value
}

func putValue(entry *Entry, value int) {
	entry.SecondField = value
}
