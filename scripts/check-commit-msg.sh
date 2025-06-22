#!/bin/bash

# Script to validate commit message format
# Expected format: [TYPE] - description
# Allowed types: ADD, CI, FEATURE, BUGFIX, FIX, INIT, DOCS, TEST, REFACTOR, STYLE, CHORE

COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Define allowed types
ALLOWED_TYPES="ADD|CI|FEATURE|BUGFIX|FIX|INIT|DOCS|TEST|REFACTOR|STYLE|CHORE"

# Check if commit message follows the pattern [TYPE] - description
if ! echo "$COMMIT_MSG" | grep -qE "^\[($ALLOWED_TYPES)\] - .+"; then
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Your commit message:"
    echo "  $COMMIT_MSG"
    echo ""
    echo "Expected format:"
    echo "  [TYPE] - description"
    echo ""
    echo "Allowed types:"
    echo "  ADD      - Adding new features or files"
    echo "  CI       - Continuous Integration changes"
    echo "  FEATURE  - New feature implementation"
    echo "  BUGFIX   - Bug fixes"
    echo "  FIX      - General fixes"
    echo "  INIT     - Initial project setup"
    echo "  DOCS     - Documentation changes"
    echo "  TEST     - Adding or modifying tests"
    echo "  REFACTOR - Code refactoring"
    echo "  STYLE    - Code style changes"
    echo "  CHORE    - Maintenance tasks"
    echo ""
    echo "Examples:"
    echo "  [ADD] - new search functionality"
    echo "  [FIX] - resolve configuration parsing error"
    echo "  [CI] - update GitHub Actions workflow"
    echo ""
    exit 1
fi

echo "✅ Commit message format is valid"
exit 0
