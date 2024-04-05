package main

import (
	"fmt"
	"wordfilter/wordfilter"
)

func main() {
	fmt.Println("Hello, world!")
	_, err := wordfilter.Init("./word.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	test := []string{
		"sb",
		"你好，傻！！sb2~~TMD， 他妈的~~~~",
		"sb2~~TMD， 他妈的",
		"傻！！sb2~~TMD， 他妈的~~~~",
		"干净!",
	}

	for _, v := range test {
		s := wordfilter.FilterSensitiveWord(v)
		fmt.Printf("old:%s\nnew:%s, hasSensitive:%t\n\n", v, s, wordfilter.HasSensitiveWord(v))
	}

}
