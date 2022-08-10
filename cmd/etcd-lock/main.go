package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

	// _, err = client.Compact(context.TODO(), 188)
	// if err != nil {
	// 	log.Fatalf("Compact: %v", err)
	// }

	resp, err := client.Get(context.TODO(), "/chapter3/option",
		clientv3.WithPrefix(),
		// clientv3.WithSort(clientv3.SortByValue, clientv3.SortAscend),
		clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend),
		// clientv3.WithKeysOnly(),
	)

	if err != nil {
		log.Fatal(err)
	}

	if resp.Count == 0 {
		log.Fatal("key /chapter3/key not found")
	}

	printResponse(resp)

	// for _, kv := range resp.Kvs {
	// 	log.Printf("%s: %s\n", kv.Key, kv.Value)
	// }

	// log.Printf(string(resp.Kvs[0].Value))

	delResp, err := client.Delete(context.TODO(), "/chapter3/option/key3")
	if err != nil {
		log.Fatal(err)
	}
	printDelResponse(delResp)

	ch := client.Watch(context.TODO(), "/chapter3/watch/", clientv3.WithPrefix())

	for resp := range ch {
		if resp.Err() != nil {
			log.Fatal(resp.Err())
		}

		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				switch {
				case ev.IsCreate():
					fmt.Printf("CREATE: %q: %q\n", ev.Kv.Key, ev.Kv.Value)
				case ev.IsModify():
					fmt.Printf("MODIFY %q: %q\n", ev.Kv.Key, ev.Kv.Value)
				}
			case clientv3.EventTypeDelete:
				fmt.Printf("DELETE: %q: %q\n", ev.Kv.Key, ev.Kv.Value)
			}
		}
	}

	preventDropout(client)
}

func printResponse(resp *clientv3.GetResponse) {
	fmt.Printf("header: %s\n", resp.Header.String())

	for i, kv := range resp.Kvs {
		fmt.Printf("kv[%d: %s] %s\n", i, kv.String(), string(kv.Value))
	}
	fmt.Println()
}

func printDelResponse(resp *clientv3.DeleteResponse) {
	fmt.Printf("header: %s\n", resp.Header.String())
	fmt.Printf("%v", resp.Deleted)

	fmt.Println()

}

// 取りこぼしをふせぐ
// https://zenn.dev/zoetro/books/560099d25d8d7f3c8449/viewer/6ddfc44f1c97b6c949b8#%E5%8F%96%E3%82%8A%E3%81%93%E3%81%BC%E3%81%97%E3%82%92%E9%98%B2%E3%81%90
func preventDropout(client *clientv3.Client) {
	rev := nextRev()
	fmt.Printf("loaded revision: %d\n", rev)
	ch := client.Watch(context.TODO(), "/chapter3/watchFile", clientv3.WithRev(rev))
	for resp := range ch {
		if resp.Err() != nil {
			log.Fatal(resp.Err())
		}

		for _, ev := range resp.Events {
			fmt.Printf("[%d] %s %q : %q\n", ev.Kv.ModRevision, ev.Type, ev.Kv.Key, ev.Kv.Value)
			err := saveRev(ev.Kv.ModRevision)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("saved: %d\n", ev.Kv.ModRevision)
		}
	}
}

func nextRev() int64 {
	p := "./lastRevision"
	f, err := os.Open(p)
	if err != nil {
		os.Remove(p)
		return 0
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		os.Remove(p)
		return 0
	}

	rev, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		os.Remove(p)
		return 0
	}

	return rev + 1
}

func saveRev(rev int64) error {
	p := "./lastRevision"
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(strconv.FormatInt(rev, 10))
	return err
}
