#!/usr/bin/env bash
set -euo pipefail

APP_DIR="/home/ubuntu/resume_ats"
BRANCH="main"

echo "Navigating to application directory..."
cd "$APP_DIR"

echo "Fetching latest code..."
git fetch origin

echo "Resetting working tree to origin/$BRANCH..."
git reset --hard "origin/$BRANCH"

echo "Cleaning untracked files..."
git clean -fd

echo "Building production binary..."
make build

echo "Restarting service..."
pm2 restart resume-ats

echo "Deployment completed successfully."