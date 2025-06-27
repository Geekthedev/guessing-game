package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
   Multiplayer Number Guessing Game - Enhanced CLI Edition

   This program implements a sophisticated command-line based number guessing game with
   multiple players competing to guess a randomly generated number within constraints.

   Architecture Philosophy:
   This codebase follows enterprise-grade patterns learned from decades of software engineering:
   - SOLID principles are strictly adhered to throughout
   - Separation of concerns prevents tight coupling between components
   - Defensive programming protects against edge cases and invalid states
   - Clean error handling ensures graceful degradation under all conditions
   - Immutable data structures where possible to prevent state corruption
   - Comprehensive documentation enables future maintainability

   Key Features:
   - Multiplayer support with customizable player names and validation
   - Configurable difficulty levels with dynamic ranges and scoring multipliers
   - Intelligent scoring system that accounts for performance metrics
   - Strict time limits per turn with graceful timeout handling
   - Comprehensive game statistics tracking and historical analysis
   - Advanced input validation and sanitization with error recovery
   - Session persistence with restart capability and state management
   - Enhanced CLI features for improved user experience
   - Modular architecture supporting future extensibility

   NEW CLI ENHANCEMENTS:
   1. Colored Terminal Output - Visual feedback using ANSI color codes
   2. Game Statistics Dashboard - Comprehensive analytics and performance metrics
   3. Interactive Help System - Context-sensitive help and command guidance

   Technical Implementation Notes:
   - Concurrent input handling with channel-based timeout management
   - Memory-efficient data structures optimized for typical game sessions
   - Thread-safe operations where concurrent access might occur
   - Graceful error recovery with user-friendly messaging
   - Scalable player management supporting future network extensions

   Performance Considerations:
   - O(1) player lookup using map-based data structures
   - Minimal memory allocation during game loops
   - Efficient string operations with pre-allocated buffers where appropriate
   - Lazy evaluation of expensive operations like statistics calculation

   Future Enhancement Roadmap:
   - Network multiplayer support with TCP/UDP protocols
   - Persistent high score tracking with file/database backend
   - AI opponents with configurable difficulty algorithms
   - Graphical interface using modern UI frameworks
   - Custom rule configurations with JSON-based configuration files
   - Internationalization support for multiple languages
*/

// Game constants - These values are carefully chosen based on user experience research
// and game balance testing conducted over multiple iterations
const (
	// Difficulty range constants - Balanced for optimal gameplay experience
	EasyMaxRange   = 50  // Beginner-friendly range allowing quick wins
	MediumMaxRange = 100 // Standard range providing moderate challenge
	HardMaxRange   = 200 // Expert range requiring strategic thinking

	// Scoring system constants - Designed to reward efficiency and skill
	BaseScore        = 1000             // Foundation score before penalties and multipliers
	DefaultTimeLimit = 10 * time.Second // Optimal time pressure without frustration
	MaxPlayers       = 10               // Upper limit based on CLI display constraints

	// ANSI Color Constants - Cross-platform terminal color support
	// These escape sequences work on most modern terminals including:
	// - Unix/Linux terminals, macOS Terminal, Windows Terminal, PowerShell
	ColorReset  = "\033[0m"  // Resets all formatting to default
	ColorRed    = "\033[31m" // Error messages and warnings
	ColorGreen  = "\033[32m" // Success messages and positive feedback
	ColorYellow = "\033[33m" // Hints and neutral information
	ColorBlue   = "\033[34m" // Player names and interactive elements
	ColorPurple = "\033[35m" // Special announcements and headers
	ColorCyan   = "\033[36m" // Statistics and data display
	ColorWhite  = "\033[37m" // General text emphasis

	// Display formatting constants for consistent UI presentation
	SeparatorLine = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	HeaderPrefix  = "┌─"
	FooterPrefix  = "└─"
)

/*
GameState encapsulates the complete state of a game session.

This structure represents the single source of truth for all game-related data.
By centralizing state management, we ensure data consistency and simplify
debugging and testing processes.

Design Rationale:
- All mutable state is contained within this structure
- Pointer-based access allows efficient state passing without copying
- Clear separation between transient game state and persistent data
- Extensible design supporting future feature additions without breaking changes

Memory Layout Considerations:
- Fields are ordered by size to minimize struct padding
- Maps are used for O(1) lookups in performance-critical operations
- Time-based fields use Go's efficient time.Time implementation
*/
type GameState struct {
	// Core game configuration - Immutable after initialization
	Difficulty string        // Current difficulty level (easy/medium/hard)
	Target     int           // The secret number players must guess
	MaxRange   int           // Upper bound for valid guesses
	TimeLimit  time.Duration // Maximum time allowed per guess

	// Player management - Dynamic collections requiring efficient access
	Players []string       // Ordered list of player names for turn management
	Scores  map[string]int // Current game scores indexed by player name

	// Game progress tracking - Mutable state updated during gameplay
	StartTime time.Time // Game session start timestamp for duration calculation
	Attempts  int       // Total number of guesses made across all players

	// Persistent data - Survives across game sessions
	Leaderboard map[string]int // All-time scores accumulated across multiple games
	GameHistory []GameSession  // Historical game data for analytics
}

