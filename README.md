2022-06-25 Sat 04:53:59

# etcd で分散ロックを素振りする
- [etcdを使った分散ロック \- Carpe Diem](https://christina04.hatenablog.com/entry/etcd-distributed-lock)
- [Go言語で学ぶetcdプログラミング](https://zenn.dev/zoetro/books/560099d25d8d7f3c8449)

# etcd お試し
docker run で起動した etcd に対して、 docker run で起動した etcdctl を使って PUT する。

- 0.0.0.0 で LISTEN してすべてのローカルIPに対して LISTEN を行うことで、どのインターフェースに対しても到達可能にする
- port forward でホストとコンテナのポートをつなぐ

``` shell
# etcd を起動
(/ º﹃º)/ < docker run -i -t --rm -p 2379:2379 --volume=etcd-data:/etcd-data --name etcd gcr.io/etcd-development/etcd:v3.4.13 /usr/local/bin/etcd --name=etcd-1 --data-dir=/etcd-data --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

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
