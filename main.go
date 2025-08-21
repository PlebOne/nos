package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/zalando/go-keyring"
)

const (
	appName       = "nos"
	keyringUser   = "nos-cli"
	keyringKey    = "nsec"
	relayListKey  = "relay-list"
)

var (
	// Charm.sh styling
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))
)

// Default relay list
var defaultRelays = []string{
	"wss://relay.damus.io",
	"wss://nos.lol",
	"wss://relay.nostr.band",
	"wss://relay.current.fyi",
	"wss://relay.snort.social",
	"wss://relay.primal.net",
}

func main() {
	// If no arguments, show interactive menu
	if len(os.Args) < 2 {
		showMainMenu()
		return
	}

	// Check for command-line commands (for backwards compatibility)
	switch os.Args[1] {
	case "reset", "-reset":
		handleReset()
	case "relay", "-relay":
		handleRelayCommand()
	case "verify", "-verify":
		handleVerify()
	default:
		// Assume it's a message to post
		quickPost(strings.Join(os.Args[1:], " "))
	}
}

func showMainMenu() {
	for {
		// Check if user has set up their key
		_, keyErr := getStoredKey()
		hasKey := keyErr == nil

		fmt.Println(titleStyle.Render("nos - Nostr CLI 🚀"))
		
		if hasKey {
			// Get and display public key
			nsec, _ := getStoredKey()
			var sk string
			_, s, _ := nip19.Decode(nsec)
			sk = s.(string)
			pub, _ := nostr.GetPublicKey(sk)
			npub, _ := nip19.EncodePublicKey(pub)
			fmt.Println(infoStyle.Render("Your npub: " + npub))
		} else {
			fmt.Println(errorStyle.Render("No account configured"))
		}
		fmt.Println()

		var choice string
		var options []huh.Option[string]

		if hasKey {
			options = []huh.Option[string]{
				huh.NewOption("Post a message", "post"),
				huh.NewOption("Verify your posts", "verify"),
				huh.NewOption("Manage relays", "relay"),
				huh.NewOption("Reset account", "reset"),
				huh.NewOption("Exit", "exit"),
			}
		} else {
			options = []huh.Option[string]{
				huh.NewOption("Setup account (add nsec)", "setup"),
				huh.NewOption("Exit", "exit"),
			}
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("What would you like to do?").
					Options(options...).
					Value(&choice),
			),
		)

		err := form.Run()
		if err != nil {
			return
		}

		switch choice {
		case "setup":
			interactiveSetup()
		case "post":
			interactivePost()
		case "verify":
			handleVerify()
			fmt.Print("\nPress Enter to continue...")
			fmt.Scanln()
		case "relay":
			showRelayMenu()
		case "reset":
			interactiveReset()
		case "exit":
			fmt.Println(infoStyle.Render("Goodbye! 👋"))
			return
		}

		// Clear screen for next iteration (optional, you can remove this if you prefer)
		fmt.Print("\033[H\033[2J")
	}
}