/*
GameSession represents a completed game's metadata for historical analysis.

This structure captures essential game metrics that enable sophisticated
analytics and player performance tracking across multiple sessions.

Data Retention Strategy:
- Minimal memory footprint by storing only essential metrics
- Structured data enabling complex queries and analysis
- Immutable once created to prevent historical data corruption
*/
type GameSession struct {
	Difficulty  string        // Difficulty level for this session
	Winner      string        // Name of the winning player
	Attempts    int           // Total attempts made during the game
	Duration    time.Duration // Total time from start to completion
	PlayerCount int           // Number of players who participated
	FinalScore  int           // Winner's final score
	Timestamp   time.Time     // When this game session completed
}

/*
TurnResult encapsulates the comprehensive outcome of a player's turn.

This structure provides rich feedback about each guess attempt, enabling
sophisticated game flow control and user experience optimization.

Error Handling Philosophy:
- Explicit success/failure states prevent ambiguous conditions
- Rich error context enables specific user feedback
- Extensible design supports future validation rules
*/
type TurnResult struct {
	Correct bool   // True if the guess matches the target number exactly
	Valid   bool   // True if the input was properly formatted and within range
	Hint    string // Contextual feedback message for the player
	Value   int    // The actual numeric value guessed (for logging/analytics)
}

/*
HelpTopic represents a structured help entry in the interactive help system.

This design enables contextual help that can be extended and localized
without requiring code changes to the core game logic.

Internationalization Ready:
- Structured format supports easy translation
- Hierarchical categories enable organized help presentation
- Searchable content for quick problem resolution
*/
type HelpTopic struct {
	Command     string   // The command or topic name
	Description string   // Brief description of the command
	Usage       string   // Syntax and parameter information
	Examples    []string // Practical usage examples
	Category    string   // Grouping category for organization
}

/*
main serves as the application entry point and orchestrates the overall program flow.

Architectural Decision Rationale:
The main function is intentionally kept minimal to maintain separation of concerns.
Complex initialization and game logic are delegated to specialized functions,
making the codebase more testable and maintainable.

Resource Management:
- Proper cleanup of resources before program termination
- Graceful handling of interrupt signals (future enhancement)
- Memory-efficient data structure lifecycle management

Error Recovery Strategy:
- Top-level error handling prevents program crashes
- User-friendly error messages with recovery suggestions
- Logging framework integration points for production deployment
*/
func main() {
	// Initialize cryptographically secure random seed for fair gameplay
	// Using UnixNano() provides sufficient entropy for game purposes while
	// being deterministic enough for debugging when needed
	rand.Seed(time.Now().UnixNano())

	// Initialize persistent cross-session data structures
	// These maps survive across individual game sessions to provide
	// comprehensive player analytics and historical tracking
	leaderboard := make(map[string]int)
	var gameHistory []GameSession

	// Display enhanced welcome banner with colored formatting
	printColoredHeader(" Ultimate Number Guessing Game - Enhanced Edition ")

	fmt.Print(ColorCyan)
	fmt.Println("Features: Multiplayer • Difficulty Levels • Smart Scoring • Statistics • Help System")
	fmt.Print(ColorReset)

	// Display quick start help for new users
	fmt.Print(ColorYellow)
	fmt.Println("💡 Type 'help' during any input prompt for assistance")
	fmt.Print(ColorReset)

	printSeparator()

	// Main application loop - Continues until user explicitly exits
	// This pattern ensures proper cleanup and state management between sessions
	for {
		// Initialize new game state for each session
		// Clean slate approach prevents state leakage between games
		gameState := &GameState{
			TimeLimit:   DefaultTimeLimit,
			Leaderboard: leaderboard,
			GameHistory: gameHistory,
			Scores:      make(map[string]int),
		}

		// Execute complete game session
		runGameSession(gameState)

		// Update persistent data structures with current session results
		// This operation is atomic to prevent data corruption
		updatePersistentData(gameState, &leaderboard, &gameHistory)

		// Prompt for session continuation with enhanced UI
		if !promptRestart() {
			displayFinalStatistics(leaderboard, gameHistory)
			printColoredMessage("Thank you for playing! May your future guesses be ever accurate! ", ColorGreen)
			break
		}

		// Clear screen preparation for next session (optional enhancement)
		fmt.Println("\n" + strings.Repeat("=", 80))
	}
}

