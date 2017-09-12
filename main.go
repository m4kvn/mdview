package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"net/http"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		printUsage()
		return
	}

	if args[0] == "-u" {
		readFromURL(args[1], args[2:])
	} else {
		readFromFile(args[0], args[1:])
	}
}

func readFromFile(fileName string, args []string) {
	_, err := os.Stat(fileName)

	if err != nil {
		fmt.Println(err)
		printUsage()
		return
	}

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)

	if len(args) == 0 {
		displayAll(scanner)
	} else if args[0] == "-l" {
		displayHeaders(scanner)
	} else {
		displayTexts(scanner, args)
	}
}

func readFromURL(url string, args []string) {
	if !checkUrl(url) {
		fmt.Println("The url is not Markdown file.")
		printUsage()
		return
	}

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		defer resp.Body.Close()
	}

	scanner := bufio.NewScanner(resp.Body)

	if len(args) == 0 {
		displayAll(scanner)
	} else if args[0] == "-l" {
		displayHeaders(scanner)
	} else {
		displayTexts(scanner, args)
	}
}

func checkUrl(url string) bool {
	tmp := strings.Split(url, "/")
	fileName := tmp[len(tmp)-1]
	return fileName[len(fileName)-3:] == ".md"
}

func displayAll(scanner *bufio.Scanner) {
	for scanner.Scan() {
		scanner.Text()
		fmt.Println(scanner.Text())
	}
}

func displayHeaders(scanner *bufio.Scanner) {
	for scanner.Scan() {
		txt := scanner.Text()

		if strings.Index(txt, "#") != 0 {
			continue
		}

		fmt.Println(txt)
	}
}

func displayTexts(scanner *bufio.Scanner, args []string) {
	headerNum := -1

	for _, str := range args {
		for scanner.Scan() {
			txt := scanner.Text()

			if strings.Index(txt, "#") != 0 {
				continue
			}

			if !strings.Contains(txt, str) {
				continue
			}

			headerNum = getHeaderNum(txt)

			break
		}
	}

	for scanner.Scan() {
		txt := scanner.Text()
		num := getHeaderNum(txt)
		if headerNum >= num && num > 0 {
			break
		}
		fmt.Println(txt)
	}
}

func getHeaderNum(str string) int {
	for i, c := range str {
		if c != '#' {
			return i
		}
	}
	return -1
}

func printUsage() {
	fmt.Println("usage: mdview (<file_path> | -u <url>) [-l] [<header_name>]")
}
