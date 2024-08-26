// Package main provides a Discord bot that responds to messages and modifies Twitter/X links.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// init loads the environment variables from a .env file.
// It should be called automatically before the main function.
func init() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file:", err)
		} else {
			log.Println("Loaded environment variables from .env file")
		}
	} else {
		log.Println("No .env file found, using system environment variables")
	}
}

// main is the entry point of the application.
// It sets up the Discord session, registers event handlers,
// and keeps the bot running until interrupted.
func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("No token provided. Set DISCORD_BOT_TOKEN in your .env file.")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	sess.AddHandler(messageCreate)

	sess.Identify.Intents = discordgo.IntentsGuildMessages

	err = sess.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer sess.Close()

	fmt.Println("The bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// messageCreate is the callback function for the MessageCreate event.
// It handles incoming messages, responds to "hello", and modifies Twitter/X links.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Respond to "hello" messages
	if m.Content == "hello" {
		_, err := s.ChannelMessageSend(m.ChannelID, "world!")
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}

	// Check for Twitter/X links and modify them
	// now only modifies if embed array is empty (thanks @tingerrrr)
	modifiedContent := modifyTwitterLinks(m.Content)
	if modifiedContent != m.Content {
		_, err := s.ChannelMessageSend(m.ChannelID, modifiedContent)
		if err != nil {
			log.Println("Error sending modified message:", err)
		}
	}
}

// modifyTwitterLinks takes a string and replaces Twitter/X links with modified versions.
// It changes "twitter.com" to "fxtwitter.com" and "x.com" to "fixupx.com".
//
// Example:
//
//	input := "Check out https://twitter.com/user/status/123456"
//	output := modifyTwitterLinks(input)
//	// output will be "Check out https://fxtwitter.com/user/status/123456"
func modifyTwitterLinks(content string) string {
    // Define patterns for Twitter and X links, including those in angle brackets
    pattern := `(<)?https?://(www\.)?(twitter\.com|x\.com)/[^/]+/status/\d+(\?[^\s<>]*)?([^<\s]*)>?`

    re := regexp.MustCompile(pattern)
    return re.ReplaceAllStringFunc(content, func(match string) string {
        if strings.HasPrefix(match, "<") && strings.HasSuffix(match, ">") {
            return match // Preserve links in angle brackets
        }
        return modifySingleLink(match)
    })
}

func modifySingleLink(link string) string {
    // Remove query parameters
    if idx := strings.Index(link, "?"); idx != -1 {
        link = link[:idx]
    }

    // Strip protocol and www subdomain
    link = strings.TrimPrefix(link, "http://")
    link = strings.TrimPrefix(link, "https://")
    link = strings.TrimPrefix(link, "www.")

    // Replace domain
    if strings.HasPrefix(link, "twitter.com") {
        link = "https://fxtwitter.com" + strings.TrimPrefix(link, "twitter.com")
    } else if strings.HasPrefix(link, "x.com") {
        link = "https://fixupx.com" + strings.TrimPrefix(link, "x.com")
    }

    return link
}