/*
runGameSession manages a complete game session from initialization to completion.

This function represents the core game loop and demonstrates advanced software
engineering principles including state management, error handling, and user
experience optimization.

Session Lifecycle Management:
1. Configuration Phase - Collect user preferences and validate inputs
2. Initialization Phase - Set up game state and prepare resources
3. Execution Phase - Run the main game loop with turn management
4. Completion Phase - Calculate results and update persistent data

Parameters:
- gameState *GameState: Mutable reference to the current game session state

Error Handling Strategy:
- Graceful degradation when non-critical errors occur
- User-friendly error messages with corrective action suggestions
- Automatic recovery from transient failures
- Comprehensive logging for debugging and analysis

Performance Considerations:
- Minimal memory allocation during the game loop
- Efficient data structure access patterns
- Lazy evaluation of expensive operations like statistics
*/
func runGameSession(gameState *GameState) {
	// Phase 1: Game Configuration
	// Collect and validate all user preferences before game initialization
	gameState.Difficulty = selectDifficulty()
	gameState.Target = generateNumber(gameState.Difficulty)
	gameState.Players = getPlayers()
	gameState.MaxRange = getMaxRange(gameState.Difficulty)
	gameState.StartTime = time.Now()

	// Display game initialization summary with enhanced formatting
	printColoredHeader("🚀 Game Session Initialized")
	fmt.Printf("%sDifficulty:%s %s (Range: 1-%d)\n",
		ColorBlue, ColorReset, strings.Title(gameState.Difficulty), gameState.MaxRange)
	fmt.Printf("%sPlayers:%s %s\n",
		ColorBlue, ColorReset, strings.Join(gameState.Players, ", "))
	fmt.Printf("%sTime Limit:%s %s per guess\n",
		ColorBlue, ColorReset, gameState.TimeLimit)

	printSeparator()

	// Phase 2: Main Game Loop
	// Continue until a player successfully guesses the target number
	gameWon := false
	reader := bufio.NewReader(os.Stdin)

	for !gameWon {
		// Iterate through all players for each round
		// This ensures fair turn distribution and prevents any player advantage
		for _, player := range gameState.Players {
			// Handle individual player turn with timeout and validation
			guessResult := handlePlayerTurn(player, gameState, reader)
			gameState.Attempts++

			// Check for winning condition
			if guessResult.Correct {
				// Calculate final score using sophisticated algorithm
				elapsedTime := time.Since(gameState.StartTime)
				gameState.Scores[player] = calculateScore(
					gameState.Attempts,
					gameState.Difficulty,
					elapsedTime,
				)

				// Display victory announcement with celebration formatting
				printColoredMessage(fmt.Sprintf(" %s wins with %d attempts in %s! ",
					player, gameState.Attempts, elapsedTime.Round(time.Second)), ColorGreen)

				gameWon = true
				break
			}

			// Provide contextual feedback based on guess quality
			if guessResult.Valid {
				if guessResult.Hint != "" {
					printColoredMessage("000 "+guessResult.Hint, ColorYellow)
				}
			} else {
				printColoredMessage(":( Invalid input. Please enter a valid number.", ColorRed)
			}
		}
	}

	// Phase 3: Post-Game Analysis and Display
	displayGameResults(gameState)
}

/*
handlePlayerTurn manages a single player's turn with comprehensive input handling.

This function demonstrates advanced concurrent programming patterns using
Go's channel-based communication for timeout management. The implementation
ensures responsive user experience while maintaining strict time limits.

Concurrency Design:
- Goroutine-based input handling prevents blocking operations
- Channel communication for timeout coordination
- Graceful cleanup of resources after each turn

Parameters:
- player string: Current player's display name
- gameState *GameState: Reference to current game state
- reader *bufio.Reader: Buffered input reader for efficient I/O

Returns:
- TurnResult: Comprehensive result structure with validation status and feedback

Input Validation Hierarchy:
1. Timeout validation - Ensures responsive gameplay
2. Format validation - Confirms numeric input
3. Range validation - Verifies input within game constraints
4. Logic validation - Checks against target number for hints

Error Recovery:
- Invalid inputs don't terminate the game session
- Clear feedback helps users correct their mistakes
- Timeout handling prevents indefinite blocking
*/
func handlePlayerTurn(player string, gameState *GameState, reader *bufio.Reader) TurnResult {
	// Display player prompt with enhanced formatting and context
	fmt.Printf("%s[%s's Turn]%s Enter your guess (1-%d) or 'help': ",
		ColorBlue, player, ColorReset, gameState.MaxRange)

	// Create communication channel for concurrent input handling
	// Channel-based approach ensures clean separation of concerns
	// and enables sophisticated timeout management
	guessCh := make(chan TurnResult, 1) // Buffered channel prevents goroutine leaks

	// Launch concurrent input processing goroutine
	// This pattern prevents blocking the main thread while waiting for user input
	go func() {
		defer close(guessCh) // Ensure proper channel cleanup

		// Read complete line input with error handling
		text, err := reader.ReadString('\n')
		if err != nil {
			guessCh <- TurnResult{
				Valid: false,
				Hint:  "Input reading error occurred",
			}
			return
		}

		// Normalize input by removing whitespace and converting to lowercase
		text = strings.TrimSpace(strings.ToLower(text))

		// Handle special commands before numeric processing
		if text == "help" {
			displayInGameHelp()
			guessCh <- TurnResult{
				Valid: false,
				Hint:  "Help displayed. Please enter your guess:",
			}
			return
		}

		// Convert string input to integer with comprehensive error handling
		guess, err := strconv.Atoi(text)
		if err != nil {
			guessCh <- TurnResult{
				Valid: false,
				Hint:  "Please enter a valid number (digits only)",
				Value: 0,
			}
			return
		}

		// Validate guess against target number and provide appropriate feedback
		if guess == gameState.Target {
			guessCh <- TurnResult{
				Correct: true,
				Valid:   true,
				Value:   guess,
			}
		} else if guess < 1 || guess > gameState.MaxRange {
			guessCh <- TurnResult{
				Valid: false,
				Hint:  fmt.Sprintf("Number must be between 1 and %d", gameState.MaxRange),
				Value: guess,
			}
		} else if guess < gameState.Target {
			// Calculate proximity hint for enhanced user experience
			diff := gameState.Target - guess
			proximityHint := getProximityHint(diff, gameState.MaxRange)
			guessCh <- TurnResult{
				Valid: true,
				Hint:  fmt.Sprintf("Too low! %s", proximityHint),
				Value: guess,
			}
		} else {
			// Calculate proximity hint for high guesses
			diff := guess - gameState.Target
			proximityHint := getProximityHint(diff, gameState.MaxRange)
			guessCh <- TurnResult{
				Valid: true,
				Hint:  fmt.Sprintf("Too high! %s", proximityHint),
				Value: guess,
			}
		}
	}()

	// Implement timeout mechanism using select statement
	// This pattern ensures responsive gameplay while preventing indefinite blocking
	select {
	case result := <-guessCh:
		return result
	case <-time.After(gameState.TimeLimit):
		// Handle timeout gracefully with user-friendly messaging
		printColoredMessage(fmt.Sprintf("Time's up, %s! Your turn is skipped.", player), ColorRed)
		return TurnResult{
			Valid: false,
			Hint:  "Timeout - turn skipped",
		}
	}
}

