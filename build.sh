#!/bin/bash

# 设置项目名称
PROJECT_NAME="sshport"

# 创建输出目录
OUTPUT_DIR="dist"
mkdir -p $OUTPUT_DIR

# 获取当前版本（可以从git tag获取，或者手动设置）
VERSION=$(git describe --tags --always 2>/dev/null || echo "v1.0.0")

echo "🚀 开始编译 $PROJECT_NAME $VERSION"
echo "📁 输出目录: $OUTPUT_DIR"

# 编译 Linux AMD64
echo "🔨 编译 Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $OUTPUT_DIR/${PROJECT_NAME}-linux-amd64 .
if [ $? -eq 0 ]; then
    echo "✅ Linux AMD64 编译成功"
else
    echo "❌ Linux AMD64 编译失败"
    exit 1
fi

# 编译 Linux ARM64
echo "🔨 编译 Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $OUTPUT_DIR/${PROJECT_NAME}-linux-arm64 .
if [ $? -eq 0 ]; then
    echo "✅ Linux ARM64 编译成功"
else
    echo "❌ Linux ARM64 编译失败"
    exit 1
fi

# 编译 macOS AMD64 (Intel)
echo "🔨 编译 macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $OUTPUT_DIR/${PROJECT_NAME}-darwin-amd64 .
if [ $? -eq 0 ]; then
    echo "✅ macOS AMD64 编译成功"
else
    echo "❌ macOS AMD64 编译失败"
    exit 1
fi

# 编译 macOS ARM64 (Apple Silicon)
echo "🔨 编译 macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $OUTPUT_DIR/${PROJECT_NAME}-darwin-arm64 .
if [ $? -eq 0 ]; then
    echo "✅ macOS ARM64 编译成功"
else
    echo "❌ macOS ARM64 编译失败"
    exit 1
fi

# 显示编译结果
echo ""
echo "🎉 编译完成！"
echo "📦 编译产物:"
ls -lh $OUTPUT_DIR/

# 创建压缩包（可选）
echo ""
echo "📦 创建压缩包..."
cd $OUTPUT_DIR
for file in ${PROJECT_NAME}-*; do
    if [ -f "$file" ]; then
        tar -czf "${file}.tar.gz" "$file"
        echo "✅ 创建 ${file}.tar.gz"
    fi
done
cd ..

echo ""
echo "🚀 所有任务完成！"
echo "📁 文件位置: $OUTPUT_DIR/"
echo ""
echo "使用方法:"
echo "  Linux AMD64: ./$OUTPUT_DIR/${PROJECT_NAME}-linux-amd64"
echo "  Linux ARM64: ./$OUTPUT_DIR/${PROJECT_NAME}-linux-arm64"
echo "  macOS AMD64: ./$OUTPUT_DIR/${PROJECT_NAME}-darwin-amd64"
echo "  macOS ARM64: ./$OUTPUT_DIR/${PROJECT_NAME}-darwin-arm64"
