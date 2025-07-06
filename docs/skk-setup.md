# SKK動作のための必要な設定

SKKを動かすためには以下の設定が必要です：

## 1. Denoのインストール

```bash
curl -fsSL https://deno.land/install.sh | sh
```

## 2. SKK辞書のダウンロード

```bash
mkdir -p ~/.skk
curl -o ~/.skk/SKK-JISYO.L https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L
```

## 必要なコンポーネント
- Deno: SKK実行環境
- SKK-JISYO.L: SKK辞書ファイル