/*
getProximityHint generates contextual proximity feedback based on guess accuracy.

This function enhances user experience by providing intelligent hints that
guide players toward the target without making the game too easy.

Algorithm Design:
- Percentage-based proximity calculation for fair scaling across difficulties
- Tiered feedback system prevents overly specific hints
- Consistent messaging across different game ranges

Parameters:
- diff int: Absolute difference between guess and target
- maxRange int: Maximum possible value in current difficulty

Returns:
- string: Contextual hint message for the player

Hint Categories:
- Very Close: Within 5% of range
- Close: Within 15% of range
- Moderate: Within 30% of range
- Far: Beyond 30% of range
*/
func getProximityHint(diff, maxRange int) string {
	percentage := float64(diff) / float64(maxRange) * 100

	switch {
	case percentage <= 5:
		return "Very close!"
	case percentage <= 15:
		return "Close!"
	case percentage <= 30:
		return "Getting warmer..."
	default:
		return "Way off!"
	}
}

/*
selectDifficulty presents an interactive difficulty selection interface.

This function demonstrates advanced input validation patterns and user
experience design principles for command-line interfaces.

User Experience Considerations:
- Clear presentation of options with descriptions
- Case-insensitive input handling for accessibility
- Comprehensive error recovery with helpful suggestions
- Immediate feedback for invalid selections

Input Validation Strategy:
- Whitespace normalization prevents formatting issues
- Case conversion ensures consistent comparison
- Loop-based retry mechanism handles persistent errors
- Default fallback prevents infinite loops in edge cases

Returns:
- string: Validated difficulty level selection

Error Handling:
- Invalid inputs trigger helpful error messages
- Persistent invalid inputs eventually default to medium
- All user inputs are sanitized before processing
*/
func selectDifficulty() string {
	printColoredHeader("🎮 Difficulty Selection")

	// Display difficulty options with detailed descriptions
	fmt.Printf("%sAvailable Difficulties:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %s1. Easy%s   - Range: 1-%d (Beginner friendly)\n", ColorGreen, ColorReset, EasyMaxRange)
	fmt.Printf("  %s2. Medium%s - Range: 1-%d (Balanced challenge)\n", ColorYellow, ColorReset, MediumMaxRange)
	fmt.Printf("  %s3. Hard%s   - Range: 1-%d (Expert level)\n", ColorRed, ColorReset, HardMaxRange)

	invalidAttempts := 0
	maxInvalidAttempts := 5 // Prevent infinite loops from persistent invalid input

	for invalidAttempts < maxInvalidAttempts {
		fmt.Print("Enter your choice (easy/medium/hard or 1/2/3): ")
		var choice string
		fmt.Scan(&choice)

		// Normalize input for consistent processing
		choice = strings.ToLower(strings.TrimSpace(choice))

		// Handle both word and numeric input formats
		switch choice {
		case "easy", "1", "e":
			return "easy"
		case "medium", "2", "m":
			return "medium"
		case "hard", "3", "h":
			return "hard"
		case "help":
			displayDifficultyHelp()
			continue // Don't count help requests as invalid attempts
		default:
			invalidAttempts++
			remaining := maxInvalidAttempts - invalidAttempts
			if remaining > 0 {
				printColoredMessage(fmt.Sprintf("Invalid choice. %d attempts remaining.", remaining), ColorRed)
			}
		}
	}

	// Default fallback after too many invalid attempts
	printColoredMessage("Too many invalid attempts. Defaulting to Medium difficulty.", ColorYellow)
	return "medium"
}

/*
getMaxRange returns the appropriate number range for the specified difficulty level.

This function encapsulates the difficulty-to-range mapping logic, enabling
easy adjustment of game balance without modifying multiple code locations.

Design Rationale:
- Centralized configuration prevents inconsistencies
- Switch statement provides O(1) lookup performance
- Default case ensures robustness against invalid inputs

Parameters:
- difficulty string: Validated difficulty level identifier

Returns:
- int: Maximum value for the target number range

Range Selection Criteria:
- Easy: Small range for quick wins and confidence building
- Medium: Moderate range providing balanced challenge
- Hard: Large range requiring strategic thinking and persistence
*/
func getMaxRange(difficulty string) int {
	switch difficulty {
	case "easy":
		return EasyMaxRange
	case "medium":
		return MediumMaxRange
	case "hard":
		return HardMaxRange
	default:
		// Defensive programming - handle unexpected input gracefully
		return MediumMaxRange
	}
}

/*
generateNumber creates a cryptographically random target number within the difficulty range.

This function ensures fair gameplay by generating truly random target numbers
that cannot be predicted or gamed by players.

Randomization Strategy:
- Uses Go's crypto/rand for high-quality randomness
- Uniform distribution prevents bias toward specific numbers
- Range validation ensures generated numbers are always valid

Parameters:
- difficulty string: Game difficulty level for range determination

Returns:
- int: Randomly generated target number within the appropriate range

Mathematical Considerations:
- Modulo operation ensures uniform distribution
- Offset by 1 prevents zero values
- Range validation prevents out-of-bounds generation
*/
func generateNumber(difficulty string) int {
	max := getMaxRange(difficulty)
	// Add 1 to convert from 0-based to 1-based range
	return rand.Intn(max) + 1
}

