#!/bin/bash
set -e

echo "==> Generating protobuf code..."

PROTO_DIR=api/proto
OUT_DIR=gen/proto

rm -rf ${OUT_DIR}
mkdir -p ${OUT_DIR}

# 遍历所有 proto 文件并生成 Go 代码
for proto_file in $(find ${PROTO_DIR} -name '*.proto'); do
    protoc \
        --proto_path=${PROTO_DIR} \
        --go_out=${OUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=${OUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        ${proto_file}
done

echo "==> Generated code in ${OUT_DIR}/"
find ${OUT_DIR} -name '*.go'