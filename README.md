# LapTrace

エンデューロごっこのデータをもとに順位変動をグラフ化する

- https://sago35.github.io/laptrace/

画面は計測サーバー (spago-examples) の計測画面の順位表と同じ見た目で、行をクリック (タップ) すると
その人の順位変動 (折れ線と周回ごとの順位) が出る。「全員の順位変動」で全員分の折れ線 (ECharts) に
切り替えられる。

## データの形

計測サーバー ([spago-examples](https://github.com/sago35/spago-examples)) の log.json と同じ
`Name` / `Start` / `Records` に加えて、`Entrants` (そのレースのクラスの人の `Class` / `No` / `Name` /
`Team` / `Bike`) があれば、凡例・表・ツールチップにエントリー名とクラスを出す。
`Entrants` が無いデータ (古い `races/*.json` など) では番号だけを出す。