func quickPost(message string) {
	// Try to get stored key
	nsec, err := getStoredKey()
	if err != nil {
		// First time setup
		fmt.Println(titleStyle.Render("Welcome to nos! 🚀"))
		fmt.Println(infoStyle.Render("It looks like this is your first time. Let's set up your Nostr key."))
		fmt.Println()

		nsec, err = promptForKey()
		if err != nil {
			fmt.Println(errorStyle.Render("Error: " + err.Error()))
			os.Exit(1)
		}

		// Store the key
		err = storeKey(nsec)
		if err != nil {
			fmt.Println(errorStyle.Render("Error storing key: " + err.Error()))
			os.Exit(1)
		}

		fmt.Println(successStyle.Render("✓ Key stored securely!"))
		fmt.Println()
	}

	// Convert nsec to private key
	var sk string
	_, s, err := nip19.Decode(nsec)
	if err != nil {
		fmt.Println(errorStyle.Render("Error decoding key: " + err.Error()))
		os.Exit(1)
	}
	sk = s.(string)

	// Show public key for verification
	pub, _ := nostr.GetPublicKey(sk)
	npub, _ := nip19.EncodePublicKey(pub)
	fmt.Println(infoStyle.Render("Your npub: " + npub))
	
	// Post to Nostr
	fmt.Println(infoStyle.Render("Posting to Nostr..."))
	err = postToNostr(sk, message)
	if err != nil {
		fmt.Println(errorStyle.Render("Error posting: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(successStyle.Render("✓ Posted successfully!"))
}

func showUsage() {
	// Check if we have stored credentials
	_, err := getStoredKey()
	if err != nil {
		fmt.Println(titleStyle.Render("Welcome to nos! 🚀"))
		fmt.Println(infoStyle.Render("Usage:"))
		fmt.Println(infoStyle.Render("  nos <message>              - Post a message to Nostr"))
		fmt.Println(infoStyle.Render("  nos relay                  - Manage relay list"))
		fmt.Println(infoStyle.Render("  nos verify                 - Check if your posts are on relays"))
		fmt.Println(infoStyle.Render("  nos reset                  - Reset all data (change account)"))
		fmt.Println(infoStyle.Render("\nFirst time? Run 'nos' with a message to set up your key."))
	} else {
		fmt.Println(errorStyle.Render("Error: Please provide a message to post"))
		fmt.Println(infoStyle.Render("Usage:"))
		fmt.Println(infoStyle.Render("  nos <message>              - Post a message to Nostr"))
		fmt.Println(infoStyle.Render("  nos relay                  - Manage relay list"))
		fmt.Println(infoStyle.Render("  nos verify                 - Check if your posts are on relays"))
		fmt.Println(infoStyle.Render("  nos reset                  - Reset all data (change account)"))
	}
}

func handleReset() {
	// Check if user has stored credentials
	_, err := getStoredKey()
	if err != nil {
		fmt.Println(errorStyle.Render("No stored data found."))
		return
	}

	fmt.Println(titleStyle.Render("Reset nos"))
	fmt.Println(errorStyle.Render("⚠️  This will delete your stored nsec key and relay configuration!"))
	fmt.Println()

	// Confirm with user
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to reset all data?").
				Affirmative("Yes, reset everything").
				Negative("No, cancel").
				Value(&confirm),
		),
	)

	err = form.Run()
	if err != nil || !confirm {
		fmt.Println(infoStyle.Render("Reset cancelled."))
		return
	}

	// Delete nsec key
	err = keyring.Delete(appName, keyringUser)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		fmt.Println(errorStyle.Render("Error deleting key: " + err.Error()))
	}

	// Delete relay list
	err = keyring.Delete(appName, relayListKey)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		fmt.Println(errorStyle.Render("Error deleting relay list: " + err.Error()))
	}

	fmt.Println(successStyle.Render("✓ All data has been reset!"))
	fmt.Println(infoStyle.Render("\nYou can now set up nos with a different account."))
	fmt.Println(infoStyle.Render("Run 'nos <message>' to start fresh."))
}

func handleRelayCommand() {
	// If a subcommand is provided, handle it directly (for backwards compatibility)
	if len(os.Args) >= 3 {
		command := os.Args[2]
		switch command {
		case "list":
			listRelays()
			return
		case "add":
			if len(os.Args) < 4 {
				showRelayMenu()
				return
			}
			addRelay(os.Args[3])
			return
		case "remove":
			if len(os.Args) < 4 {
				showRelayMenu()
				return
			}
			removeRelay(os.Args[3])
			return
		case "reset":
			resetRelays()
			return
		}
	}
	
	// Show interactive menu
	showRelayMenu()
}

func showRelayMenu() {
	for {
		fmt.Println(titleStyle.Render("Relay Management"))
		
		// Get current relay status
		relays := getActiveRelays()
		storedRelays, _ := getStoredRelays()
		usingDefaults := len(storedRelays) == 0
		
		if usingDefaults {
			fmt.Println(infoStyle.Render("Currently using default relays"))
		} else {
			fmt.Println(infoStyle.Render(fmt.Sprintf("Using %d custom relays", len(relays))))
		}
		fmt.Println()

		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("What would you like to do?").
					Options(
						huh.NewOption("View current relays", "list"),
						huh.NewOption("Add a relay", "add"),
						huh.NewOption("Remove a relay", "remove"),
						huh.NewOption("Reset to defaults", "reset"),
						huh.NewOption("Back to main menu", "exit"),
					).
					Value(&choice),
			),
		)

		err := form.Run()
		if err != nil {
			return
		}

		switch choice {
		case "list":
			fmt.Println()
			listRelays()
			fmt.Println()
			fmt.Print("Press Enter to continue...")
			fmt.Scanln()
		case "add":
			interactiveAddRelay()
		case "remove":
			interactiveRemoveRelay()
		case "reset":
			interactiveResetRelays()
		case "exit":
			return
		}
	}
}

