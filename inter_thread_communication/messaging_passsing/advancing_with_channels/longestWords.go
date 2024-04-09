/* Adicionando mais uma goroutine na pipeline, dessa vez para retornar as 10 palavras mais longas */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const downloaders2 = 20

func main() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls3(quit)

	pages := make([]<-chan string, downloaders2)

	for i := 0; i < downloaders2; i++ {
		pages[i] = downloadPages3(quit, urls)
	}

	paginasUnificadas := FanIn2(quit, pages...)
	words := extractWords3(quit, paginasUnificadas)
	results := longestWords(quit, words)

	for result := range results {
		fmt.Println(result)
	}
}

func generateUrls3(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i < 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	return urls
}

func downloadPages3(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	return pages
}

func extractWords3(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, pg := true, ""
		for moreData {
			select {
			case pg, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}

func FanIn2[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))
	output := make(chan K)
	for _, c := range allChannels {
		go func(channel <-chan K) {
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

/*
curiosidade: se eu tivesse 5 goroutines para essa função,
cada goroutine iria retornar 1 canal com 10 palavras
*/
func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}
		sort.Slice(uniqueWords, func(a, b int) bool {
			return len(uniqueWords[a]) > len(uniqueWords[b])
		})
		longWords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longWords
}
