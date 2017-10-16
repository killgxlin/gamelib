package discover

import (
	"context"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type Stub struct {
	lock   string
	key    string
	value  string
	sess   *concurrency.Session
	locker sync.Locker
}

func NewStub(lock string, key string, value string) (*Stub, error) {
	key = Namespace + key
	lock = Namespace + lock

	sess, e := concurrency.NewSession(Client, concurrency.WithTTL(3))
	if e != nil {
		sess.Close()
		return nil, e
	}

	locker := concurrency.NewLocker(sess, lock)

	locker.Lock()

	_, e = sess.Client().Put(context.TODO(), key, value, clientv3.WithLease(sess.Lease()))
	if e != nil {
		locker.Unlock()
		sess.Close()
		return nil, e
	}

	_, e = sess.Client().KeepAlive(context.TODO(), sess.Lease())
	if e != nil {
		locker.Unlock()
		sess.Close()
		return nil, e
	}

	return &Stub{lock: lock, key: key, value: value, sess: sess, locker: locker}, nil

}

func (s *Stub) Close() {
	s.locker.Unlock()
	s.sess.Close()
}
