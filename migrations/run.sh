#!/bin/bash

# Get the current directory
SCRIPT_DIR_REL=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# Specify the name of the SQLite database file
DATABASE_FILE="database.db"
MIGRATIONS_FILE="schema.sql"

# Create file if it doesn't exist
touch "$DATABASE_FILE"

# Run the SQL commands to create tables
sqlite3 "$DATABASE_FILE" < "$SCRIPT_DIR_REL/$MIGRATIONS_FILE"
