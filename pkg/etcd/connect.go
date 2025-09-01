package etcd

//
//import (
//	"github.com/star-find-cloud/star-mall/conf"
//	log "github.com/star-find-cloud/star-mall/pkg/logger"
//	"go.etcd.io/etcd/clientv3"
//)
//
//var _etcdClient *clientv3.Client
//
//func init() {
//	conf := conf.GetConfig()
//	etcdConfig := clientv3.Config{
//		Endpoints:          conf.Etcd.Endpoints,
//		DialTimeout:        conf.Etcd.DialTimeout,
//		AutoSyncInterval:   conf.Etcd.AutoSyncInterval,
//		MaxCallSendMsgSize: conf.Etcd.MaxCallSendMsgSize,
//		MaxCallRecvMsgSize: conf.Etcd.MaxCallRecvMsgSize,
//	}
//	//if conf.Etcd.EnableTLS != false {
//	//	etcdConfig.TLS = conf.Etcd.TLS
//	//}
//
//	client, err := clientv3.New(etcdConfig)
//	if err != nil {
//		log.AppLogger.Errorln("Etcd error")
//		log.EtcdLogger.Errorf("Etcd connect faild: %s \n", err)
//	}
//	_etcdClient = client
//}
//
//func GetEtcdClient() *clientv3.Client {
//	return _etcdClient
//}
