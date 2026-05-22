package store_driver

//
//import (
//	"context"
//	"fmt"
//	"time"
//
//	"github.com/cruvie/kk-scheduler/go/internal/g_config"
//	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"google.golang.org/protobuf/encoding/protojson"
//)
//
///*
//[StoreEtcd]
//UserName = "root"
//Password = "root"
//Endpoints = ["http://host.docker.internal:2379"]
//
//*/
//const (
//	storeServiceKey = "kk-scheduler-service"
//	storeJobKey     = "kk-scheduler-job"
//)
//
//type StoreEtcd struct {
//	Client *clientv3.Client
//}
//
//func NewStoreEtcd() *StoreEtcd {
//	config := clientv3.Config{
//		Endpoints:   g_config.Config.StoreEtcd.Endpoints,
//		Username:    g_config.Config.StoreEtcd.UserName,
//		Password:    g_config.Config.StoreEtcd.Password,
//		DialTimeout: 2 * time.Second,
//	}
//	client, err := clientv3.New(config)
//	if err != nil {
//		panic(err)
//	}
//	return &StoreEtcd{
//		Client: client,
//	}
//}
//
//func (x *StoreEtcd) getJobKey(entry *kk_scheduler.PBJob) string {
//	return fmt.Sprintf("%s/%s/%s", storeJobKey, entry.GetServiceName(), entry.GetFuncName())
//}
//
//func (x *StoreEtcd) JobList(serviceName string) ([]*kk_scheduler.PBJob, error) {
//	resp, err := x.Client.Get(context.Background(), fmt.Sprintf("%s/%s", storeJobKey, serviceName), clientv3.WithPrefix())
//	if err != nil {
//		return nil, err
//	}
//
//	var jobs []*kk_scheduler.PBJob
//	for _, kv := range resp.Kvs {
//		var v kk_scheduler.PBJob
//		if err := protojson.Unmarshal(kv.Value, &v); err != nil {
//			return nil, err
//		}
//		jobs = append(jobs, &v)
//	}
//
//	return jobs, nil
//}
//
//func (x *StoreEtcd) JobGet(jobId string) (*kk_scheduler.PBJob, error) {
//	key := fmt.Sprintf("%s/%s/%s", storeJobKey, jobId)
//	resp, err := x.Client.Get(context.Background(), key)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(resp.Kvs) == 0 {
//		return nil, kk_scheduler.ErrJobNotFount
//	}
//
//	var entry kk_scheduler.PBJob
//	if err := protojson.Unmarshal(resp.Kvs[0].Value, &entry); err != nil {
//		return nil, err
//	}
//
//	return &entry, nil
//}
//
//func (x *StoreEtcd) JobPut(entry *kk_scheduler.PBJob) error {
//	key := x.getJobKey(entry)
//	value, err := protojson.Marshal(entry)
//	if err != nil {
//		return err
//	}
//
//	_, err = x.Client.Put(context.Background(), key, string(value))
//	return err
//}
//
//func (x *StoreEtcd) JobDelete(jobId string) error {
//	key := fmt.Sprintf("%s/%s/%s", storeJobKey, jobId)
//	_, err := x.Client.Delete(context.Background(), key)
//	return err
//}
//
//func (x *StoreEtcd) ServicePut(v *kk_scheduler.PBRegisterService) error {
//	key := fmt.Sprintf("%s/%s", storeServiceKey, v.GetServiceName())
//	value, err := protojson.Marshal(v)
//	if err != nil {
//		return err
//	}
//
//	_, err = x.Client.Put(context.Background(), key, string(value))
//	return err
//}
//
//func (x *StoreEtcd) ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error) {
//	key := fmt.Sprintf("%s/%s", storeServiceKey, serviceName)
//	resp, err := x.Client.Get(context.Background(), key)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(resp.Kvs) == 0 {
//		return nil, kk_scheduler.ErrServiceNotFount
//	}
//
//	var v kk_scheduler.PBRegisterService
//	err = protojson.Unmarshal(resp.Kvs[0].Value, &v)
//	return &v, err
//}
//
//func (x *StoreEtcd) ServiceList() ([]*kk_scheduler.PBRegisterService, error) {
//	key := storeServiceKey
//	resp, err := x.Client.Get(context.Background(), key, clientv3.WithPrefix())
//	if err != nil {
//		return nil, err
//	}
//
//	var services []*kk_scheduler.PBRegisterService
//	for _, kv := range resp.Kvs {
//		var v kk_scheduler.PBRegisterService
//		err := protojson.Unmarshal(kv.Value, &v)
//		if err != nil {
//			return nil, err
//		}
//		services = append(services, &v)
//	}
//	return services, nil
//}
//
//func (x *StoreEtcd) ServiceDelete(serviceName string) error {
//	if serviceName == "" {
//		return kk_scheduler.ErrServiceNameEmpty
//	}
//	key := fmt.Sprintf("%s/%s", storeServiceKey, serviceName)
//	_, err := x.Client.Delete(context.Background(), key)
//	return err
//}
