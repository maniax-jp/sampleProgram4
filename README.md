# Go言語ブロック崩しゲーム

Go言語で実装されたブロック崩し（Breakout）ゲームです。ebitenライブラリを使用して2Dグラフィックスを描画します。

## 機能

### 基本機能
- **ボール**: 物理演算による反射動作
- **パドル**: 左右キーによる操作
- **ブロック**: 6行12列の色分けされたブロック配置
- **スコア管理**: ブロック破壊時のスコア加算
- **ゲーム状態**: ゲームオーバー/勝利条件の判定

### 対応プラットフォーム
- **デスクトップ**: Windows、Linux、macOS対応
- **ブラウザ**: WebAssembly（WASM）でブラウザ実行対応
- **モバイル**: タッチ操作対応（WASM版）

### ゲームルール
- 左右矢印キーでパドルを操作
- ボールをパドルで跳ね返し、すべてのブロックを破壊する
- ボールが底面に落ちるとゲームオーバー
- すべてのブロックを破壊すると勝利
- スペースキーでゲーム再開

## ビルド方法

### デスクトップ版

#### 必要な依存関係（Linux）
```bash
sudo apt-get update
sudo apt-get install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

> **注意**: サーバー環境やCI環境では、デスクトップ版のビルドに必要なX11ライブラリが利用できない場合があります。その場合はWASM版をご利用ください。

#### ビルド
```bash
# Makefileを使用
make build

# または直接ビルド
go build -o breakout main.go
```

### WASM版（ブラウザ対応）

#### WASM ビルド
```bash
# 完全なWASMビルド（HTML含む）
make wasm-full

# WASMファイルのみビルド
make wasm
```

#### ローカルテスト
```bash
# ローカルサーバーで実行（自動でブラウザ用サーバーを起動）
make serve
# http://localhost:8080 でアクセス可能
```

#### GitHub Pages デプロイ
1. リポジトリの Settings > Pages で GitHub Pages を有効化
2. Source を "GitHub Actions" に設定
3. main ブランチにプッシュすると自動デプロイ

## 実行方法

### デスクトップ版
```bash
# ビルド後に実行
./breakout

# または直接実行
make run
```

### WASM版
- GitHub Pages: `https://username.github.io/repository-name/`
- ローカル: `make serve` 後に `http://localhost:8080`

## 開発環境

- Go 1.24.4以上
- ebiten v2.8.8
- グラフィック環境（X11対応）

## ゲーム画面

- 画面サイズ: 640x480ピクセル
- ブロック配置: 6行12列（全72個）
- ブロック色: 行ごとに異なる色（赤、オレンジ、黄、緑、青、紫）

## 技術詳細

### 主要構造体
- `Game`: ゲーム全体の状態管理
- `Ball`: ボールの位置と速度
- `Paddle`: パドルの位置とサイズ
- `Block`: ブロックの位置、サイズ、可視性

### 衝突検出
矩形同士の衝突検出アルゴリズムを使用してボールとパドル/ブロックの衝突を判定

### 物理演算
- ボールの反射角度計算
- パドル位置に応じた反射角度調整
- 壁との衝突処理
