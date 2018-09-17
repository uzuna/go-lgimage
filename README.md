# golang image test

nodejs で画像生成サーバーを作っていたが、canvasのbuild周りの取り扱いが難しいことと
書きやすさの関係でgolangの使用を始めたので、golangで置き換えが可能か検証を行う


## referece

- [ease function](https://gist.github.com/gre/1650294)


http://foxcodex.html.xdomain.jp/index.html


## task

ヒストグラム表示の追加

min -> maxにむけったのでテーブル情報
```
```


coloscale 20分割なのでhistgramをn倍の粒度で入れる?

そもそもbinの幅が決まらない
https://note.mu/utaka233/n/n797c1d92ec78

Scottの公式で正規分布に従って得られるなら幅が3.49*sd/math.Cbrt(n)

Newrelicのように大きなテールが出る場合はhistgrmは95percentileで切って作ることもある
https://newrelic.degica.com/docs/data-analysis/user-interface-functions/view-your-data/histograms-viewing-data-distribution

- エイリアス周りの処理ピクセル小数点以下の取り扱いを考える
    chartとしては辺に混ざると濁るので整数値に切りなおすのがベター