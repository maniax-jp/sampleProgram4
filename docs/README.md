# Go Breakout Game - WebAssembly版

このディレクトリにはGo言語で作成されたブロック崩しゲームのWebAssembly版が含まれています。

## ファイル構成

- `index.html` - ゲームのメインページ
- `breakout.wasm` - Go言語からコンパイルされたWebAssemblyバイナリ
- `wasm_exec.js` - Go WebAssembly実行環境

## 遊び方

1. ブラウザで `index.html` を開いてください
2. ゲームが自動的に読み込まれます
3. 左右の矢印キーでパドルを操作します
4. すべてのブロックを破壊してください！

## GitHub Pages

このゲームはGitHub Pagesでプレイできます：
[https://maniax-jp.github.io/sampleProgram4/](https://maniax-jp.github.io/sampleProgram4/)

## 技術詳細

- Go 1.24.4
- Ebiten v2.8.8 (2Dゲームエンジン)
- WebAssembly (WASM)
- ブラウザ互換性: Chrome, Firefox, Safari, Edge (WebAssembly対応ブラウザ)