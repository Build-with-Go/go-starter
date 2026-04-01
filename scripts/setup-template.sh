#!/bin/bash

# Go Starter Template Setup Script
# This script helps customize the template for a new project

set -e

echo "🚀 Setting up Go Starter template for your project..."

# Get project information
echo "Please enter your project information:"
read -p "Project name (e.g., my-awesome-app): " PROJECT_NAME
read -p "GitHub username: " GITHUB_USERNAME
read -p "Project description: " PROJECT_DESCRIPTION
read -p "Author name: " AUTHOR_NAME

# Validate inputs
if [[ -z "$PROJECT_NAME" || -z "$GITHUB_USERNAME" ]]; then
    echo "❌ Project name and GitHub username are required!"
    exit 1
fi

echo ""
echo "🔧 Customizing project..."

# Update go.mod
echo "📝 Updating go.mod..."
sed -i.bak "s|github.com/Build-with-Go/go-starter|github.com/${GITHUB_USERNAME}/${PROJECT_NAME}|g" go.mod
rm go.mod.bak

# Update all Go files
echo "📝 Updating Go files..."
find . -name "*.go" -type f -exec sed -i.bak "s|github.com/Build-with-Go/go-starter|github.com/${GITHUB_USERNAME}/${PROJECT_NAME}|g" {} +
find . -name "*.go.bak" -delete

# Update README.md
echo "📝 Updating README.md..."
sed -i.bak "s|Go Starter|${PROJECT_NAME}|g" README.md
sed -i.bak "s|A production-ready, batteries-included Go project template for the Build-with-Go organization|${PROJECT_DESCRIPTION}|g" README.md
sed -i.bak "s|Build-with-Go/go-starter|${GITHUB_USERNAME}/${PROJECT_NAME}|g" README.md
rm README.md.bak

# Update Dockerfile
echo "📝 Updating Dockerfile..."
sed -i.bak "s|Build-with-Go/go-starter|${GITHUB_USERNAME}/${PROJECT_NAME}|g" Dockerfile
rm Dockerfile.bak

# Update GitHub Actions workflows
echo "📝 Updating GitHub Actions..."
find .github/workflows -name "*.yml" -type f -exec sed -i.bak "s|Build-with-Go/go-starter|${GITHUB_USERNAME}/${PROJECT_NAME}|g" {} +
find .github/workflows -name "*.yml.bak" -delete

# Update configuration
echo "📝 Updating configuration..."
cp configs/config.example.yaml configs/config.yaml

# Create initial commit
echo "📝 Creating initial commit..."
git add .
git commit -m "Initial commit: ${PROJECT_NAME} setup from Go Starter template"

echo ""
echo "✅ Template setup completed!"
echo ""
echo "📋 Next steps:"
echo "1. Review the changes: git diff HEAD~1"
echo "2. Update configs/config.yaml with your settings"
echo "3. Add your business logic in internal/ directories"
echo "4. Update README.md with more project details"
echo "5. Run 'make deps && make run' to test your new project"
echo ""
echo "🚀 Your project is ready to develop!"
echo ""
echo "📊 Project info:"
echo "   Name: ${PROJECT_NAME}"
echo "   GitHub: https://github.com/${GITHUB_USERNAME}/${PROJECT_NAME}"
echo "   Module: github.com/${GITHUB_USERNAME}/${PROJECT_NAME}"
