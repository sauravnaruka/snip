#!/usr/bin/env python3
import subprocess
import sys

def main():
        
    # Forward all command-line arguments after this script to Go
    args = sys.argv[1:]
    if not args:
        print("Usage: python cli/keyword_search_cli.py <command> [args...]")
        sys.exit(1)

    # Construct command: go run . <args...>
    cmd = ["go", "run", ".", *args]

    # Run the Go command and stream its output
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as e:
        print(f"Command failed: {e}")
        sys.exit(e.returncode)

if __name__ == "__main__":
    main()