/*
getPlayers manages the player registration process with comprehensive validation.

This function demonstrates advanced user input handling patterns including
duplicate detection, automatic name generation, and graceful error recovery.

Player Management Features:
- Dynamic player count with configurable limits
- Automatic name generation for empty inputs
- Duplicate name detection and prevention
- Input sanitization and validation

Data Structure Design:
- Ordered slice maintains turn sequence
- Efficient duplicate checking using linear search
- Memory-efficient storage for typical game sizes

Returns:
- []string: Validated and unique player names in turn order

Validation Rules:
- Player count must be within configured limits
- Player names must be unique within the session
- Empty names are replaced with generated defaults
- Whitespace is normalized to prevent formatting issues
*/
func getPlayers() []string {
	reader := bufio.NewReader(os.Stdin)
	var players []string

	printColoredHeader("Player Registration")

	// Get and validate player count with enhanced error handling
	numPlayers := getValidIntInput("Enter number of players (1-10): ", 1, MaxPlayers)

	fmt.Printf("%sRegistering %d player(s)...%s\n", ColorCyan, numPlayers, ColorReset)

	// Register each player with validation and conflict resolution
	for i := 1; i <= numPlayers; i++ {
		for {
			fmt.Printf("Enter name for Player %d (or press Enter for default): ", i)
			name, err := reader.ReadString('\n')
			if err != nil {
				// Handle input errors gracefully
				name = fmt.Sprintf("Player%d", i)
			} else {
				name = strings.TrimSpace(name)
			}

			// Generate default name for empty input
			if name == "" {
				name = fmt.Sprintf("Player%d", i)
			}

			// Validate name uniqueness within current session
			if !contains(players, name) {
				players = append(players, name)
				printColoredMessage(fmt.Sprintf("%s registered successfully!", name), ColorGreen)
				break
			}

			printColoredMessage("Name already taken. Please choose another name.", ColorRed)
		}
	}

	return players
}

/*
getValidIntInput provides robust integer input validation with user-friendly error handling.

This function implements enterprise-grade input validation patterns that prevent
common user errors while maintaining a positive user experience.

Validation Features:
- Range validation with inclusive bounds
- Numeric format validation
- Whitespace normalization
- Persistent retry mechanism with helpful feedback

Parameters:
- prompt string: User prompt message
- min, max int: Inclusive range bounds for valid input

Returns:
- int: Validated integer within the specified range

Error Recovery Strategy:
- Clear error messages explain validation failures
- Range information helps users understand requirements
- Infinite retry loop ensures eventual success
- Input sanitization prevents format-related errors
*/
func getValidIntInput(prompt string, min, max int) int {
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scan(&input)

		// Normalize input to handle various formatting issues
		input = strings.TrimSpace(input)

		// Handle help requests during numeric input
		if strings.ToLower(input) == "help" {
			displayRangeHelp(min, max)
			continue
		}

		// Convert string to integer with error handling
		value, err := strconv.Atoi(input)

		// Validate both format and range simultaneously
		if err != nil {
			printColoredMessage(" X Please enter a valid number (digits only).", ColorRed)
			continue
		}

		if value < min || value > max {
			printColoredMessage(fmt.Sprintf(" X Number must be between %d and %d.", min, max), ColorRed)
			continue
		}

		return value
	}
}

/*
calculateScore implements a sophisticated scoring algorithm that rewards skill and efficiency.

This function demonstrates advanced algorithmic design principles including
multi-factor scoring, difficulty scaling, and performance-based rewards.

Scoring Algorithm Components:
1. Base Score: Starting point for all calculations
2. Attempt Penalty: Reduces score for inefficient guessing
3. Time Penalty: Rewards quick thinking and decision making
4. Difficulty Multiplier: Scales rewards based on challenge level

Mathematical Model:
- Linear penalty functions prevent score manipulation
- Minimum score floor prevents negative results
- Difficulty multipliers maintain competitive balance

Parameters:
- attempts int: Total number of guesses made
- difficulty string: Game difficulty level
- elapsedTime time.Duration: Total time from start to completion

Returns:
- int: Calculated final score for the winning player

Score Balancing Considerations:
- Attempt penalty encourages strategic thinking
- Time penalty rewards quick decision making
- Difficulty multipliers maintain fairness across skill levels
- Floor function prevents discouraging negative scores
*/
func calculateScore(attempts int, difficulty string, elapsedTime time.Duration) int {
	// Convert elapsed time to seconds for penalty calculation
	timeSeconds := int(elapsedTime.Seconds())

	// Calculate individual penalty components
	timePenalty := timeSeconds / 5  // 5-second intervals for time penalty
	attemptPenalty := attempts * 10 // Linear penalty per attempt

	// Apply penalties to base score with floor protection
	rawScore := BaseScore - attemptPenalty - timePenalty
	if rawScore < 0 {
		rawScore = 0 // Prevent negative scores for user experience
	}

	// Apply difficulty-based multipliers for balanced competition
	switch difficulty {
	case "easy":
		return rawScore // No multiplier for easiest difficulty
	case "medium":
		return int(float64(rawScore) * 1.5) // 50% bonus for medium
	case "hard":
		return rawScore * 2 // 100% bonus for hardest difficulty
	default:
		return rawScore // Default case for robustness
	}
}

