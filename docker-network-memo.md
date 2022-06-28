# etcd お試し
docker run で起動した etcd に対して、 docker run で起動した etcdctl を使って PUT する。

- 0.0.0.0 で LISTEN してすべてのローカルIPに対して LISTEN を行うことで、どのインターフェースに対しても到達可能にする
``` shell
# etcd を起動
(/ º﹃º)/ < docker run -i -t --rm -p 2379:2379 --volume=etcd-data:/etcd-data --name etcd gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcd --name=etcd-1 --data-dir=/etcd-data --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

# コンテナ -> コンテナ通信
## IP を直接指定する
``` shell
# etcd コンテナの IP アドレスを確認
(/ º﹃º)/ < docker inspect etcd | grep IPAddress
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.2",
                    "IPAddress": "172.17.0.2",
# コンテナ間通信: コンテナの IP を指定して etcdctl を実行
(/ º﹃º)/ < docker run -i -t --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints 172.17.0.2:2379 put /chapter2/hello "Hello"
OK
```

## コンテナの名前(ホスト名)を指定する
- docker でデフォルトのブリッジネットワークに繋がれたコンテナは名前解決ができない

``` shell
(/ º﹃º)/ < docker run -i -t --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints etcd:2379 put /chapter2/hello "Hello"
{"level":"warn","ts":"2022-06-24T18:46:54.659Z","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-fe925c26-06f4-4a77-86d4-58dbb5269338/etcd:2379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = latest balancer error: all SubConns are in TransientFailure, latest connection error: connection error: desc = \"transport: Error while dialing dial tcp: lookup etcd on 192.168.65.5:53: no such host\""}
Error: context deadline exceeded
```

cf.
> デフォルトの bridge ネットワークでは自動的に名前解決が行われません。
> defaultネットワークに接続されたコンテナ同士はコンテナ名で通信をすることができません。
- [network コマンドを使う](https://man.plustar.jp/docker/engine/userguide/networking/work-with-networks.html)
- [Dockerコンテナのネットワーク周りについて \- 文鳥大好きエンジニアのガラクタ置き場](https://ponteru.hatenablog.com/entry/2019/05/04/185950)
- [Docker Compose入門 \(3\) ～ネットワークの理解を深める～ \| さくらのナレッジ](https://knowledge.sakura.ad.jp/23899/)
- [【docker network】Dockerコンテナのネットワークまとめ【コンテナ間通信・名前解決】 \| RARA Land](https://rara-world.com/docker-network/)

## コンテナの名前(ホスト名)を指定する ユーザー定義したネットワークを使用する
- ユーザー定義したネットワークにコンテナをつなげるだけで、接続先のコンテナ名を内部 DNS サーバ(127.0.0.11) が名前解決してくれる。

``` shell
(/ º﹃º)/ < docker network create test
0ba7af1a17a8a511849b3e2fbbb39264a99045ca8d2edb552880cf4e101f3838
(/ º﹃º)/ < docker network ls
NETWORK ID     NAME      DRIVER    SCOPE
d9ec819fb0a3   bridge    bridge    local
4d78df565126   host      host      local
e12598525c78   none      null      local
0ba7af1a17a8   test      bridge    local
(/ º﹃º)/ < docker network inspect test
[
    {
        "Name": "test",
        "Id": "0ba7af1a17a8a511849b3e2fbbb39264a99045ca8d2edb552880cf4e101f3838",
        "Created": "2022-06-24T19:26:51.920364574Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": {},
            "Config": [
                {
                    "Subnet": "172.18.0.0/16",
                    "Gateway": "172.18.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {},
        "Options": {},
        "Labels": {}
    }
]
```

etcd server を 作成したネットワークに接続して起動
``` shell
(/ º﹃º)/ < docker run -i -t --rm --network test -p 2379:2379 --volume=etcd-data:/etcd-data --name etcd gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcd --name=etcd-1 --data-dir=/etcd-data --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

クライアントを作成したネットワークに接続して起動
``` shell
(/ º﹃º)/ < docker run -i -t --network test --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints etcd:2379 put /chapter2/hello "Hello"
OK
```

cf.
- [Docker 小技・ノウハウ集 \- とほほのWWW入門](https://www.tohoho-web.com/docker/howto.html#name-resolution)


# ホスト->コンテナ通信
## localhost を指定する
- etcd コンテナは、 port forward でホストとつながっているので、ホストの localhost に対して通信を行えばつながる
``` shell
# etcdctl の --entdpoints の default は "http://127.0.0.1:2379"
(/ º﹃º)/ < etcdctl put /chapter2/hello 'Hello, World!'
OK
(/ º﹃º)/ < etcdctl --endpoints localhost:2379 put /chapter2/hello "Hello"
OK
(/ º﹃º)/ < etcdctl --endpoints 127.0.0.1:2379 put /chapter2/hello "Hello"
OK
```

## 直接コンテナの IP を指定する
- ホスト側から コンテナの IP を指定すると届かない。(別ネットワークなので届かないのはそれはそう)
``` shell
(/ º﹃º)/ < etcdctl --endpoints 172.17.0.2:2379 put /chapter2/hello "Hello"
{"level":"warn","ts":"2022-06-12T16:28:08.381+0900","logger":"etcd-client","caller":"v3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000356700/172.17.0.2:2379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = context deadline exceeded"}
Error: context deadline exceeded
```

## docker-compose で etcd お試し

``` shell
(/ º﹃º)/ < docker compose up -d
(/ º﹃º)/ < docker inspect -f "{{.NetworkSettings.Networks.NETWORK.IPAddress}}" etcd-lock_etcd1_1
<no value> # https://qiita.com/Mii4a_Shota/items/4ef7427a1118842769d6
(/ º﹃º)/ < docker compose down
```

``` shell
(/ º﹃º)/ < docker run -i -t --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints ...:2379 put /chapter2/hello "Hello"
```

- [\[Docker / Docker Compose\] コンテナのIPアドレスを固定する方法 \- zaki work log](https://zaki-hmkc.hatenablog.com/entry/2021/02/26/234357)

コンテナ間通信の問題なので、同一 ブリッジのネットワークにそれぞれのコンテナが存在すれば直接通信できるはず。
- [Dockerのコンテナ間通信をする方法をまとめる \- きり丸の技術日記](https://nainaistar.hatenablog.com/entry/2021/06/14/120000)


``` shell
(/ º﹃º)/ < docker-compose up -d

(/ º﹃º)/ < docker container ls
CONTAINER ID   IMAGE                                 COMMAND                  CREATED          STATUS          PORTS                               NAMES
6d50a82b8e03   gcr.io/etcd-development/etcd:v3.4.7   "/usr/local/bin/etcd…"   13 minutes ago   Up 13 minutes   0.0.0.0:2379->2379/tcp, 2380/tcp    etcd-lock_etcd1_1
4ef32f09ff71   gcr.io/etcd-development/etcd:v3.4.7   "/usr/local/bin/etcd…"   13 minutes ago   Up 13 minutes   2380/tcp, 0.0.0.0:22379->2379/tcp   etcd-lock_etcd3_1
1ed18230af26   gcr.io/etcd-development/etcd:v3.4.7   "/usr/local/bin/etcd…"   13 minutes ago   Up 13 minutes   2380/tcp, 0.0.0.0:12379->2379/tcp   etcd-lock_etcd2_1

# コンテナ間通信したいコンテナの内容を見る
(/ º﹃º)/ < docker inspect etcd-lock_etcd1_1

# そのコンテナが所属する仮想ブリッジの設定を見る
(/ º﹃º)/ < docker inspect etcd-lock_etcd1_1 | grep Network
            "NetworkMode": "etcd-lock_default",
        "NetworkSettings": {
            "Networks": {
                    "NetworkID": "b95660fd048d50ece05003bb0539b46bb4e1cf8b3189ab6fc394851034c4bc0d",

# 上記で調べた同一ネットワークを指定してコンテナ間通信を行う
(/ º﹃º)/ < docker run --network etcd-lock_default -i -t --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints http://6d50a82b8e03:2379 put /chapter2/hello "Hello"
OK
(/ º﹃º)/ < docker run --network etcd-lock_default -i -t --rm gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcdctl --endpoints http://etcd-lock_default:2379 put /chapter2/hello "Hello"
OK
```
