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
	str := flag.String("string", "anagram", "The string to find anagrams for")
	maxWords := flag.Int("maxwords", -1, "The maximum number of words to use")
	flag.Parse()

	// Read the dictionary file
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

	// Add enabled words to the word list
	words := []string{}

	for key, value := range data {
		if value == 1 {
			words = append(words, key)
		}
	}

	sort.Slice(words, func(i, j int) bool {
		if len(words[i]) == len(words[j]) {
			return words[i] < words[j]
		}

		return len(words[i]) < len(words[j])
	})

	// Find anagrams
	fmt.Println("Finding anagrams for", "'"+*str+"'", "with", *maxWords, "words using", *dictionaryFile)
	anagrams := Anagrammer(words, *str, *maxWords)

	for _, anagram := range anagrams {
		fmt.Println(anagram)
	}
}

func Anagrammer(words []string, str string, maxWords int) []string {
	return GenerateAnagrams(words, map[int]map[string]struct{}{}, str, -1, maxWords)
}

// String comparison function accounting for
// length first and lexigraphical order second
func CompareStrings(str1 string, str2 string) int {
	if len(str1) == len(str2) {
		return strings.Compare(str1, str2)
	}

	return len(str1) - len(str2)
}

// Check if str2 could be built from characters in str1
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
	}

	for _, count := range chars {
		if count < 0 {
			return false
		}
	}

	return true
}

// Generates a slice of strings containing all possible
// "subanagrams" of str of length length and number of words
// at most maxWords.
// work is used to store previous results from smaller
// subproblems.
func GenerateAnagrams(words []string, work map[int]map[string]struct{}, str string, length int, maxWords int) []string {
	str = strings.ToLower(strings.Replace(str, " ", "", -1))

	// Base case: cannot use any words
	if maxWords == 0 {
		return []string{}
	}

	// Auto-detect length if not specified
	if length == -1 {
		length = len(str)
	}

	// Check if we've already computed this
	// subproblem, and return it if we have
	_, ok := work[length]

	if ok {
		return maps.Keys(work[length])
	}

	// For different lengths, generate an anagram of that length
	// from the set of characters in str and add a word of appropriate
	// length from the dictionary to it.
	anagrams := []string{}

	for i := length - 1; i >= 0; i-- {
		left := GenerateAnagrams(words, work, str, i, maxWords-1)

		// Find the words with the same length as the remaining length
		startString := ""

		for j := 0; j < length-i; j++ {
			startString += " "
		}

		start, _ := slices.BinarySearchFunc(words, startString, CompareStrings)

		endString := ""

		for j := 0; j < length-i; j++ {
			endString += "~"
		}

		end, _ := slices.BinarySearchFunc(words, endString, CompareStrings)

		right := words[start:end]

		// Check if any of the potential anagrams work,
		// and add them if they do.
		for j := 0; j < len(left); j++ {
			for k := 0; k < len(right); k++ {
				if IsInString(str, left[j]+right[k]) {
					anagrams = append(anagrams, left[j]+" "+right[k])
				}
			}
		}

		if i == 0 {
			for j := 0; j < len(right); j++ {
				if IsInString(str, right[j]) {
					anagrams = append(anagrams, right[j])
				}
			}
		}
	}

	// Add our anagrams to the work map
	work[length] = make(map[string]struct{})

	for _, word := range anagrams {
		work[length][word] = struct{}{}
	}

	return maps.Keys(work[length])
}
