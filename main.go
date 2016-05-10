package main

import (
	"encoding/json"
	"fmt"
	"github.com/lapingvino/readline"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	lists := getLists()
	categories := countCategories(lists)
	datelists := byDate(lists)
	todolists, nodo := filter(datelists)
	if len(todolists) > 0 {
		// Show lists that are to do and distill them
		lists = distill(todolists, nodo)
	} else {
		newList(lists, categories)
	}
}

func getLists() []List {
	js, err := ioutil.ReadFile("goldlist.json")
	if err != nil {
		return []List{}
	}
	var goldlists Goldlist
	json.Unmarshal(js, &goldlists)
	return goldlists.Lists
}

func byDate(lists []List) map[time.Time]List {
	datelists := map[time.Time]List{}
	for _, l := range lists {
		datelists[l.Date] = l
	}
	return datelists
}

type mydate []time.Time

func (t mydate) Len() int {
	return len(t)
}

func filter(datelists map[time.Time]List) (now []List, nodo []List) {
	now = []List{}
	nodo = []List{}

	for _, l := range datelists {
		if time.Now().Unix() > l.Date.AddDate(0, 0, 14).Unix() && l.Distill > -1 {
			now = append(now, l)
		} else {
			nodo = append(nodo, l)
		}
	}
	return now, nodo
}

func countCategories(lists []List) (categories map[string]int) {
	categories = map[string]int{}
	for _, v := range lists {
		categories[v.Category]++
	}
	return categories
}

func distill(todo []List, nodo []List) []List {
	return append(todo, nodo...)
}

func newList(lists []List, categories map[string]int) {
	fmt.Println("Cool! We are going to learn new stuff!")
	category := chooseCategory(categories)
	list := inputList()
	if len(list) > 0 {
		saveList(category, list, 0, lists)
		fmt.Println("Done! Doesn't that feel good?")
	} else {
		fmt.Println("Well then... Maybe we add some words another day?")
	}
}

func chooseCategory(categories map[string]int) string {
	fmt.Println("Choose one from your existing categories or choose a new one:")
	for cat, n := range categories {
		fmt.Printf(" %v (%v entries)\n", cat, n)
	}
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
	return list[:len(list)-1]
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
