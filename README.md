# LapTrace

エンデューロごっこのデータをもとに順位変動をグラフ化する

- https://sago35.github.io/laptrace/

画面は計測サーバー ([spago-examples](https://github.com/sago35/spago-examples)) の計測画面と同じ見た目で、
「順位表」「周回記録」「全員の順位変動」をボタンで切り替える。順位表と周回記録の行をクリック (タップ) すると
その人の順位変動 (折れ線と周回ごとの順位) が出る。スマホでも見られる。
観戦している人向けの画面の見かたは `help.html` (画面上部の「見かた」リンク) にある。

## 使い方

URL の `path` でどのデータを見るかを指定する。

```
https://sago35.github.io/laptrace/?path=<パス>
```

パスは `races/20251130_01` のように `races/` を含めた形で、静的ファイルも Firebase も同じパスを使う。

1. まず `<パス>.json` を静的ファイルとして探す (このページと同じ場所からの相対パス)。
   `?path=races/20251130_01` なら `races/20251130_01.json`。見つかればそれを表示して終わり
2. 無ければ Firebase Realtime Database の `<パス>` (`races/20251130_01`) を購読し、更新があるたびに
   描き直す (リアルタイム表示)。このときは「順位は仮のものです」の注意を出す

### レース当日 (リアルタイム)

計測サーバー (spago-examples) に Firebase への送信が入っている。計測画面の「設定」で
「LapTrace に送る」を ON にして送り先のパス (例 `races/20251130_01`) を入れると、周回を記録するたびに
下記の形のデータが Firebase のそのパスへ送られる (サーバー起動時に環境変数 `LAPTRACE_SERVICE_ACCOUNT`
でサービスアカウント鍵を渡す)。観客は `https://sago35.github.io/laptrace/?path=races/20251130_01` を開く。

同じことをするコマンド版が `cmd/putlapjson` (log.json を監視して送る。`serviceAccount.json` が要る):

```
$ go run ./cmd/putlapjson races/20251130_01 /path/to/log.json
```

### レース後 (静的ファイル)

Firebase に送ったデータは REST で取れる (読み取りは公開)。これを同じパス `<パス>.json`
(`races/20251130_01.json`) として置いてコミットすれば、以後は静的ファイルとして表示される。

```
$ curl -o races/20251130_01.json \
    https://laptrace-live01-default-rtdb.asia-southeast1.firebasedatabase.app/races/20251130_01.json
```

計測サーバーの log.json をそのまま置いてもよい (エントリー名は出ない)。

### ローカルで見る

リポジトリのルートで静的ファイルを配るサーバーを立てる。`path` は Firebase や Pages と同じ。

```
$ python -m http.server 9000
$ # ブラウザで http://localhost:9000/?path=races/20251130_01
```

`races/` に無いパスを指定すると Firebase に接続する (ネットが要る)。

## データの形

計測サーバーの log.json と同じ `Name` / `Start` / `Records` に、計測サーバーが `Entrants` を足したもの。

```json
{
    "Name": "2025/11/30 ごっこビギナー",
    "Start": { "No": 0, "Time": 1764466190919318500, "Lap": 0, "LapTime": 0, "Index": 0 },
    "Records": [
        { "No": 8, "Time": 1764466338423456500, "Lap": 1, "LapTime": 147504138000, "Index": 0, "Class": "ビギナー" },
        { "No": 44, "Time": 1764466345026685000, "Lap": 1, "LapTime": 154107366500, "Index": 1, "Class": "ビギナー" }
    ],
    "Entrants": [
        { "Class": "ビギナー", "No": 8, "Name": "名前", "Team": "チーム名", "Bike": "CRF250R" },
        { "Class": "ビギナー", "No": 44, "Name": "名前", "Team": "", "Bike": "" }
    ]
}
```

- `Start.Time` と `Records[].Time` は Unix 時刻のナノ秒、`LapTime` はナノ秒の経過時間
- `Finish` (Unix 時刻のナノ秒) があれば計測サーバーでレースを終了している。画面の更新時刻の横に「レース終了」と出す
- `Records[].Lap` はその人の周回数、`LapTime` は前の周回からの時間。無ければページ側で数え直す
- `Records[].Class` は NFC で記録したときのクラス名 (無くてもよい)
- `Entrants` はそのレースのクラスの人の `Class` / `No` / `Name` / `Team` / `Bike`。`Class` と `No` で
  レコードと結び付け、凡例・表・ツールチップに名前とクラスを出す。無いデータ (古い `races/*.json` など)
  では名前の列が空欄になるだけで動く
- 順位は周回数が多い順、同じなら最後の通過が早い順。クラス順位は同じクラス名の中での順位
