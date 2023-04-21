package main

import (
	"bingchat"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var spinChars = `|/-\`

func main() {
	cookiesPath := flag.String("c", "cookies.json", "cookies path")
	flag.Parse()
	if os.Getenv("https_proxy") != "" || os.Getenv("http_proxy") != "" {
		fmt.Println("local http proxy is set")
	}
	data, err := os.ReadFile(*cookiesPath)
	if err != nil {
		log.Fatalln(err)
	}
	help := `Bing Ai Chat Copilot
/reset [styles]
      1  more create
      2  more balance
      3  more precise
/quit
  quit
/help
  print help
`

	fmt.Printf(help+"Current Style: %s \n", bingchat.ConversationBalanceStyle.String())
	bingChat, err := bingchat.NewBingChat(data, bingchat.ConversationBalanceStyle)
	if err != nil {
		log.Fatalln(err)
	}
	var suggest []string
	for {
		var input string
		var style bingchat.ConversationStyle
		fmt.Printf("Ask> ")
		fmt.Scanln(&input, &style)
		if input == "" {
			continue
		}
		if strings.HasPrefix(input, "/reset") {
			var styles []bingchat.ConversationStyle
			if style >= 1 && style <= 3 {
				styles = append(styles, style)
			}
			bingChat.Reset(styles...)
			continue
		} else if strings.HasPrefix(input, "/quit") {
			fmt.Println("bye")
			os.Exit(0)
		} else if strings.HasPrefix(input, "/help") {
			fmt.Printf(help+"Current Style: %s \n", bingChat.Style())
			continue
		}
		index, err := strconv.Atoi(input)
		if err == nil {
			if index-1 >= 0 && index-1 < len(suggest) {
				input = suggest[index-1]
			}
		}
		msgFinished := make(chan struct{})
		go func() {
			i := 0
			for {
				select {
				case <-msgFinished:
					return
				default:
				}
				fmt.Printf("\r%s", string(spinChars[i%len(spinChars)]))
				time.Sleep(time.Second / 10)
				i++
			}
		}()
		resp, err := bingChat.SendMessage(input)
		close(msgFinished)
		fmt.Printf("\r")
		if err != nil {
			fmt.Println(err)
			continue
		}

		for {
			msg, ok := <-resp.Notify
			if !ok {
				fmt.Println()
				break
			}
			fmt.Printf("%s", msg)
		}
		if len(resp.Suggest) > 0 {
			fmt.Print("Prompt suggest\n")
			for i, suggest := range resp.Suggest {
				fmt.Printf("%d: %s\n", i+1, suggest)
			}
			suggest = resp.Suggest
		}
	}
}
