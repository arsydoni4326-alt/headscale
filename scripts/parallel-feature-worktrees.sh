#!/bin/bash
# Parallel Feature Worktrees Setup Script
# Creates Git worktrees for parallel feature development

set -e

# Configuration
REPO_DIR="$(pwd)"
PARENT_DIR="$(dirname "$REPO_DIR")"
REPO_NAME="$(basename "$REPO_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Usage
usage() {
    echo "Usage: $0 <feature1> <feature2> [feature3] [feature4] [feature5]"
    echo ""
    echo "Example:"
    echo "  $0 plugin-system multi-instance-dashboard monitoring-integrations"
    echo ""
    echo "This will create:"
    echo "  - Feature branches: feature/plugin-system, feature/multi-instance-dashboard, ..."
    echo "  - Worktrees: ../${REPO_NAME}-plugin-system, ../${REPO_NAME}-multi-instance-dashboard, ..."
    exit 1
}

# Check arguments
if [ $# -lt 2 ] || [ $# -gt 5 ]; then
    echo -e "${RED}Error: Provide 2-5 feature names${NC}"
    usage
fi

# Verify we're in a git repository
if [ ! -d .git ]; then
    echo -e "${RED}Error: Not in a git repository root${NC}"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes. Commit or stash them first.${NC}"
    git status --short
    exit 1
fi

# Record base commit
BASE_COMMIT=$(git rev-parse HEAD)
CURRENT_BRANCH=$(git branch --show-current)

echo -e "${GREEN}Base commit: $BASE_COMMIT${NC}"
echo -e "${GREEN}Current branch: $CURRENT_BRANCH${NC}"
echo ""

# Create branches and worktrees
for feature in "$@"; do
    BRANCH_NAME="feature/$feature"
    WORKTREE_PATH="${PARENT_DIR}/${REPO_NAME}-${feature}"
    
    echo -e "${YELLOW}Creating feature: $feature${NC}"
    
    # Check if branch already exists
    if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
        echo -e "${RED}  Branch $BRANCH_NAME already exists. Skipping.${NC}"
        continue
    fi
    
    # Check if worktree path already exists
    if [ -d "$WORKTREE_PATH" ]; then
        echo -e "${RED}  Directory $WORKTREE_PATH already exists. Skipping.${NC}"
        continue
    fi
    
    # Create branch
    git branch "$BRANCH_NAME"
    echo -e "${GREEN}  ✓ Created branch: $BRANCH_NAME${NC}"
    
    # Create worktree
    git worktree add "$WORKTREE_PATH" "$BRANCH_NAME" > /dev/null 2>&1
    echo -e "${GREEN}  ✓ Created worktree: $WORKTREE_PATH${NC}"
    echo ""
done

# Summary
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Worktree setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Worktrees created:"
git worktree list
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Navigate to each worktree: cd $PARENT_DIR/${REPO_NAME}-<feature>"
echo "2. Implement the feature"
echo "3. Test independently"
echo "4. Commit and push"
echo "5. Request merge approval"
echo ""
echo "See docs/parallel-development.md for detailed workflow."