func showRelayUsage() {
	fmt.Println(titleStyle.Render("Relay Management"))
	fmt.Println(infoStyle.Render("Usage:"))
	fmt.Println(infoStyle.Render("  nos relay list             - List current relays"))
	fmt.Println(infoStyle.Render("  nos relay add <url>        - Add a relay"))
	fmt.Println(infoStyle.Render("  nos relay remove <url>     - Remove a relay"))
	fmt.Println(infoStyle.Render("  nos relay reset            - Reset to default relays"))
}

func promptForKey() (string, error) {
	var nsec string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your nsec key").
				Description("Your private key (starts with 'nsec1')").
				Placeholder("nsec1...").
				EchoMode(huh.EchoModePassword).
				Value(&nsec).
				Validate(func(str string) error {
					if !strings.HasPrefix(str, "nsec1") {
						return fmt.Errorf("key must start with 'nsec1'")
					}
					// Try to decode it
					_, _, err := nip19.Decode(str)
					if err != nil {
						return fmt.Errorf("invalid nsec key format")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return nsec, nil
}

func getStoredKey() (string, error) {
	secret, err := keyring.Get(appName, keyringUser)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func storeKey(nsec string) error {
	return keyring.Set(appName, keyringUser, nsec)
}

func postToNostr(sk string, content string) error {
	// Create event
	pub, err := nostr.GetPublicKey(sk)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}
	
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      []nostr.Tag{},
		Content:   content,
	}

	// Calculate ID before signing
	ev.ID = ev.GetID()

	// Sign the event
	err = ev.Sign(sk)
	if err != nil {
		return fmt.Errorf("failed to sign event: %v", err)
	}

	// Verify the event is valid
	ok, err := ev.CheckSignature()
	if !ok || err != nil {
		return fmt.Errorf("invalid event signature: %v", err)
	}

	// Show event details for verification
	fmt.Println(infoStyle.Render("Event ID: " + ev.ID))
	fmt.Println(infoStyle.Render("Created at: " + time.Unix(ev.CreatedAt.Time().Unix(), 0).Format(time.RFC3339)))
	fmt.Println(infoStyle.Render("Content: " + content))

	// Get relay list
	relays := getActiveRelays()
	fmt.Println(infoStyle.Render(fmt.Sprintf("Publishing to %d relays...", len(relays))))

	// Connect to relays and publish
	successCount := 0
	failedRelays := []string{}
	
	for _, url := range relays {
		fmt.Printf("  %s Connecting to %s... ", infoStyle.Render("→"), url)
		
		// Create a timeout context for each connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Println(errorStyle.Render("failed: " + err.Error()))
			failedRelays = append(failedRelays, fmt.Sprintf("%s (connection failed: %v)", url, err))
			cancel()
			continue
		}
		
		// Publish with timeout
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = relay.Publish(pubCtx, ev)
		pubCancel()
		
		if err == nil {
			fmt.Println(successStyle.Render("✓ published"))
			successCount++
		} else {
			fmt.Println(errorStyle.Render("failed: " + err.Error()))
			failedRelays = append(failedRelays, fmt.Sprintf("%s (publish failed: %v)", url, err))
		}
		
		relay.Close()
		cancel()
	}

	fmt.Println()
	if successCount == 0 {
		fmt.Println(errorStyle.Render("Failed relays:"))
		for _, fr := range failedRelays {
			fmt.Println(errorStyle.Render("  - " + fr))
		}
		return fmt.Errorf("failed to publish to any relay")
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("Successfully published to %d/%d relays", successCount, len(relays))))
	return nil
}

// Relay management functions
func getActiveRelays() []string {
	relays, err := getStoredRelays()
	if err != nil || len(relays) == 0 {
		return defaultRelays
	}
	return relays
}

