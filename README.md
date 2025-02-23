# X AutoLiker

This Go application automatically likes the latest posts from accounts you follow on X (formerly Twitter).

## Features

- Browser-based automation (visible browser mode)
- Automatic login to X
- Scrolls through your timeline
- Automatically likes posts that haven't been liked yet
- Built-in delays to avoid rate limiting

## Prerequisites

- Go 1.16 or higher
- Chrome/Chromium browser installed

## Installation

1. Clone this repository
2. Install dependencies:
```bash
go mod download
```

## Usage

Run the program with your X credentials:

```bash
go run main.go -username "your_username_or_email" -password "your_password"
```

## How it Works

1. The program launches a Chrome browser
2. Logs into X using provided credentials
3. Navigates to your home timeline
4. Automatically finds and likes unread posts
5. Scrolls down to load more posts as needed

## Safety Features

- Visible browser mode for monitoring
- Built-in delays between actions
- Error handling for failed interactions

## Note

This is a basic implementation. Use responsibly and be aware of X's terms of service regarding automation.

## Future Improvements

- Headless browser mode option
- Customizable delay between likes
- Like limit configuration
- Specific user targeting
- OAuth authentication instead of password