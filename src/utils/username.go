package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GenerateUsername creates a unique username from fullname and email
func GenerateUsername(db *gorm.DB, fullName, email string) (string, error) {
	// // Clean and prepare the base username from fullname
	// baseUsername := cleanUsername(extractBaseFromFullName(fullName))

	// // If base username is too short, use email prefix
	// if len(baseUsername) < 3 {
	// 	baseUsername = cleanUsername(extractBaseFromEmail(email))
	// }

	// // Ensure minimum length
	// if len(baseUsername) < 3 {
	// 	baseUsername = "user"
	// }

	// // Try the base username first
	// username := baseUsername
	// if isUsernameAvailable(db, username) {
	// 	return username, nil
	// }

	// // Generate variations with numbers
	// for i := 1; i <= 9999; i++ {
	// 	// Try with random number
	// 	randomNum := rand.Intn(9999) + 1
	// 	username = fmt.Sprintf("%s%d", baseUsername, randomNum)
	// 	if isUsernameAvailable(db, username) {
	// 		return username, nil
	// 	}

	// 	// Try with sequential number
	// 	username = fmt.Sprintf("%s%d", baseUsername, i)
	// 	if isUsernameAvailable(db, username) {
	// 		return username, nil
	// 	}
	// }

	// If all else fails, generate completely random username
	return generateRandomUsername(db), nil
}

// cleanUsername removes special characters and converts to lowercase
func cleanUsername(username string) string {
	// Convert to lowercase
	username = strings.ToLower(username)

	// Remove special characters except letters, numbers, and underscores
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	username = reg.ReplaceAllString(username, "")

	// Remove multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	username = reg.ReplaceAllString(username, "_")

	// Remove leading/trailing underscores
	username = strings.Trim(username, "_")

	return username
}

// extractBaseFromFullName extracts a base username from full name
func extractBaseFromFullName(fullName string) string {
	// Split by spaces and take first name
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// extractBaseFromEmail extracts a base username from email
func extractBaseFromEmail(email string) string {
	// Get the part before @
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// isUsernameAvailable checks if username is available in database
func isUsernameAvailable(db *gorm.DB, username string) bool {
	var count int64
	db.Model(&struct {
		Username string `gorm:"column:username"`
	}{}).Where("username = ?", username).Count(&count)
	return count == 0
}

// generateRandomUsername creates a completely random username
func generateRandomUsername(db *gorm.DB) string {
	rand.Seed(time.Now().UnixNano())

	// List of random words for username generation
	words := []string{"user", "member", "player", "gamer", "trader", "bidder", "seller", "buyer"}

	for {
		// Pick random word
		word := words[rand.Intn(len(words))]

		// Add random number
		randomNum := rand.Intn(999999) + 1
		username := fmt.Sprintf("%s%d", word, randomNum)

		if isUsernameAvailable(db, username) {
			return username
		}
	}
}
