package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func main() {
	dictionaryFile := flag.String("dictionary", "dictionary.json", "The dictionary file to use")
	word := flag.String("string", "anagram", "The string to find anagrams for")
	num := flag.Int("maxnum", 1, "The maximum number of words to use")
	flag.Parse()

	fmt.Println("Finding anagrams for", "'"+*word+"'", "with", *num, "words using", *dictionaryFile)

	// Read the file
	dictionary, err := ioutil.ReadFile(*dictionaryFile)

	if err != nil {
		panic(err)
	}

	// Unmarshal the data
	var data map[string]int
	err = json.Unmarshal(dictionary, &data)

	if err != nil {
		panic(err)
	}

	// Print the result
	words := []string{}

	for key, value := range data {
		// if key == "a" || key == "u" || key == "i" {
		// 	fmt.Println(key, value)
		// }

		if value == 1 {
			words = append(words, key)
		}
	}

	// words := maps.Keys(data)

	sort.Slice(words, func(i, j int) bool {
		if len(words[i]) == len(words[j]) {
			return words[i] < words[j]
		}

		return len(words[i]) < len(words[j])
	})

	anagrams := GenerateAnagrams(words, map[int]map[string]struct{}{}, *word, -1, *num)

	for _, anagram := range anagrams {
		fmt.Println(anagram)
	}
	// fmt.Println(words, Anagramer(words, "amanap"))
}

func CompareStrings(str1 string, str2 string) bool {
	if len(str1) == len(str2) {
		return str1 < str2
	}

	return len(str1) < len(str2)
}

// func Anagrammer(words []string, str string) []string {
// 	str = strings.Replace(str, " ", "", -1)

// 	sort.Slice(words, func(i, j int) bool {
// 		if len(words[i]) == len(words[j]) {
// 			return words[i] < words[j]
// 		}

// 		return len(words[i]) < len(words[j])
// 	})

// 	return GetAnagrams(words, map[int][]string{}, str, "", len(str))
// }

// func GetAnagrams(words []string, work map[int][]string, str string, current string, length int) []string {
// 	anagrams := []string{}

// 	for i := 0; i < len(str); i++ {
// 		startString := ""

// 		for j := 0; j < i; j++ {
// 			startString += "a"
// 		}

// 		start, _ := slices.BinarySearchFunc(words, startString, func(str1 string, str2 string) int {
// 			if len(str1) == len(str2) {
// 				return strings.Compare(str1, str2)
// 			}

// 			return len(str1) - len(str2)
// 		})

// 		for j := start; len(words[j]) == i; j++ {
// 			anagrams = append(anagrams, GetAnagrams(words)...)
// 		}
// 	}

// 	return work[len(str)]
// }

func GenerateAnagrams(words []string, work map[int]map[string]struct{}, str string, length int, num int) []string {
	str = strings.ToLower(strings.Replace(str, " ", "", -1))
	// fmt.Println("str:", str)

	if num == 0 {
		return []string{}
	}

	if length == -1 {
		length = len(str)
	}

	_, ok := work[length]

	if ok {
		// fmt.Println("hit", length)

		return maps.Keys(work[length])
	}

	anagrams := []string{}

	for i := length - 1; i >= 0; i-- {
		// fmt.Println(i)
		left := GenerateAnagrams(words, work, str, i, num-1)

		startString := ""

		for j := 0; j < length-i; j++ {
			startString += "A"
		}

		start, _ := slices.BinarySearchFunc(words, startString, func(str1 string, str2 string) int {
			if len(str1) == len(str2) {
				return strings.Compare(str1, str2)
			}

			return len(str1) - len(str2)
		})

		endString := ""

		for j := 0; j < length-i; j++ {
			endString += "z"
		}

		end, _ := slices.BinarySearchFunc(words, endString, func(str1 string, str2 string) int {
			if len(str1) == len(str2) {
				return strings.Compare(str1, str2)
			}

			return len(str1) - len(str2)
		})

		right := words[start:end]

		// fmt.Println(start, startString, right, endString, end)
		// fmt.Println(i, len(left), len(right), words[0:26])
		// fmt.Println(startString, endString)
		for j := 0; j < len(left); j++ {
			for k := 0; k < len(right); k++ {
				// fmt.Println("^", i, left[j]+right[k])
				if IsInString(str, left[j]+right[k]) {
					// fmt.Println(right)
					anagrams = append(anagrams, left[j]+" "+right[k])
				}
			}
		}

		if i == 0 {
			for j := 0; j < len(right); j++ {
				// fmt.Println("^", i, right[j])
				if IsInString(str, right[j]) {
					// fmt.Println(right)
					anagrams = append(anagrams, right[j])
				}
			}
		}
	}

	work[length] = make(map[string]struct{})

	for _, word := range anagrams {
		work[length][word] = struct{}{}
	}

	return maps.Keys(work[length])
}

func IsInString(str1 string, str2 string) bool {
	str1 = strings.Replace(str1, " ", "", -1)
	str2 = strings.Replace(str2, " ", "", -1)
	chars := map[rune]int{}

	if len(str2) > len(str1) {
		return false
	}

	for _, char := range str1 {
		chars[char]++
	}

	for _, char := range str2 {
		chars[char]--

		if chars[char] < 0 {
			return false
		}
	}

	return true
}

func IsAnagram(str1 string, str2 string) bool {
	str1 = strings.Replace(str1, " ", "", -1)
	str2 = strings.Replace(str2, " ", "", -1)

	if len(str1) != len(str2) {
		return false
	}

	chars := map[rune]int{}

	for _, char := range str1 {
		chars[char]++
	}

	for _, char := range str2 {
		chars[char]--
	}

	for _, count := range chars {
		if count != 0 {
			return false
		}
	}

	return true
}
