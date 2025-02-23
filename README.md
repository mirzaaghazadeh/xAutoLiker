# X Auto Liker

An automated tool for liking posts on X (formerly Twitter) with configurable settings and session management.

## Features

- 👍 Automatically likes posts in your timeline
- ⚙️ Highly configurable settings (delays, refresh cycles, timeouts)
- 🔒 Secure session management
- 📊 Real-time statistics
- 🌐 Supports both visible and headless browser modes
- 💾 Persistent sessions for quick reuse

## Installation

```bash
go install github.com/mirzaaghazadeh/xAutoLiker@latest
```

This will install `xAutoLiker` to your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is added to your `PATH`.

## Usage

### Basic Command

```bash
xAutoLiker -username "your_username" -password "your_password"
```

### Available Settings

1. Required:
   - `-username`: Your X (Twitter) username/email
   - `-password`: Your X (Twitter) password

2. Optional with defaults:
   - `-headless`: Run in headless mode without browser window (default: false)
   - `-refresh`: Refresh page every N cycles (default: 5, 0 to disable)
   - `-like-delay`: Delay between likes (default: 1s)
   - `-scroll-delay`: Delay after scrolling (default: 2s)
   - `-timeout`: Total session duration (default: 5m)

### Example Commands

```bash
# Run in headless mode with fast liking
xAutoLiker -username "user" -password "pass" -headless -like-delay 500ms

# Visible browser with longer delays
xAutoLiker -username "user" -password "pass" -refresh 0 -like-delay 2s

# Long session in headless mode
xAutoLiker -username "user" -password "pass" -headless -timeout 2h -like-delay 1s
```

## Session Management

Sessions are stored in `~/.xautoliker/sessions/[username]`. This allows for:
- Faster subsequent logins
- Persistent authentication
- Multiple account support

## Statistics

The tool provides real-time statistics:
- Total likes count
- Session duration
- Like rate (likes per minute)
- Stats are shown every 10 likes and after page refreshes

## Requirements

- Go 1.21 or higher
- Chrome/Chromium browser installed

## Security Note

- Your credentials are never stored, only browser session data is saved
- Sessions are stored securely in your home directory
- Run in headless mode for background operation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details