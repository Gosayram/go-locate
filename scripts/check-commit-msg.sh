#!/bin/bash

# Script to validate commit message format
# Supports two formats:
# 1. [TYPE] - description (custom format)
# 2. type: description (conventional commits format)

COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Define allowed types for custom format [TYPE] - description
CUSTOM_TYPES="ADD|CI|FEATURE|BUGFIX|FIX|INIT|DOCS|TEST|REFACTOR|STYLE|CHORE"

# Define allowed types for conventional format type: description
CONVENTIONAL_TYPES="feat|fix|docs|style|refactor|test|chore|ci|build|perf|revert|add|feature|bugfix|init"

# Check if commit message follows either pattern
if echo "$COMMIT_MSG" | grep -qE "^\[($CUSTOM_TYPES)\] - .+"; then
    echo "✅ Commit message format is valid (custom format)"
    exit 0
elif echo "$COMMIT_MSG" | grep -qE "^($CONVENTIONAL_TYPES)(\(.+\))?: .+"; then
    echo "✅ Commit message format is valid (conventional format)"
    exit 0
else
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Your commit message:"
    echo "  $COMMIT_MSG"
    echo ""
    echo "Supported formats:"
    echo ""
    echo "Format 1 - Custom format:"
    echo "  [TYPE] - description"
    echo ""
    echo "Format 2 - Conventional Commits:"
    echo "  type: description"
    echo "  type(scope): description"
    echo ""
    echo "Custom format types:"
    echo "  ADD, CI, FEATURE, BUGFIX, FIX, INIT, DOCS, TEST, REFACTOR, STYLE, CHORE"
    echo ""
    echo "Conventional format types:"
    echo "  feat, fix, docs, style, refactor, test, chore, ci, build, perf, revert"
    echo ""
    echo "Examples:"
    echo "  [ADD] - new search functionality"
    echo "  [FIX] - resolve configuration parsing error"
    echo "  feat: add new search functionality"
    echo "  fix(config): resolve parsing error"
    echo "  chore: update version to 0.1.1"
    echo ""
    exit 1
fi
