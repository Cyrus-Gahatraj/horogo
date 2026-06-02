# 🌌 Horogo

Horogo calculates birth charts and provides AI-powered astrological readings. It consists of a Go CLI tool and a Python FastAPI server.

---

## Requirements
*   **Go** (1.20 or later)
*   **Python** (3.10 or later)
*   **Hugging Face Account & Token** (for AI predictions)

---

## Setup & Running Instructions

### 1. Configuration (All Operating Systems)
1. Copy `.env.example` to `.env`:
   *   **macOS/Linux:** `cp .env.example .env`
   *   **Windows:** `copy .env.example .env`
2. Open `.env` and fill in:
   *   `HF_TOKEN`: Your Hugging Face user access token.
   *   `PORT`: `8000`

---

### 2. Start the AI Server
Open a terminal in the project root to install dependencies and run the server.

#### Option A: Using `uv` (Recommended)
```bash
uv run main.py
```

#### Option B: Standard Python (macOS / Linux)
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -e .
python main.py
```

#### Option C: Standard Python (Windows)
```powershell
python -m venv .venv
.venv\Scripts\activate
pip install -e .
python main.py
```
*Note: Keep this terminal window open and running.*

---

### 3. Build and Run the CLI
Open a **new** terminal window in the project root to run the CLI.

#### Build the Executable:
*   **macOS / Linux:**
    ```bash
    go build -o horogo cmd/main.go
    ```
*   **Windows:**
    ```powershell
    go build -o horogo.exe cmd/main.go
    ```

#### Commands:
1.  **Generate a new profile chart:**
    *   *macOS/Linux:* `./horogo`
    *   *Windows:* `.\horogo.exe`
    *   *Saves calculations to `data/<Name>/raw/chart.json`.*

2.  **Interactive AI chat about a chart:**
    *   *macOS/Linux:* `./horogo ask <Name>`
    *   *Windows:* `.\horogo.exe ask <Name>`
    *   *Commands within chat:* type `help` for command info, `exit` to quit.

3.  **List existing profiles:**
    *   *macOS/Linux:* `./horogo ls`
    *   *Windows:* `.\horogo.exe ls`

---

## Profiles Layout
User profiles are saved in:
`data/<Name>/`
*   `raw/chart.json` — Calculated degrees, signs, houses, and nakshatras.
*   `index.md` — Main landing index for your personal astrology profile.
*   `wiki/` — Detailed analyses (personality, career, relationships, and house placements).
