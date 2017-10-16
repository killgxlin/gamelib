package discover

import (
	"gamelib/base/util"

	"github.com/coreos/etcd/clientv3"
)

var (
	Client    *clientv3.Client
	Namespace string
)

func Init(etcdaddr string) {
	c, e := clientv3.New(clientv3.Config{Endpoints: []string{etcdaddr}})
	util.PanicOnErr(e)

	Client = c
}