func getStoredRelays() ([]string, error) {
	data, err := keyring.Get(appName, relayListKey)
	if err != nil {
		return nil, err
	}
	
	var relays []string
	err = json.Unmarshal([]byte(data), &relays)
	if err != nil {
		return nil, err
	}
	
	return relays, nil
}

func storeRelays(relays []string) error {
	data, err := json.Marshal(relays)
	if err != nil {
		return err
	}
	
	return keyring.Set(appName, relayListKey, string(data))
}

func listRelays() {
	relays := getActiveRelays()
	usingDefaults := false
	
	// Check if we're using defaults
	storedRelays, err := getStoredRelays()
	if err != nil || len(storedRelays) == 0 {
		usingDefaults = true
	}
	
	fmt.Println(titleStyle.Render("Current Relay List"))
	if usingDefaults {
		fmt.Println(infoStyle.Render("(Using default relays)"))
	}
	fmt.Println()
	
	for i, relay := range relays {
		fmt.Printf("%s %d. %s\n", infoStyle.Render("•"), i+1, relay)
	}
}

func addRelay(url string) {
	// Validate URL
	if !strings.HasPrefix(url, "wss://") && !strings.HasPrefix(url, "ws://") {
		fmt.Println(errorStyle.Render("Error: Relay URL must start with wss:// or ws://"))
		os.Exit(1)
	}
	
	// Get current relays
	relays, err := getStoredRelays()
	if err != nil || len(relays) == 0 {
		// If no custom relays, start with defaults
		relays = make([]string, len(defaultRelays))
		copy(relays, defaultRelays)
	}
	
	// Check if already exists
	for _, r := range relays {
		if r == url {
			fmt.Println(infoStyle.Render("Relay already in list: " + url))
			return
		}
	}
	
	// Add relay
	relays = append(relays, url)
	
	// Store updated list
	err = storeRelays(relays)
	if err != nil {
		fmt.Println(errorStyle.Render("Error storing relay list: " + err.Error()))
		os.Exit(1)
	}
	
	fmt.Println(successStyle.Render("✓ Added relay: " + url))
}

func removeRelay(url string) {
	// Get current relays
	relays, err := getStoredRelays()
	if err != nil || len(relays) == 0 {
		fmt.Println(errorStyle.Render("Error: No custom relay list found. Use 'nos relay add' to create one."))
		os.Exit(1)
	}
	
	// Find and remove relay
	found := false
	newRelays := make([]string, 0)
	for _, r := range relays {
		if r != url {
			newRelays = append(newRelays, r)
		} else {
			found = true
		}
	}
	
	if !found {
		fmt.Println(errorStyle.Render("Relay not found in list: " + url))
		os.Exit(1)
	}
	
	if len(newRelays) == 0 {
		fmt.Println(errorStyle.Render("Error: Cannot remove all relays. Use 'nos relay reset' to restore defaults."))
		os.Exit(1)
	}
	
	// Store updated list
	err = storeRelays(newRelays)
	if err != nil {
		fmt.Println(errorStyle.Render("Error storing relay list: " + err.Error()))
		os.Exit(1)
	}
	
	fmt.Println(successStyle.Render("✓ Removed relay: " + url))
}

func resetRelays() {
	// Remove stored relay list
	err := keyring.Delete(appName, relayListKey)
	if err != nil {
		// Ignore error if key doesn't exist
		if !strings.Contains(err.Error(), "not found") {
			fmt.Println(errorStyle.Render("Error resetting relays: " + err.Error()))
			os.Exit(1)
		}
	}
	
	fmt.Println(successStyle.Render("✓ Reset to default relays"))
	listRelays()
}

func interactiveAddRelay() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Add Relay"))
	
	var url string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter relay URL").
				Description("Must start with wss:// or ws://").
				Placeholder("wss://relay.example.com").
				Value(&url).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("URL cannot be empty")
					}
					if !strings.HasPrefix(str, "wss://") && !strings.HasPrefix(str, "ws://") {
						return fmt.Errorf("URL must start with wss:// or ws://")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println(infoStyle.Render("Cancelled"))
		return
	}

	// Get current relays
	relays, err := getStoredRelays()
	if err != nil || len(relays) == 0 {
		// If no custom relays, start with defaults
		relays = make([]string, len(defaultRelays))
		copy(relays, defaultRelays)
	}
	
	// Check if already exists
	for _, r := range relays {
		if r == url {
			fmt.Println(infoStyle.Render("\nRelay already in list: " + url))
			fmt.Print("Press Enter to continue...")
			fmt.Scanln()
			return
		}
	}
	
	// Add relay
	relays = append(relays, url)
	
	// Store updated list
	err = storeRelays(relays)
	if err != nil {
		fmt.Println(errorStyle.Render("\nError storing relay list: " + err.Error()))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}
	
	fmt.Println(successStyle.Render("\n✓ Added relay: " + url))
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
}

