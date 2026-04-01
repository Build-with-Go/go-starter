#!/bin/bash
set -e

echo "🏗️  Building go-starter..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf bin/
mkdir -p bin/

# Build for multiple platforms
echo "📦 Building binaries..."
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"

for platform in $PLATFORMS; do
    GOOS=$(echo $platform | cut -d'/' -f1)
    GOARCH=$(echo $platform | cut -d'/' -f2)
    
    OUTPUT_NAME="server-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "  📦 Building for ${GOOS}/${GOARCH}..."
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-w -s" \
        -o bin/${OUTPUT_NAME} \
        ./cmd/server
    
    echo "  ✅ Built bin/${OUTPUT_NAME}"
done

# Create checksums
echo "🔐 Creating checksums..."
cd bin/
sha256sum * > checksums.txt
cd ..

echo "✅ Build completed successfully!"
echo "📁 Binaries available in: bin/"
echo "📋 Checksums: bin/checksums.txt"

# Show build info
echo ""
echo "📊 Build information:"
ls -la bin/
