#!/usr/bin/env bash
set -e

cd /home/ubuntu/resume_ats

echo "Pulling latest code..."
git pull origin main

echo "Building..."
make build

echo "Restarting service..."
pm2 restart resume-ats

echo "Deploy complete"