func interactiveRemoveRelay() {
	// Get current relays
	relays, err := getStoredRelays()
	if err != nil || len(relays) == 0 {
		fmt.Println()
		fmt.Println(errorStyle.Render("No custom relay list found."))
		fmt.Println(infoStyle.Render("Using default relays. Add a custom relay first."))
		fmt.Print("\nPress Enter to continue...")
		fmt.Scanln()
		return
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("Remove Relay"))
	
	// Create options for select
	options := make([]huh.Option[string], len(relays))
	for i, relay := range relays {
		options[i] = huh.NewOption(relay, relay)
	}
	
	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select relay to remove").
				Options(options...).
				Value(&selected),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Println(infoStyle.Render("Cancelled"))
		return
	}

	if len(relays) == 1 {
		fmt.Println(errorStyle.Render("\nError: Cannot remove all relays."))
		fmt.Println(infoStyle.Render("Use 'Reset to defaults' instead."))
		fmt.Print("\nPress Enter to continue...")
		fmt.Scanln()
		return
	}

	// Remove selected relay
	newRelays := make([]string, 0)
	for _, r := range relays {
		if r != selected {
			newRelays = append(newRelays, r)
		}
	}

	// Store updated list
	err = storeRelays(newRelays)
	if err != nil {
		fmt.Println(errorStyle.Render("\nError storing relay list: " + err.Error()))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}
	
	fmt.Println(successStyle.Render("\n✓ Removed relay: " + selected))
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
}

func interactiveResetRelays() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Reset Relays"))
	fmt.Println(errorStyle.Render("This will restore the default relay list."))
	
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to reset to default relays?").
				Affirmative("Yes, reset").
				Negative("No, cancel").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil || !confirm {
		fmt.Println(infoStyle.Render("\nReset cancelled."))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Remove stored relay list
	err = keyring.Delete(appName, relayListKey)
	if err != nil {
		// Ignore error if key doesn't exist
		if !strings.Contains(err.Error(), "not found") {
			fmt.Println(errorStyle.Render("\nError resetting relays: " + err.Error()))
			fmt.Print("Press Enter to continue...")
			fmt.Scanln()
			return
		}
	}
	
	fmt.Println(successStyle.Render("\n✓ Reset to default relays"))
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
}