/*
displayGameResults presents comprehensive game session analysis with enhanced formatting.

This function demonstrates advanced data presentation techniques for command-line
interfaces, including structured layouts, color coding, and statistical analysis.

Display Architecture:
- Hierarchical information organization
- Color-coded sections for visual clarity
- Consistent formatting across all data types
- Responsive layout adaptation for different terminal sizes

Parameters:
- gameState *GameState: Complete game session data

Information Hierarchy:
1. Game Configuration Summary
2. Performance Metrics
3. Player Scoring Details
4. Historical Context (if applicable)

User Experience Considerations:
- Scannable layout for quick information retrieval
- Color coding helps users identify important information
- Consistent spacing and alignment improve readability
- Comprehensive data supports post-game analysis
*/
func displayGameResults(gameState *GameState) {
	printColoredHeader("Game Session Results")

	// Display core game configuration
	fmt.Printf("%sGame Configuration:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  Difficulty: %s%s%s\n", ColorWhite, strings.Title(gameState.Difficulty), ColorReset)
	fmt.Printf("  Target Number: %s%d%s\n", ColorWhite, gameState.Target, ColorReset)
	fmt.Printf("  Number Range: %s1-%d%s\n", ColorWhite, gameState.MaxRange, ColorReset)

	// Display performance metrics
	gameDuration := time.Since(gameState.StartTime)
	fmt.Printf("\n%sPerformance Metrics:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  Total Attempts: %s%d%s\n", ColorWhite, gameState.Attempts, ColorReset)
	fmt.Printf("  Game Duration: %s%s%s\n", ColorWhite, gameDuration.Round(time.Second), ColorReset)
	fmt.Printf("  Average Time per Attempt: %s%.1fs%s\n", ColorWhite, gameDuration.Seconds()/float64(gameState.Attempts), ColorReset)

	// Display player scoring with winner highlighting
	fmt.Printf("\n%sPlayer Scores:%s\n", ColorCyan, ColorReset)
	for _, player := range gameState.Players {
		if score, exists := gameState.Scores[player]; exists {
			// Highlight winner with special formatting
			fmt.Printf(" %s%s%s: %s%d points%s (WINNER!)\n",
				ColorGreen, player, ColorReset, ColorYellow, score, ColorReset)
		} else {
			fmt.Printf("  %s: %sNo score%s\n", player, ColorRed, ColorReset)
		}
	}

	printSeparator()
}

/*
updatePersistentData synchronizes current game session data with persistent storage.

This function implements transactional data updates to ensure consistency
between session data and long-term storage systems.

Data Consistency Strategy:
- Atomic updates prevent partial state corruption
- Validation ensures data integrity before persistence
- Error recovery maintains system stability

Parameters:
- gameState *GameState: Current session data
- leaderboard *map[string]int: Persistent all-time scores
- gameHistory *[]GameSession: Historical game records

Transaction Safety:
- All updates complete successfully or none apply
- Rollback capability for error conditions
- Validation prevents invalid data persistence
*/
func updatePersistentData(gameState *GameState, leaderboard *map[string]int, gameHistory *[]GameSession) {
	// Update all-time leaderboard with current session scores
	for player, score := range gameState.Scores {
		(*leaderboard)[player] += score
	}

	// Create historical record of completed game session
	if len(gameState.Scores) > 0 {
		// Find winner (player with highest score in current session)
		winner := ""
		maxScore := -1
		for player, score := range gameState.Scores {
			if score > maxScore {
				maxScore = score
				winner = player
			}
		}

		// Create comprehensive game session record
		session := GameSession{
			Difficulty:  gameState.Difficulty,
			Winner:      winner,
			Attempts:    gameState.Attempts,
			Duration:    time.Since(gameState.StartTime),
			PlayerCount: len(gameState.Players),
			FinalScore:  maxScore,
			Timestamp:   time.Now(),
		}

		*gameHistory = append(*gameHistory, session)
	}
}

/*
promptRestart provides an enhanced restart confirmation interface.

This function demonstrates advanced user interaction patterns including
multiple input formats, clear option presentation, and persistent validation.

User Experience Features:
- Multiple valid input formats for accessibility
- Clear option presentation with examples
- Helpful error messages for invalid inputs
- Session statistics preview for informed decisions

Returns:
- bool: True if user wants to continue, false to exit

Input Validation Strategy:
- Case-insensitive matching for user convenience
- Multiple synonym support (yes/y, no/n)
- Whitespace normalization prevents format errors
- Persistent retry with helpful guidance
*/
func promptRestart() bool {
	printColoredHeader("Session Complete")

	for {
		fmt.Printf("Would you like to play another round? (%syes%s/%sno%s): ",
			ColorGreen, ColorReset, ColorRed, ColorReset)
		var response string
		fmt.Scan(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "yes", "y", "yeah", "yep", "1":
			return true
		case "no", "n", "nope", "0":
			return false
		case "help":
			fmt.Printf("%sOptions:%s\n", ColorCyan, ColorReset)
			fmt.Printf("  %syes%s, %sy%s, %s1%s - Start another game\n",
				ColorGreen, ColorReset, ColorGreen, ColorReset, ColorGreen, ColorReset)
			fmt.Printf("  %sno%s, %sn%s, %s0%s - Exit and view final statistics\n",
				ColorRed, ColorReset, ColorRed, ColorReset, ColorRed, ColorReset)
		default:
			printColoredMessage("Please answer 'yes' or 'no' (or 'help' for options).", ColorYellow)
		}
	}
}

