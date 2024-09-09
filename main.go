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
		}
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
        // Log detailed information about the message and its embeds
        logTwitterMessage(m)

        // Check if the message has any valid Twitter embeds or attachments
        hasValidPreview := hasValidTwitterPreview(m)

        if !hasValidPreview {
            modifiedContent := modifyTwitterLinks(m.Content)
            
            if modifiedContent != m.Content {
                _, err := s.ChannelMessageSend(m.ChannelID, modifiedContent)
                if err != nil {
                    log.Println("Error sending modified message:", err)
                }
            }
        }
    }
}

// logTwitterMessage logs detailed information about a message containing a Twitter link
func logTwitterMessage(m *discordgo.MessageCreate) {
    twitterLinks := extractTwitterLinks(m.Content)
    
    for _, link := range twitterLinks {
        log.Printf("Original Twitter Link: %s\n", link)
        log.Printf("Total Embeds in Message: %d\n", len(m.Embeds))
        
        for i, embed := range m.Embeds {
            log.Printf("Embed %d:\n", i+1)
            log.Printf("  Type: %s\n", embed.Type)
            log.Printf("  Title: %s\n", embed.Title)
            log.Printf("  Description: %s\n", embed.Description)
            
            if embed.Image != nil {
                log.Printf("  Image URL: %s\n", embed.Image.URL)
            }
            
            if embed.Thumbnail != nil {
                log.Printf("  Thumbnail URL: %s\n", embed.Thumbnail.URL)
            }
            
            log.Printf("  Fields: %d\n", len(embed.Fields))
        }
        
        log.Printf("Total Attachments in Message: %d\n", len(m.Attachments))
        
        for i, attachment := range m.Attachments {
            log.Printf("Attachment %d:\n", i+1)
            log.Printf("  Filename: %s\n", attachment.Filename)
            log.Printf("  URL: %s\n", attachment.URL)
            log.Printf("  Size: %d bytes\n", attachment.Size)
        }
    }
}

func containsTwitterLink(content string) bool {
    pattern := `https?:\/\/(www\.)?(twitter\.com|x\.com)\/[a-zA-Z0-9_]+\/status\/[0-9]+`
    match, _ := regexp.MatchString(pattern, content)
    return match
}

func extractTwitterLinks(content string) []string {
    pattern := `https?://(www\.)?(twitter\.com|x\.com)/[^/]+/status/\d+`
    re := regexp.MustCompile(pattern)
    return re.FindAllString(content, -1)
}

func hasValidTwitterPreview(m *discordgo.MessageCreate) bool {
    // Check embeds
    for _, embed := range m.Embeds {
        if isWorkingTwitterEmbed(embed) {
            return true
        }
    }

    // Check attachments
    for _, attachment := range m.Attachments {
        if isWorkingTwitterAttachment(attachment) {
            return true
        }
    }

    return false
}

func isWorkingTwitterEmbed(embed *discordgo.MessageEmbed) bool {
    // List of Twitter CDN domains
    twitterCDNs := []string{
        "pbs.twimg.com",
        "video.twimg.com",
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
            if strings.HasSuffix(u.Hostname(), "abs.twimg.com") {
                return false
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
            if strings.HasSuffix(u.Hostname(), "abs.twimg.com") {
                return false
            }
        }
    }

    // Check thumbnail URL
    if embed.Thumbnail != nil && embed.Thumbnail.URL != "" {
        u, err := url.Parse(embed.Thumbnail.URL)
        if err == nil {
            // Check for tweet_video_thumb or amplify_video_thumb in the thumbnail URL
            if strings.Contains(embed.Thumbnail.URL, "tweet_video_thumb") || 
                    strings.Contains(embed.Thumbnail.URL, "amplify_video_thumb") {
                    return false
                }
            for _, cdn := range twitterCDNs {
                if strings.HasSuffix(u.Hostname(), cdn) {
                    return true
                }
            }
            if strings.HasSuffix(u.Hostname(), "abs.twimg.com") {
                return false
            }
        }
    }

    return false
}

func isWorkingTwitterAttachment(attachment *discordgo.MessageAttachment) bool {
    // List of Twitter CDN domains
    twitterCDNs := []string{
        "pbs.twimg.com",
        "video.twimg.com",
        "ton.twimg.com",
    }

    u, err := url.Parse(attachment.URL)
    if err == nil {
        for _, cdn := range twitterCDNs {
            if strings.HasSuffix(u.Hostname(), cdn) {
                return true
            }
        }
        if strings.HasSuffix(u.Hostname(), "abs.twimg.com") {
            return false
        }
    }

    return false
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