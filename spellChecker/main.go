package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// TrieNode represents a node in our Trie structure
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Trie is our main Trie structure for word storage and lookup
type Trie struct {
	root  *TrieNode
	depth int // Maximum depth for suggestion generation
}

// NewTrie initializes a new Trie with root node and depth for suggestions
func NewTrie() *Trie {
	return &Trie{
		root:  &TrieNode{children: make(map[rune]*TrieNode)},
		depth: 2, // adjust based on how many edits you want to consider
	}
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// SearchDirect checks if a word exists in the trie
func (t *Trie) SearchDirect(word string) bool {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

// generateVariations creates possible variations of a word by changing one character
func generateVariations(word string, depth int) []string {
	if depth <= 0 {
		return []string{word}
	}
	var variations []string
	for i := 0; i < len(word); i++ {
		for j := 'a'; j <= 'z'; j++ {
			if rune(word[i]) != j {
				newWord := word[:i] + string(j) + word[i+1:]
				variations = append(variations, newWord)
			}
		}
	}
	// Here's where you'd add more complex variations like insertions, deletions, etc. if needed
	return variations
}

// downloadFile fetches the content from the specified URL
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	url := "https://raw.githubusercontent.com/makifdb/spellcheck/main/words.txt"
	content, err := downloadFile(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	trie := NewTrie()
	for _, line := range bytes.Split(content, []byte("\n")) {
		word := strings.TrimSpace(string(line))
		if word != "" {
			trie.Insert(word)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a word to check spelling:")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	word := strings.TrimSpace(input)

	if trie.SearchDirect(word) {
		fmt.Printf("%s is spelled correctly.\n", word)
	} else {
		fmt.Printf("%s is not in the dictionary. Suggestions:\n", word)
		suggestions := generateVariations(word, trie.depth)
		suggestionsFound := false
		for _, suggestion := range suggestions {
			if trie.SearchDirect(suggestion) {
				fmt.Println(suggestion)
				suggestionsFound = true
			}
		}
		if !suggestionsFound {
			fmt.Println("No suggestions found.")
		}
	}
}