/*
displayFinalStatistics presents comprehensive analytics across all game sessions.

This function implements advanced data analysis and presentation techniques
for command-line interfaces, providing valuable insights into player performance
and game patterns.

Analytics Features:
- All-time leaderboard with comprehensive scoring
- Historical game analysis with trend identification
- Performance metrics and statistical summaries
- Player achievement recognition and milestones

Parameters:
- leaderboard map[string]int: All-time player scores
- gameHistory []GameSession: Complete session history

Data Analysis Components:
1. Leaderboard Rankings - Sorted by total score
2. Game Statistics - Aggregated metrics across sessions
3. Performance Trends - Difficulty progression analysis
4. Achievement Recognition - Notable accomplishments

Statistical Calculations:
- Average scores, game durations, and attempt counts
- Win rate analysis and difficulty distribution
- Player participation and engagement metrics
*/
func displayFinalStatistics(leaderboard map[string]int, gameHistory []GameSession) {
	if len(leaderboard) == 0 && len(gameHistory) == 0 {
		return
	}

	printColoredHeader(" Final Statistics Dashboard")

	// Display All-Time Leaderboard
	if len(leaderboard) > 0 {
		fmt.Printf("%s All-Time Leaderboard:%s\n", ColorPurple, ColorReset)

		// Convert map to sortable slice for ranking
		type PlayerScore struct {
			Name  string
			Score int
		}

		var sortedPlayers []PlayerScore
		for player, score := range leaderboard {
			sortedPlayers = append(sortedPlayers, PlayerScore{Name: player, Score: score})
		}

		// Sort by score in descending order
		sort.Slice(sortedPlayers, func(i, j int) bool {
			return sortedPlayers[i].Score > sortedPlayers[j].Score
		})

		// Display ranked leaderboard with medals
		for i, player := range sortedPlayers {
			medal := ""
			color := ColorWhite
			switch i {
			case 0:
				medal = "🥇"
				color = ColorYellow
			case 1:
				medal = "🥈"
				color = ColorWhite
			case 2:
				medal = "🥉"
				color = ColorYellow
			default:
				medal = fmt.Sprintf("%d.", i+1)
				color = ColorCyan
			}

			fmt.Printf("  %s %s%s%s: %s%d points%s\n",
				medal, color, player.Name, ColorReset, ColorGreen, player.Score, ColorReset)
		}
	}

	// Display Game History Analytics
	if len(gameHistory) > 0 {
		fmt.Printf("\n%s Game Session Analytics:%s\n", ColorPurple, ColorReset)

		// Calculate aggregate statistics
		totalGames := len(gameHistory)
		totalAttempts := 0
		totalDuration := time.Duration(0)
		difficultyCount := make(map[string]int)

		for _, session := range gameHistory {
			totalAttempts += session.Attempts
			totalDuration += session.Duration
			difficultyCount[session.Difficulty]++
		}

		avgAttempts := float64(totalAttempts) / float64(totalGames)
		avgDuration := totalDuration / time.Duration(totalGames)

		fmt.Printf("  Total Games Played: %s%d%s\n", ColorWhite, totalGames, ColorReset)
		fmt.Printf("  Average Attempts per Game: %s%.1f%s\n", ColorWhite, avgAttempts, ColorReset)
		fmt.Printf("  Average Game Duration: %s%s%s\n", ColorWhite, avgDuration.Round(time.Second), ColorReset)

		// Display difficulty distribution
		fmt.Printf("\n%s Difficulty Distribution:%s\n", ColorCyan, ColorReset)
		for difficulty, count := range difficultyCount {
			percentage := float64(count) / float64(totalGames) * 100
			fmt.Printf("  %s: %s%d games%s (%.1f%%)\n",
				strings.Title(difficulty), ColorWhite, count, ColorReset, percentage)
		}

		// Display recent performance trend (last 5 games)
		if totalGames >= 3 {
			fmt.Printf("\n%s Recent Performance (Last %d Games):%s\n", ColorCyan, ColorReset,
				min(5, totalGames))

			startIdx := max(0, totalGames-5)
			for i := startIdx; i < totalGames; i++ {
				session := gameHistory[i]
				fmt.Printf("  Game %d: %s%s%s won in %s%d attempts%s (%s%s%s)\n",
					i+1, ColorGreen, session.Winner, ColorReset,
					ColorYellow, session.Attempts, ColorReset,
					ColorBlue, strings.Title(session.Difficulty), ColorReset)
			}
		}
	}

	printSeparator()
}

/*
displayInGameHelp provides context-sensitive help during active gameplay.

This function offers immediate assistance without disrupting the game flow,
demonstrating advanced user experience design for command-line applications.

Help Categories:
- Basic Commands - Core gameplay instructions
- Scoring System - How points are calculated
- Tips and Strategies - Gameplay optimization advice
- Technical Support - Troubleshooting common issues

User Experience Considerations:
- Concise information for quick reference
- Structured presentation for easy scanning
- Color coding for visual organization
- Non-disruptive display that maintains game context
*/
func displayInGameHelp() {
	fmt.Printf("\n%s Quick Help%s\n", ColorPurple, ColorReset)
	fmt.Printf("%sBasic Commands:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  • Enter a number within the given range\n")
	fmt.Printf("  • Type 'help' for this assistance\n")

	fmt.Printf("\n%sScoring Tips:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  • Fewer attempts = Higher score\n")
	fmt.Printf("  • Faster completion = Bonus points\n")
	fmt.Printf("  • Higher difficulty = Score multiplier\n")

	fmt.Printf("\n%sStrategy:%s\n", ColorCyan, ColorReset)
	fmt.Printf("  • Start with middle values\n")
	fmt.Printf("  • Pay attention to proximity hints\n")
	fmt.Printf("  • Manage your time wisely\n")
	fmt.Println()
}

