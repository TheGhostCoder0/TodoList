package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type listItem struct {
	id   int
	desc string
}

type todoList struct {
	name     string
	numItems int
	items    []listItem
}

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createNewList(filename string, listName string) {
	data := []byte(listName + "\n0\n")

	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func readListFromFile(filename string) (*todoList, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	l := todoList{}

	scanner.Scan()
	l.name = scanner.Text()

	scanner.Scan()
	numItems, _ := strconv.Atoi(scanner.Text())
	l.numItems = numItems

	for i := 0; i < numItems; i++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		id, _ := strconv.Atoi(parts[0])
		l.items = append(l.items, listItem{id: id, desc: parts[1]})
	}

	return &l, nil
}

func (l *todoList) showList() {
	fmt.Println(l.name)
	for x := 0; x < len(l.name); x++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	for i := 0; i < l.numItems; i++ {
		fmt.Printf("%d. %s\n", l.items[i].id, l.items[i].desc)
	}
}

func (l *todoList) addItem(n string) {
	l.numItems++
	i := listItem{id: l.numItems, desc: n}
	l.items = append(l.items, i)
}

func (l *todoList) deleteItem(n int) {
	if n < 0 || n >= l.numItems {
		fmt.Println("Item index out of range")
		return
	}

	l.items = append(l.items[:n], l.items[n+1:]...)

	l.numItems--

	for i := n; i < l.numItems; i++ {
		l.items[i].id--
	}
}

func (l *todoList) saveList(filename string) error {

	var sb strings.Builder

	sb.WriteString(l.name + "\n")
	sb.WriteString(strconv.Itoa(l.numItems) + "\n")

	for _, item := range l.items {
		sb.WriteString(strconv.Itoa(item.id) + " " + item.desc + "\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)

}

func beginApp() {
	fmt.Println("Welcome to your Todo List Viewer!")
	reader := bufio.NewReader(os.Stdin)

	if !fileExists("todolist.txt") {
	outerLoop:
		for {
			answer, _ := getInput("You don't have a Todo list, would you like to create one (Y / N)?: ", reader)

			switch answer {
			case "Y":
				answer, _ = getInput("Enter Todo List name: ", reader)
				createNewList("todolist.txt", answer)
				break outerLoop
			case "N":
				fmt.Println("Okay, take care.")
				return
			default:
				fmt.Println("Invalid input, please enter a Y or N")
			}
		}
	}
	list, err := readListFromFile("todolist.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for {
		list.showList()
		opt, _ := getInput("Would you like to add an item, delete an item, or save and exit (A / D / S)?: ", reader)
		switch opt {
		case "A":
			response, _ := getInput("Please enter item to add: ", reader)
			list.addItem(response)
		case "D":
			response, _ := getInput("Please enter the number of the item you want to delete: ", reader)
			num, _ := strconv.Atoi(response)
			list.deleteItem(num - 1)
		case "S":
			err := list.saveList("todolist.txt")
			if err != nil {
				fmt.Println("Error occurred when attempting to save.")
			}
			fmt.Println("Saving and exiting.")
			return
		default:
			fmt.Println("Invalid option, please choose A, D, or S.")
		}
	}
}

func main() {
	beginApp()
}
