package main

import (
	"encoding/json"
	"fmt"
	"github.com/lapingvino/readline"
	"os"
	"time"
)

func main() {
	lists := getLists()
	// Show lists that are to do
	// Option to distill list for those
	// Option to add new list
	newList(lists)
}

func getLists() []List {
	return []List{}
}

func newList(lists []List) {
	fmt.Println("Cool! We are going to learn new stuff!")
	category := chooseCategory()
	list := inputList()
	saveList(category, list, 0, lists)
	fmt.Println("Done! Doesn't that feel good?")
}

func chooseCategory() string {
	fmt.Println("Choose one from your existing categories or choose a new one:")
	var category string
	readline.Get(&category)
	return category
}

func inputList() []string {
	fmt.Println("Let's list the 25 things you want to learn now. If you want to stop early, that's fine; just leave the line empty. You can also go on for a bit more, but don't stretch it too far... it will probably not help very much...")
	list := []string{}
	for v, i := ".", 1; v != ""; list, i = append(list, v), i+1 {
		if i == 25 {
			fmt.Println("Now comes your 25th line. It's recommended that you stop here. You can stop by leaving the 26th line empty.")
		}
		fmt.Printf("%v: ", i)
		readline.Get(&v)
	}
	return list
}

type Goldlist struct {
	Lists []List `json:"lists"`
}

type List struct {
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
	Distill  int       `json:"book"`
	List     []string  `json:"list"`
}

func saveList(category string, list []string, book int, lists []List) {
	lists = append(lists, List{time.Now(), category, book, list})
	v, err := json.MarshalIndent(Goldlist{lists}, "", "  ")
	if err != nil {
		panic("JSON marchal failed")
	}
	file, err := os.Create("goldlist.json")
	if err != nil {
		panic("Failed to open goldlist.json and failed to create it")
	}
	defer file.Close()
	fmt.Fprintln(file, string(v))
}