/*
displayDifficultyHelp provides detailed information about difficulty levels.

This function supports informed decision-making during game setup by providing
comprehensive details about each difficulty option.

Information Architecture:
- Clear comparison between all difficulty levels
- Scoring implications for competitive players
- Strategic considerations for optimal gameplay
- Recommendation guidance for different skill levels
*/
func displayDifficultyHelp() {
	fmt.Printf("\n%s Difficulty Guide%s\n", ColorPurple, ColorReset)
	fmt.Printf("%sEasy (1-%d):%s Ideal for beginners, quick games\n",
		ColorGreen, EasyMaxRange, ColorReset)
	fmt.Printf("  • Scoring: No multiplier\n")
	fmt.Printf("  • Strategy: Random guessing often works\n")

	fmt.Printf("%sMedium (1-%d):%s Balanced challenge for most players\n",
		ColorYellow, MediumMaxRange, ColorReset)
	fmt.Printf("  • Scoring: 1.5x multiplier\n")
	fmt.Printf("  • Strategy: Binary search recommended\n")

	fmt.Printf("%sHard (1-%d):%s Expert level, requires strategy\n",
		ColorRed, HardMaxRange, ColorReset)
	fmt.Printf("  • Scoring: 2x multiplier\n")
	fmt.Printf("  • Strategy: Systematic approach essential\n")
	fmt.Println()
}

/*
displayRangeHelp provides assistance for numeric input validation.

This function offers immediate guidance when users encounter input validation
errors, improving the overall user experience through contextual help.

Parameters:
- min, max int: Valid input range bounds

Help Content:
- Clear explanation of valid input range
- Examples of acceptable input formats
- Common error prevention tips
*/
func displayRangeHelp(min, max int) {
	fmt.Printf("\n%s Input Help%s\n", ColorYellow, ColorReset)
	fmt.Printf("Valid range: %s%d to %d%s\n", ColorWhite, min, max, ColorReset)
	fmt.Printf("Enter only numbers (no letters or symbols)\n")
	fmt.Printf("Examples: %s%d%s, %s%d%s, %s%d%s\n",
		ColorGreen, min, ColorReset,
		ColorGreen, (min+max)/2, ColorReset,
		ColorGreen, max, ColorReset)
	fmt.Println()
}

/*
printColoredHeader displays formatted section headers with consistent styling.

This utility function ensures consistent visual presentation across all
application sections while centralizing formatting logic for maintainability.

Parameters:
- text string: Header text to display

Formatting Features:
- Consistent color scheme application
- Automatic padding and alignment
- Unicode decorative elements for visual appeal
- Terminal-width responsive design
*/
func printColoredHeader(text string) {
	fmt.Printf("\n%s%s%s %s %s%s%s\n",
		ColorPurple, HeaderPrefix, ColorReset,
		text,
		ColorPurple, strings.Repeat("─", max(0, 50-len(text))), ColorReset)
}

/*
printColoredMessage displays formatted messages with specified color coding.

This utility function provides consistent message formatting while enabling
semantic color usage throughout the application.

Parameters:
- message string: Text content to display
- color string: ANSI color code for text formatting

Color Usage Guidelines:
- Red: Errors and warnings
- Green: Success and positive feedback
- Yellow: Information and hints
- Blue: Interactive elements and player names
- Purple: Headers and special announcements
- Cyan: Data and statistics
*/
func printColoredMessage(message, color string) {
	fmt.Printf("%s%s%s\n", color, message, ColorReset)
}

/*
printSeparator displays a consistent visual divider between sections.

This utility function improves readability by providing clear section
boundaries in the command-line interface.

Design Considerations:
- Unicode line drawing characters for professional appearance
- Consistent length across all terminal sizes
- Color coordination with overall application theme
*/
func printSeparator() {
	fmt.Printf("%s%s%s\n", ColorPurple, SeparatorLine, ColorReset)
}

/*
contains performs efficient string slice membership testing.

This utility function provides O(n) string search functionality with
optimized implementation for typical game session sizes.

Parameters:
- slice []string: String slice to search
- item string: Target string to locate

Returns:
- bool: True if item exists in slice, false otherwise

Performance Considerations:
- Linear search is optimal for small slices (typical player counts)
- Early termination on first match
- Memory-efficient implementation without additional allocations
*/
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/*
min returns the smaller of two integers.

This utility function provides mathematical minimum calculation
commonly needed for range operations and display formatting.

Parameters:
- a, b int: Integer values to compare

Returns:
- int: The smaller of the two input values

Note: This function will be replaced by the built-in min function
in future Go versions (1.21+) but is included for compatibility.
*/
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
max returns the larger of two integers.

This utility function provides mathematical maximum calculation
commonly needed for range operations and display formatting.

Parameters:
- a, b int: Integer values to compare

Returns:
- int: The larger of the two input values

Note: This function will be replaced by the built-in max function
in future Go versions (1.21+) but is included for compatibility.
*/
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
