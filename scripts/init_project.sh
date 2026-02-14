#!/bin/bash

# ==============================================================================
# PROJECT INITIALIZATION SCRIPT
# Role: Automated Refactoring & Setup
# Compliance: POSIX Standard
# ==============================================================================

set -e

# Warna untuk output log
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}[AUDIT] Starting Project Initialization Protocol...${NC}"

# 1. Validasi Keberadaan go.mod
if [ ! -f go.mod ]; then
    echo -e "${RED}[ERROR] go.mod not found! Please run this script from the project root.${NC}"
    exit 1
fi

# 2. Deteksi Module Saat Ini
CURRENT_MODULE=$(head -n 1 go.mod | awk '{print $2}')
echo -e "${YELLOW}[INFO] Current Module detected: ${CURRENT_MODULE}${NC}"

# 3. Input Module Baru
read -p "Enter new module name (e.g., github.com/mycompany/myproject): " NEW_MODULE

if [ -z "$NEW_MODULE" ]; then
    echo -e "${RED}[ERROR] Module name cannot be empty.${NC}"
    exit 1
fi

if [ "$CURRENT_MODULE" == "$NEW_MODULE" ]; then
    echo -e "${RED}[ERROR] New module name is identical to current one. No changes needed.${NC}"
    exit 0
fi

# 4. Konfirmasi Eksekusi
echo -e "${YELLOW}[WARN] This will replace all occurrences of '${CURRENT_MODULE}' with '${NEW_MODULE}'${NC}"
read -p "Are you sure? (y/n): " CONFIRM
if [[ "$CONFIRM" != "y" ]]; then
    echo -e "${RED}[ABORTED] Operation cancelled by user.${NC}"
    exit 0
fi

# 5. Eksekusi Find & Replace (Cross-platform compatible: Linux/MacOS)
echo -e "${GREEN}[EXEC] Refactoring codebase...${NC}"

# Deteksi OS untuk kompatibilitas sed
OS="$(uname)"
if [ "$OS" == "Darwin" ]; then
    # MacOS requires empty string for -i
    grep -rl "$CURRENT_MODULE" . --exclude-dir=.git --exclude=init_project.sh | xargs sed -i '' "s|$CURRENT_MODULE|$NEW_MODULE|g"
else
    # Linux standard
    grep -rl "$CURRENT_MODULE" . --exclude-dir=.git --exclude=init_project.sh | xargs sed -i "s|$CURRENT_MODULE|$NEW_MODULE|g"
fi

echo -e "${GREEN}[EXEC] Renaming complete.${NC}"

# 6. Reset Git History (Optional but Recommended for clean start)
read -p "Do you want to re-initialize Git (remove old history)? (y/n): " INIT_GIT
if [[ "$INIT_GIT" == "y" ]]; then
    echo -e "${GREEN}[EXEC] Re-initializing Git repository...${NC}"
    rm -rf .git
    git init
    git branch -m main
    echo -e "${GREEN}[INFO] Git initialized. New repository ready.${NC}"
fi

# 7. Go Mod Tidy
echo -e "${GREEN}[EXEC] Running go mod tidy...${NC}"
go mod tidy

echo -e "----------------------------------------------------------------"
echo -e "${GREEN}✅ PROJECT INITIALIZED SUCCESSFULLY${NC}"
echo -e "New Module: ${NEW_MODULE}"
echo -e "Next Step: Update your .env file and run 'make run'"
echo -e "----------------------------------------------------------------"
