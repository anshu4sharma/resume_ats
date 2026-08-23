#!/bin/bash
set -euo pipefail

USERNAME="anshu4sharma"
IMAGE_NAME="resume_ats"
REGISTRY="ghcr.io/$USERNAME/$IMAGE_NAME"
GIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "manual")
TAG_LATEST="$REGISTRY:latest"
TAG_SHA="$REGISTRY:$GIT_SHA"

echo "=== 🚀 Building and Pushing Docker Image to GHCR ==="

echo "🔨 Building ($TAG_SHA)..."
docker build \
  -t "$TAG_LATEST" \
  -t "$TAG_SHA" \
  -f Dockerfile .

echo "📤 Pushing to GHCR..."
docker push "$TAG_LATEST"
docker push "$TAG_SHA"

echo "✅ Pushed $TAG_LATEST and $TAG_SHA"