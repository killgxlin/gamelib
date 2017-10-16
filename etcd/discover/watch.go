package discover

import (
	"context"
	"log"
	"gamelib/base/util"
	"sort"
	"strconv"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// define ----------------------------------------------------------------------------
type WatchType string

const (
	List  WatchType = "list"
	Array WatchType = "array"
)

type WatchFun func(prefix string, values []string)

// indexValue ----------------------------------------------------------------------------
type indexValue struct {
	Index int
	Value string
}

func (v *indexValue) IsEmpty() bool {
	return v.Index == 0 && v.Value == ""
}

type byIndex []indexValue

func (a byIndex) Len() int           { return len(a) }
func (a byIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

// watcher ----------------------------------------------------------------------------
type Watcher struct {
	prefix string
	typ    WatchType

	values []indexValue

	ctx context.Context
	can context.CancelFunc
	wg  sync.WaitGroup
}

func (w *Watcher) getValues() []string {
	values := []string{}
	for _, v := range w.values {
		values = append(values, v.Value)
	}
	return values
}
func (w *Watcher) initValues(r *clientv3.GetResponse) {
	maxIndex := 0
	for _, kv := range r.Kvs {
		key := string(kv.Key)
		pl := len(w.prefix)
		sindex := key[pl:]
		nindex, e := strconv.Atoi(sindex)
		util.PanicOnErr(e)

		if nindex > maxIndex {
			maxIndex = nindex
		}
		w.values = append(w.values, indexValue{nindex, string(kv.Value)})
	}

	sort.Stable(byIndex(w.values))
	switch w.typ {
	case Array:
		rets := make([]indexValue, maxIndex, maxIndex)
		for _, v := range w.values {
			if rets[v.Index-1].IsEmpty() {
				rets[v.Index-1] = v
			}
		}
		w.values = rets
	}
}

func (w *Watcher) updateValues(r *clientv3.WatchResponse) {
	switch w.typ {
	case List:
		for _, ev := range r.Events {
			log.Println(ev)
			key := string(ev.Kv.Key)
			pl := len(w.prefix)
			index := key[pl:]
			nindex, e := strconv.Atoi(index)
			util.PanicOnErr(e)

			idx := -1
			for i, v := range w.values {
				if v.Index == nindex {
					idx = i
					break
				}
			}

			switch ev.Type {
			case mvccpb.PUT:
				if idx < 0 {
					w.values = append(w.values, indexValue{nindex, string(ev.Kv.Value)})
				}
			case mvccpb.DELETE:
				l := len(w.values)
				if idx >= 0 {
					for i := idx; i < l-1; i++ {
						w.values[i] = w.values[i+1]
					}
					w.values = w.values[0 : l-1]
				}
			}
		}
	case Array:
		for _, ev := range r.Events {
			log.Println(ev)
			key := string(ev.Kv.Key)
			pl := len(w.prefix)
			index := key[pl:]
			nindex, e := strconv.Atoi(index)
			util.PanicOnErr(e)

			nvalues := len(w.values)
			switch ev.Type {
			case mvccpb.PUT:
				needPut := nindex - nvalues
				for i := 0; i < needPut; i++ {
					w.values = append(w.values, indexValue{})
				}
				w.values[nindex-1] = indexValue{nindex, string(ev.Kv.Value)}
			case mvccpb.DELETE:
				w.values[nindex-1] = indexValue{}
			}
		}
	}
}

func NewWatch(prefix string, typ WatchType, f WatchFun) *Watcher {
	w := &Watcher{prefix: Namespace + prefix, typ: typ}
	w.ctx, w.can = context.WithCancel(context.TODO())
	rev := int64(0)
	{
		r, e := Client.Get(
			w.ctx,
			w.prefix,
			clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))
		util.PanicOnErr(e)
		//log.Println("watchget", r)
		w.initValues(r)
		rev = r.Header.Revision
		//log.Println(w.values)
		f(w.prefix, w.getValues())
	}
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		c := Client.Watch(w.ctx, w.prefix, clientv3.WithPrefix(), clientv3.WithRev(rev+1))
		for {
			select {
			case r := <-c:
				log.Println(r)
				w.updateValues(&r)
				f(w.prefix, w.getValues())
			case <-w.ctx.Done():
				return
			}
		}
	}()
	return w
}

func (w *Watcher) Close() {
	w.can()
	w.wg.Wait()
}
