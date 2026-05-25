# Insider Football League Simulation API

This project is a RESTful API built with Go that simulates a 4-team football league (based on Premier League rules). It handles match simulations, dynamic league standings, and implements an advanced probabilistic model to predict championship chances as the league progresses.

---

## 🚀 Architecture & Engineering Decisions

This project is designed with a strong emphasis on clean code, maintainability, and Go idioms, avoiding unnecessary third-party frameworks.

* **Layered Architecture (Clean Architecture):** The project is strictly divided into three layers: `Handler` (HTTP layer), `Service` (Business/Simulation logic), and `Repository` (Data access). 
* **Struct Composition & Interface-Based Design:** Since Go does not support classical OOP inheritance, the system relies heavily on Interface-based design and Dependency Injection. This decouples the modules, making the Service and Handler layers completely agnostic of the underlying database implementation.
* **Standard Library Mastery:** Routing and HTTP server operations are handled purely via `net/http`. JSON serialization/deserialization is seamlessly mapped to Go structs using the standard `encoding/json` package.

### Database Strategy: Dynamic Calculation vs. Snapshot
*PostgreSQL* is used via the standard `database/sql` package. 
For the league standings, I opted for a **Dynamic Calculation Approach** rather than keeping a separate `Standings` table or historical snapshots. The database only stores `Teams` and `Matches`. Standings, goal differences, and points are calculated dynamically in the Service layer based on the match history.
* *The "Why":* This decision was made specifically to support the **"Edit Match Results"** extra feature. If I had used a weekly snapshot approach, editing a past match would require complex cascading updates across all subsequent weeks. 
* *The Trade-off:* While a snapshot approach offers faster read queries for historical data, calculating a 4-team, 6-week league dynamically is an O(1) operation with negligible performance cost. 

### 🧠 The Simulation Engine
* **Match Score Generation (Poisson Distribution):** Match scores are generated using a Poisson distribution algorithm based on the respective `strength` metrics of the home and away teams.
* **Championship Prediction (Monte Carlo Simulation):** From the 4th week onwards, the system predicts championship probabilities. Instead of mutating the actual database, the current standings and unplayed matches are copied to memory (RAM). A **Monte Carlo simulation** runs 10,000 iterations of the remaining fixtures using the Poisson distribution. The algorithm tallies the champion of each iteration and calculates the exact percentage probability for each team, returning it instantly in the API response.

---

## 🛠 Tech Stack
* **Language:** Go (Golang) 1.22+
* **Database:** PostgreSQL 15
* **Containerization:** Docker & Docker Compose
* **Libraries:** `net/http`, `encoding/json`, `database/sql`, `github.com/lib/pq`

---

## ⚙️ Setup & Installation (Docker)

The project is fully containerized. You do not need to install Go or PostgreSQL on your local machine to run it. The database schema and initial mock data (Teams and 6-week fixtures) are automatically embedded into the Go binary using `//go:embed` and executed during startup.

1. Clone the repository and navigate to the project directory:
   ```bash
   git clone <your-repository-url>
   cd insider-league
