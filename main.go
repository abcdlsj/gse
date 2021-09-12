package main

import (
	"bufio"
	"fmt"
	"go_search/engine"
	"os"
	"time"
)

var dumpPath = "/Users/songjian.li/Desktop/zhwiki-latest-abstract.xml"

func main() {
	fmt.Println("Starting load")
	start := time.Now()

	documentChan := make(chan int)
	dictChan := make(chan int)

	var docs []engine.Documents

	go func() {
		docs = engine.LoadDocument(dumpPath)
		fmt.Printf("Loaded %d documents in %v\n", len(docs), time.Since(start))
		documentChan <- 1
	}()

	go func() {
		engine.LoadDict()
		dictChan <- 1
	}()

	_, _ = <-documentChan, <-dictChan

	start = time.Now()
	idx := make(engine.Index)

	idx.Add(docs)
	fmt.Printf("Indexed %d documents in %v\n", len(docs), time.Since(start))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		query := scanner.Text()
		fmt.Printf("\n> Input Search Query: ")
		if query == ".exit" {
			break
		}
		start = time.Now()
		matchedIDs := idx.Search(query)
		fmt.Printf("Search found %d documents in %v\n", len(matchedIDs), time.Since(start))
		matchedIDs = engine.SortResult(matchedIDs, docs, query)
		for _, id := range matchedIDs {
			fmt.Printf("%d\t%s\n", id, docs[id].Abstract)
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}
