# Ultimate Number Guessing Game  

A feature-rich multiplayer number guessing game built in Go with an enhanced command-line interface.  

---

## **Features**  

### **Game Modes & Setup**  
- **Multiplayer Support**: Play with up to 10 players  
- **Difficulty Levels**:  
  - **Easy**: Numbers 1-50 (beginner-friendly)  
  - **Medium**: Numbers 1-100 (balanced challenge)  
  - **Hard**: Numbers 1-200 (expert level)  
- **Smart Scoring System**: Points based on attempts, time, and difficulty  
- **Time Limits**: 10-second limit per guess to keep games fast-paced  

### **User Experience**  
- **Colored Terminal Output**: Enhanced visuals with ANSI colors  
- **Interactive Help System**: Get guidance anytime during gameplay  
- **Smart Hints**: Proximity feedback ("Very close!", "Getting warmer...")  
- **Range Validation**: Helpful error messages for invalid inputs  

### **Statistics & Analytics**  
- **Game Statistics**: Track performance across multiple sessions  
- **All-Time Leaderboard**: Compare scores with other players  
- **Performance Metrics**: Average attempts, duration, and win rates  
- **Game History**: Detailed records of past matches  

---

## **How to Play**  

1. **Start the Game**: Run the program and follow the setup prompts.  
2. **Choose Difficulty**: Select **Easy**, **Medium**, or **Hard**.  
3. **Register Players**: Enter names or use auto-generated ones.  
4. **Take Turns**: Each player guesses the secret number within 10 seconds.  
5. **Win the Game**: The first correct guess wins, with points calculated based on performance.  
6. **View Results**: Check the leaderboard and session statistics.  

---

## **Installation & Running**  

```bash
# Clone or download the code
# Create a folder and run:
go mod init yourname/yourfolder
go run .  # or go run ./file.go
```

---

## **Gameplay Commands**  

- Enter any number within the selected difficulty range.  
- Type `help` during any input for assistance.  
- Invalid inputs and timeouts are handled automatically.  

---

## **Scoring System**  

| Factor               | Points Adjustment         |  
|----------------------|---------------------------|  
| **Base Score**       | 1000 points               |  
| **Attempt Penalty**  | -10 points per guess      |  
| **Time Penalty**     | -1 point per 5 seconds    |  
| **Difficulty Bonus** | Medium: 1.5x, Hard: 2x    |  

---

## **Example Game Session**  

```
Ultimate Number Guessing Game - Enhanced Edition  
Features: Multiplayer • Difficulty Levels • Smart Scoring • Statistics • Help System  

Difficulty Selection  
Available Difficulties:  
  1. Easy   - Range: 1-50 (Beginner friendly)  
  2. Medium - Range: 1-100 (Balanced challenge)  
  3. Hard   - Range: 1-200 (Expert level)  

Player Registration  
Enter number of players (1-10): 2  
Enter name for Player 1: Alice  
Enter name for Player 2: Bob  

Game Session Initialized  
Difficulty: Medium (Range: 1-100)  
Players: Alice, Bob  
Time Limit: 10s per guess  

[Alice's Turn] Enter your guess (1-100): 50  
Too high! Getting warmer...  

[Bob's Turn] Enter your guess (1-100): 25  
Too low! Close!  

Alice wins with 3 attempts in 45s!  
```

---

## **Technical Details**  

- **Language**: Go (Golang)  
- **Architecture**: Modular design with clear separation of concerns  
- **Concurrency**: Goroutine-based timeout handling  
- **Data Structures**: Efficient maps and slices for game state  
- **Error Handling**: Comprehensive validation and recovery  
- **Memory Management**: Optimized for typical game sessions  

---

## **Requirements**  

- **Go 1.16+**  
- **Terminal with ANSI color support** (most modern terminals)  
- **No external dependencies**  

---

**Happy guessing! May your numbers be ever accurate!** 🎯
