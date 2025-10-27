# 🧠 Snip CLI Utilities

This folder contains **Python-based utilities** used for data preprocessing, experimentation, or model-related tasks that complement the main Go application.

The Python environment here is **fully isolated** and managed using [`uv`](https://github.com/astral-sh/uv), a modern and fast Python package manager and environment tool.

---

## 📦 Project Structure

cli/
├── pyproject.toml # Python project configuration
├── README.md # This file
├── .venv/ # Local virtual environment (auto-created)
├── .gitignore
├── convert_data.py # Example script (optional)
└── ...

## ⚙️ Setup Instructions

Follow these steps to set up and activate the Python environment.

### 1️⃣ Prerequisites

- Python 3.10 or newer
- `uv` installed (see [uv installation guide](https://docs.astral.sh/uv/getting-started/installation/))

To install `uv`:

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

### 2️⃣ Initialize the Environment

From the project root, initialize the Python CLI folder (only once):

```bash
uv init cli
```

Then navigate into it:

```bash
cd cli
```

Create a virtual environment:

```bash
uv venv
```

Activate it:

```bash
source .venv/bin/activate
```

You should now see (.venv) in your shell prompt.

### 3️⃣ Install Dependencies

To install dependencies listed in pyproject.toml:

```bash
uv sync
```

To add a new dependency:

```bash
uv add <package-name>
```

### 4️⃣ Deactivate the Environment

When done working inside the CLI environment:

```bash
deactivate
```

This switches back to your system Python.

### 5️⃣ Folder Access and Data Sharing

The CLI scripts can access project data stored outside this folder (for example, under ../data/).
To reference project paths from Python scripts, use relative paths:

```python
from pathlib import Path

ROOT_DIR = Path(**file**).resolve().parents[1]
DATA_DIR = ROOT_DIR / "data"

print(f"Data directory: {DATA_DIR}")
```

This ensures the scripts work correctly even if run from different locations.

### 💡 Tips

Never commit the .venv/ folder — it is ignored by default.

If you want to start fresh, simply delete .venv/ and recreate it using uv venv.

The pyproject.toml + uv.lock files define your full reproducible Python setup.

### 🧩 Example Workflow

1. Activate the environment:

```bash
cd cli
source .venv/bin/activate
```

2. Run a CLI script:

```bash
uv run semantic_search_cli.py verify
```

3. Deactivate when done:

```bash
deactivate
```
