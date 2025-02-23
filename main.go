package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type Credentials struct {
	Username string
	Password string
}

type Config struct {
	Username       string
	Password       string
	RefreshCycles  int
	LikeDelay      time.Duration
	ScrollDelay    time.Duration
	SessionTimeout time.Duration
	Headless       bool
}

type Stats struct {
	totalLikes  int
	startTime   time.Time
	lastRefresh time.Time
	cycleCount  int
}

func getSessionDir(username string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get home directory: ", err)
	}
	return filepath.Join(homeDir, ".xautoliker", "sessions", username)
}

func ensureSessionDir(username string) string {
	sessionDir := getSessionDir(username)
	err := os.MkdirAll(sessionDir, 0755)
	if err != nil {
		log.Fatal("Could not create session directory: ", err)
	}
	return sessionDir
}

func isLoggedIn(ctx context.Context) bool {
	var isLoggedIn bool
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://twitter.com/home"),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`!document.querySelector('a[href="/login"]')`, &isLoggedIn),
	)
	return err == nil && isLoggedIn
}

func main() {
	config := Config{}

	// Required flags
	flag.StringVar(&config.Username, "username", "", "X (Twitter) username or email")
	flag.StringVar(&config.Password, "password", "", "X (Twitter) password")

	// Optional settings with defaults
	flag.IntVar(&config.RefreshCycles, "refresh", 5, "Number of cycles before refreshing page (0 to disable)")
	flag.DurationVar(&config.LikeDelay, "like-delay", time.Second, "Delay between likes (e.g., 1s, 500ms)")
	flag.DurationVar(&config.ScrollDelay, "scroll-delay", 2*time.Second, "Delay after scrolling (e.g., 2s, 1500ms)")
	flag.DurationVar(&config.SessionTimeout, "timeout", 5*time.Minute, "Total session duration (e.g., 5m, 1h)")
	flag.BoolVar(&config.Headless, "headless", false, "Run in headless mode (browser window hidden)")

	flag.Parse()

	if config.Username == "" || config.Password == "" {
		flag.Usage()
		log.Fatal("Please provide both username and password")
	}

	fmt.Printf("Starting auto-liker with settings:\n")
	fmt.Printf("- Headless mode: %v\n", config.Headless)
	fmt.Printf("- Refresh every %d cycles\n", config.RefreshCycles)
	fmt.Printf("- Like delay: %v\n", config.LikeDelay)
	fmt.Printf("- Scroll delay: %v\n", config.ScrollDelay)
	fmt.Printf("- Session timeout: %v\n", config.SessionTimeout)

	sessionDir := ensureSessionDir(config.Username)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.Flag("disable-gpu", config.Headless),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.UserDataDir(sessionDir),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, config.SessionTimeout)
	defer cancel()

	if !isLoggedIn(ctx) {
		fmt.Println("Logging in...")
		err := loginToX(ctx, &config)
		if err != nil {
			log.Fatal("Failed to login: ", err)
		}
		fmt.Println("Login successful")
	} else {
		fmt.Println("Using existing session")
	}

	err := autoLikePosts(ctx, &config)
	if err != nil {
		log.Fatal("Error while liking posts: ", err)
	}
}

func loginToX(ctx context.Context, config *Config) error {
	return chromedp.Run(ctx,
		chromedp.Navigate("https://twitter.com/login"),
		chromedp.WaitVisible(`//input[@autocomplete="username"]`),
		chromedp.SendKeys(`//input[@autocomplete="username"]`, config.Username),
		chromedp.Click(`//span[contains(text(),"Next")]`),
		chromedp.WaitVisible(`//input[@autocomplete="current-password"]`),
		chromedp.SendKeys(`//input[@autocomplete="current-password"]`, config.Password),
		chromedp.Click(`//span[contains(text(),"Log in")]`),
		chromedp.WaitVisible(`//div[@data-testid="primaryColumn"]`),
	)
}

func autoLikePosts(ctx context.Context, config *Config) error {
	stats := &Stats{
		startTime:   time.Now(),
		lastRefresh: time.Now(),
	}

	for {
		// Handle page refresh if enabled
		if config.RefreshCycles > 0 && stats.cycleCount%config.RefreshCycles == 0 && stats.cycleCount > 0 {
			err := chromedp.Run(ctx,
				chromedp.Navigate("https://twitter.com/home"),
				chromedp.WaitVisible(`//div[@data-testid="primaryColumn"]`),
			)
			if err != nil {
				return err
			}
			stats.lastRefresh = time.Now()
			fmt.Printf("\nPage refreshed. Stats: %d likes in %d minutes\n",
				stats.totalLikes,
				int(time.Since(stats.startTime).Minutes()))
		}

		var likeButtons []string
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`
Array.from(document.querySelectorAll('article[data-testid="tweet"]')).map(article => {
  const button = article.querySelector('button[data-testid="like"][role="button"]:not([aria-label*="Unlike"])');
  if (button) {
    return button.getAttribute('aria-label') || 'no-label-' + Math.random();
  }
  return null;
}).filter(Boolean)
`, &likeButtons),
		)
		if err != nil {
			return err
		}

		if len(likeButtons) == 0 {
			err = chromedp.Run(ctx,
				chromedp.Evaluate(`window.scrollBy(0, 500)`, nil),
				chromedp.Sleep(config.ScrollDelay),
			)
			if err != nil {
				return err
			}
			stats.cycleCount++
			continue
		}

		for _, buttonLabel := range likeButtons {
			err = chromedp.Run(ctx,
				chromedp.Click(fmt.Sprintf(`//button[@data-testid="like"][@aria-label=%q]`, buttonLabel)),
				chromedp.Sleep(config.LikeDelay),
			)
			if err != nil {
				continue
			}
			stats.totalLikes++

			// Print statistics every 10 likes
			if stats.totalLikes%10 == 0 {
				duration := time.Since(stats.startTime)
				rate := float64(stats.totalLikes) / duration.Minutes()
				fmt.Printf("\nStats: %d likes in %d minutes (%.1f likes/minute)\n",
					stats.totalLikes,
					int(duration.Minutes()),
					rate)
			}
		}

		stats.cycleCount++
	}
}
