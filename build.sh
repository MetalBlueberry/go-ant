#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./docs/go-ant-ebiten.wasm   ./cmd/go-ant-ebiten/