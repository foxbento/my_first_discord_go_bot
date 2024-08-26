package main

import (
    "testing"
)

func TestModifyTwitterLinks(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Twitter link",
            input:    "Check out https://twitter.com/user/status/123456",
            expected: "Check out https://fxtwitter.com/user/status/123456",
        },
        {
            name:     "X link",
            input:    "Look at https://x.com/user/status/789012",
            expected: "Look at https://fixupx.com/user/status/789012",
        },
        {
            name:     "Multiple links",
            input:    "Twitter: https://twitter.com/user1/status/123 and X: https://x.com/user2/status/456",
            expected: "Twitter: https://fxtwitter.com/user1/status/123 and X: https://fixupx.com/user2/status/456",
        },
        {
            name:     "No links",
            input:    "Just a regular message",
            expected: "Just a regular message",
        },
        {
            name:     "Link in angle brackets",
            input:    "Don't modify this: <https://twitter.com/user/status/123456>",
            expected: "Don't modify this: <https://twitter.com/user/status/123456>",
        },
        {
            name:     "Mixed links",
            input:    "Modify this: https://twitter.com/user1/status/123 but not this: <https://x.com/user2/status/456>",
            expected: "Modify this: https://fxtwitter.com/user1/status/123 but not this: <https://x.com/user2/status/456>",
        },
        {
            name:     "Twitter link with www subdomain",
            input:    "Check out https://www.twitter.com/user/status/123456",
            expected: "Check out https://fxtwitter.com/user/status/123456",
        },
        {
            name:     "X link with www subdomain",
            input:    "Look at https://www.x.com/user/status/789012",
            expected: "Look at https://fixupx.com/user/status/789012",
        },
        {
            name:     "Links at start and end of string",
            input:    "https://twitter.com/user1/status/123 is interesting and so is https://x.com/user2/status/456",
            expected: "https://fxtwitter.com/user1/status/123 is interesting and so is https://fixupx.com/user2/status/456",
        },
        {
            name:     "HTTP links",
            input:    "Old link http://twitter.com/user/status/123456 and http://x.com/user/status/789012",
            expected: "Old link https://fxtwitter.com/user/status/123456 and https://fixupx.com/user/status/789012",
        },
        {
            name:     "Links with angle brackets in text",
            input:    "This <link> https://twitter.com/user/status/123456 and this <one> https://x.com/user/status/789012",
            expected: "This <link> https://fxtwitter.com/user/status/123456 and this <one> https://fixupx.com/user/status/789012",
        },
        {
            name:     "Mixed www and non-www links",
            input:    "Check https://www.twitter.com/user1/status/123 and https://x.com/user2/status/456",
            expected: "Check https://fxtwitter.com/user1/status/123 and https://fixupx.com/user2/status/456",
        },
        {
            name:     "X link with query parameters",
            input:    "https://x.com/Nefarious_Foxx/status/1827343634091409773?t=vz1CxWwkTUyboeZhODW_yw&s=19",
            expected: "https://fixupx.com/Nefarious_Foxx/status/1827343634091409773",
        },
        {
            name:     "Twitter link with query parameters",
            input:    "https://twitter.com/SinSquaredArt/status/1825669070588354674?t=NzSvTTYTI773iZgSnwTHpQ&s=19",
            expected: "https://fxtwitter.com/SinSquaredArt/status/1825669070588354674",
        },
        {
            name:     "X link with www and query parameters",
            input:    "https://www.x.com/CandySharkie/status/1826132464814682482",
            expected: "https://fixupx.com/CandySharkie/status/1826132464814682482",
        },
        {
            name:     "Twitter link with www and query parameters",
            input:    "https://www.twitter.com/CandySharkie/status/1826132464814682482",
            expected: "https://fxtwitter.com/CandySharkie/status/1826132464814682482",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := modifyTwitterLinks(tc.input)
            if result != tc.expected {
                t.Errorf("modifyTwitterLinks(%q) = %q; want %q", tc.input, result, tc.expected)
            }
        })
    }
}