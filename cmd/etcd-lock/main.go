package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 3 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	_, err = client.Put(context.TODO(), "/chapter3/option/key1", "my-value3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/option/key2", "my-value1")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/option/key3", "my-value2")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Get(context.TODO(), "/chapter3/option/",
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByValue, clientv3.SortAscend),
		clientv3.WithKeysOnly(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Count == 0 {
		log.Fatal("key /chapter3/key not found")
	}

	for _, kv := range resp.Kvs {
		log.Printf("%s: %s\n", kv.Key, kv.Value)
	}

	log.Printf(string(resp.Kvs[0].Value))

	_, err = client.Delete(context.TODO(), "/chapter3/key")
	if err != nil {
		log.Fatal(err)
	}
}

func printResponse(resp *clientv3.GetResponse) {
	fmt.Printf("header: %s\n", resp.Header.String())

	for i, kv := range resp.Kvs {
		fmt.Printf("kv[%d: %s\n]", i, kv.String())
	}
	fmt.Println()
}
