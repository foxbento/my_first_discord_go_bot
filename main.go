// Package main provides a Discord bot that responds to messages and modifies Twitter/X links.
package main

import (
	"fmt"
	"log"
	"net/url"
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
        return
    }

    // Check for Twitter/X links
    if containsTwitterLink(m.Content) {
        // Log raw embed array
        logRawEmbeds(m.Embeds)

        // Log detailed information for messages with Twitter links
        logTwitterMessage(m)

        // Check if the message has any valid Twitter embeds
        hasValidEmbed := hasValidTwitterEmbed(m.Embeds)
        log.Printf("Has valid Twitter embed: %v\n", hasValidEmbed)

        if !hasValidEmbed {
            log.Println("No valid Twitter embed found, modifying Twitter link")
            modifiedContent := modifyTwitterLinks(m.Content)
            
            if modifiedContent != m.Content {
                _, err := s.ChannelMessageSend(m.ChannelID, modifiedContent)
                if err != nil {
                    log.Println("Error sending modified message:", err)
                }
            }
        } else {
            log.Println("Valid Twitter embed found, not modifying message")
        }
    }
}

func logRawEmbeds(embeds []*discordgo.MessageEmbed) {
    log.Printf("Raw Embed Array: %+v\n", embeds)
    for i, embed := range embeds {
        log.Printf("Embed %d raw data: %+v\n", i, embed)
    }
}

func containsTwitterLink(content string) bool {
    pattern := `https?:\/\/(www\.)?(twitter\.com|x\.com)\/[a-zA-Z0-9_]+\/status\/[0-9]+`
    match, _ := regexp.MatchString(pattern, content)
    return match
}

func hasValidTwitterEmbed(embeds []*discordgo.MessageEmbed) bool {
    for _, embed := range embeds {
        if isTwitterEmbed(embed) {
            return true
        }
    }
    return false
}

func isTwitterEmbed(embed *discordgo.MessageEmbed) bool {
    // List of Twitter CDN domains
    twitterCDNs := []string{
        "pbs.twimg.com",
        "video.twimg.com",
        "abs.twimg.com",
        "ton.twimg.com",
    }

    // Check embed URL
    if embed.URL != "" {
        u, err := url.Parse(embed.URL)
        if err == nil {
            for _, cdn := range twitterCDNs {
                if strings.HasSuffix(u.Hostname(), cdn) {
                    return true
                }
            }
        }
    }

    // Check image URL
    if embed.Image != nil && embed.Image.URL != "" {
        u, err := url.Parse(embed.Image.URL)
        if err == nil {
            for _, cdn := range twitterCDNs {
                if strings.HasSuffix(u.Hostname(), cdn) {
                    return true
                }
            }
        }
    }

    // Check thumbnail URL
    if embed.Thumbnail != nil && embed.Thumbnail.URL != "" {
        u, err := url.Parse(embed.Thumbnail.URL)
        if err == nil {
            for _, cdn := range twitterCDNs {
                if strings.HasSuffix(u.Hostname(), cdn) {
                    return true
                }
            }
        }
    }

    return false
}

func logTwitterMessage(m *discordgo.MessageCreate) {
    log.Printf("Twitter link detected in message ID: %s\n", m.ID)
    log.Printf("Author: %s (ID: %s)\n", m.Author.Username, m.Author.ID)
    log.Printf("Channel ID: %s\n", m.ChannelID)
    log.Printf("Content: %s\n", m.Content)
    log.Printf("Timestamp: %s\n", m.Timestamp)
    log.Printf("Number of embeds: %d\n", len(m.Embeds))
    
    for i, embed := range m.Embeds {
        logEmbedDetails(i, embed)
    }
}

func logEmbedDetails(index int, embed *discordgo.MessageEmbed) {
    log.Printf("Embed %d details:\n", index)
    log.Printf("  Type: %s\n", embed.Type)
    log.Printf("  Title: %s\n", embed.Title)
    log.Printf("  Description: %s\n", embed.Description)
    log.Printf("  URL: %s\n", embed.URL)
    if embed.Image != nil {
        log.Printf("  Image URL: %s\n", embed.Image.URL)
    }
    if embed.Thumbnail != nil {
        log.Printf("  Thumbnail URL: %s\n", embed.Thumbnail.URL)
    }
    log.Printf("  Is Twitter Embed: %v\n", isTwitterEmbed(embed))
}

// modifyTwitterLinks takes a string and replaces Twitter/X links with modified versions.
// It changes "twitter.com" to "fxtwitter.com" and "x.com" to "fixupx.com".
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