func handleVerify() {
	// Get stored key
	nsec, err := getStoredKey()
	if err != nil {
		fmt.Println(errorStyle.Render("No stored key found. Please set up nos first."))
		os.Exit(1)
	}

	// Convert nsec to private key
	var sk string
	_, s, err := nip19.Decode(nsec)
	if err != nil {
		fmt.Println(errorStyle.Render("Error decoding key: " + err.Error()))
		os.Exit(1)
	}
	sk = s.(string)

	// Get public key
	pub, _ := nostr.GetPublicKey(sk)
	npub, _ := nip19.EncodePublicKey(pub)
	
	fmt.Println(titleStyle.Render("Verifying Posts"))
	fmt.Println(infoStyle.Render("Your npub: " + npub))
	fmt.Println(infoStyle.Render("Checking relays for your recent posts..."))
	fmt.Println()

	// Get relay list
	relays := getActiveRelays()
	postsFound := 0

	for _, url := range relays {
		fmt.Printf("%s Checking %s... ", infoStyle.Render("→"), url)
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Println(errorStyle.Render("connection failed"))
			cancel()
			continue
		}

		// Create filter for user's posts
		filter := nostr.Filter{
			Authors: []string{pub},
			Kinds:   []int{nostr.KindTextNote},
			Limit:   5,
		}

		// Subscribe to events
		sub, err := relay.Subscribe(ctx, []nostr.Filter{filter})
		if err != nil {
			fmt.Println(errorStyle.Render("subscription failed"))
			relay.Close()
			cancel()
			continue
		}

		// Collect events with timeout
		events := make([]nostr.Event, 0)
		timeout := time.After(3 * time.Second)
		
	loop:
		for {
			select {
			case ev := <-sub.Events:
				if ev != nil {
					events = append(events, *ev)
				}
			case <-timeout:
				break loop
			}
		}

		sub.Close()
		relay.Close()
		cancel()

		if len(events) > 0 {
			fmt.Println(successStyle.Render(fmt.Sprintf("✓ found %d posts", len(events))))
			postsFound += len(events)
			
			// Show recent posts from this relay
			for i, ev := range events {
				if i >= 3 { // Only show first 3
					break
				}
				timestamp := time.Unix(ev.CreatedAt.Time().Unix(), 0).Format("2006-01-02 15:04:05")
				content := ev.Content
				if len(content) > 50 {
					content = content[:50] + "..."
				}
				fmt.Printf("    %s [%s] %s\n", infoStyle.Render("•"), timestamp, content)
			}
		} else {
			fmt.Println(infoStyle.Render("no posts found"))
		}
	}

	fmt.Println()
	if postsFound == 0 {
		fmt.Println(errorStyle.Render("No posts found on any relay."))
		fmt.Println(infoStyle.Render("This could mean:"))
		fmt.Println(infoStyle.Render("  - Your posts haven't propagated yet (wait a few seconds)"))
		fmt.Println(infoStyle.Render("  - The relays rejected your posts"))
		fmt.Println(infoStyle.Render("  - There's an issue with your key"))
	} else {
	fmt.Println(successStyle.Render(fmt.Sprintf("Total posts found: %d", postsFound)))
	}
}

func interactiveSetup() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Setup Account"))
	fmt.Println(infoStyle.Render("Let's set up your Nostr account."))
	fmt.Println()

	nsec, err := promptForKey()
	if err != nil {
		fmt.Println(errorStyle.Render("\nSetup cancelled."))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Store the key
	err = storeKey(nsec)
	if err != nil {
		fmt.Println(errorStyle.Render("\nError storing key: " + err.Error()))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Show npub
	var sk string
	_, s, _ := nip19.Decode(nsec)
	sk = s.(string)
	pub, _ := nostr.GetPublicKey(sk)
	npub, _ := nip19.EncodePublicKey(pub)

	fmt.Println(successStyle.Render("\n✓ Account setup complete!"))
	fmt.Println(infoStyle.Render("Your npub: " + npub))
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}

func interactivePost() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Post to Nostr"))
	
	var message string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("What would you like to post?").
				Description("Your message will be posted to Nostr").
				Placeholder("Hello Nostr!").
				Value(&message).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("message cannot be empty")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println(infoStyle.Render("\nPost cancelled."))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Get stored key
	nsec, err := getStoredKey()
	if err != nil {
		fmt.Println(errorStyle.Render("\nError: No stored key found."))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Convert nsec to private key
	var sk string
	_, s, err := nip19.Decode(nsec)
	if err != nil {
		fmt.Println(errorStyle.Render("\nError decoding key: " + err.Error()))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}
	sk = s.(string)

	fmt.Println()
	// Post to Nostr
	err = postToNostr(sk, message)
	if err != nil {
		fmt.Println(errorStyle.Render("\nError posting: " + err.Error()))
	} else {
		fmt.Println(successStyle.Render("\n✓ Posted successfully!"))
	}
	
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}

func interactiveReset() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Reset Account"))
	fmt.Println(errorStyle.Render("⚠️  This will delete your stored nsec key and relay configuration!"))
	
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to reset all data?").
				Affirmative("Yes, reset everything").
				Negative("No, cancel").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil || !confirm {
		fmt.Println(infoStyle.Render("\nReset cancelled."))
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	// Delete nsec key
	err = keyring.Delete(appName, keyringUser)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		fmt.Println(errorStyle.Render("\nError deleting key: " + err.Error()))
	}

	// Delete relay list
	err = keyring.Delete(appName, relayListKey)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		fmt.Println(errorStyle.Render("\nError deleting relay list: " + err.Error()))
	}

	fmt.Println(successStyle.Render("\n✓ All data has been reset!"))
	fmt.Println(infoStyle.Render("You can now set up nos with a different account."))
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}
