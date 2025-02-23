package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Credentials struct {
	Username string
	Password string
}

func main() {
	// Command line flags for credentials
	username := flag.String("username", "", "X (Twitter) username or email")
	password := flag.String("password", "", "X (Twitter) password")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Please provide both username and password")
	}

	creds := &Credentials{
		Username: *username,
		Password: *password,
	}

	// Create context with visible browser
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Login to X
	err := loginToX(ctx, creds)
	if err != nil {
		log.Fatal("Failed to login: ", err)
	}

	// Start liking posts
	err = autoLikePosts(ctx)
	if err != nil {
		log.Fatal("Error while liking posts: ", err)
	}
}

func loginToX(ctx context.Context, creds *Credentials) error {
	// Navigate to login page
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://twitter.com/login"),
		chromedp.WaitVisible(`//input[@autocomplete="username"]`),
		chromedp.SendKeys(`//input[@autocomplete="username"]`, creds.Username),
		chromedp.Click(`//span[contains(text(),"Next")]`),
		chromedp.WaitVisible(`//input[@autocomplete="current-password"]`),
		chromedp.SendKeys(`//input[@autocomplete="current-password"]`, creds.Password),
		chromedp.Click(`//span[contains(text(),"Log in")]`),
		chromedp.WaitVisible(`//div[@data-testid="primaryColumn"]`),
	)
	return err
}

func autoLikePosts(ctx context.Context) error {
	// Navigate to Following timeline
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://twitter.com/home"),
		chromedp.WaitVisible(`//div[@data-testid="primaryColumn"]`),
	)
	if err != nil {
		return err
	}

	// Main liking loop
	for {
		var likeButtons []string
		err = chromedp.Run(ctx,
			// Find all unlike buttons (which means posts that haven't been liked yet)
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll('div[data-testid="like"]')).map(el => el.id)
			`, &likeButtons),
		)
		if err != nil {
			return err
		}

		if len(likeButtons) == 0 {
			// Scroll down to load more posts
			err = chromedp.Run(ctx,
				chromedp.Evaluate(`window.scrollBy(0, 500)`, nil),
				chromedp.Sleep(2*time.Second), // Wait for new posts to load
			)
			if err != nil {
				return err
			}
			continue
		}

		// Like each post
		for _, buttonID := range likeButtons {
			err = chromedp.Run(ctx,
				chromedp.Click(`//div[@id="`+buttonID+`"]`),
				chromedp.Sleep(1*time.Second), // Delay between likes to avoid rate limiting
			)
			if err != nil {
				log.Printf("Failed to like a post: %v", err)
				continue
			}
		}
	}
}
