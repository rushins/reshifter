package restore

import (
	"context"
	"os"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
)

func TestRestore(t *testing.T) {
	port := "4001"
	// testing insecure etcd 2 and 3:
	tetcd := "http://localhost:" + port
	// backing up to remote https://play.minio.io:9000:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
	etcd2Restore(t, port, tetcd)
	etcd3Restore(t, port, tetcd)
	// testing secure etcd 2 and 3:
	// tetcd := "https://localhost:" + port
	// TBD
}

func etcd2Restore(t *testing.T, port, tetcd string) {
	target := types.DefaultWorkDir
	err := util.Etcd2Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd2 at %s: %s", tetcd, err)
		return
	}
	c2, err := util.NewClient2(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}
	kapi := client.NewKeysAPI(c2)
	err = util.SetKV2(kapi,
		types.KubernetesPrefix+"/namespaces/kube-system",
		"{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}",
	)
	if err != nil {
		t.Errorf("Can't create key %snamespaces/kube-system: %s", types.KubernetesPrefix, err)
		return
	}
	backupid, err := backup.Backup(tetcd, target, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}
	_ = util.EtcdDown()
	err = util.Etcd2Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd2 at %s: %s", tetcd, err)
		return
	}
	_, err = Restore(tetcd, backupid, target, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during restore: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(backupid + ".zip")
	_ = util.EtcdDown()
}

func etcd3Restore(t *testing.T, port, tetcd string) {
	target := types.DefaultWorkDir
	err := util.Etcd3Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd3 at %s: %s", tetcd, err)
		return
	}
	c3, err := util.NewClient3(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd3 at %s: %s", tetcd, err)
		return
	}
	_, err = c3.Put(context.Background(),
		types.KubernetesPrefix+"/namespaces/kube-system",
		"{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}",
	)
	if err != nil {
		t.Errorf("Can't create key %snamespaces/kube-system: %s", types.KubernetesPrefix, err)
		return
	}
	backupid, err := backup.Backup(tetcd, target, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}
	_ = util.EtcdDown()
	err = util.Etcd3Up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd3 at %s: %s", tetcd, err)
		return
	}
	_, err = Restore(tetcd, backupid, target, "play.minio.io:9000", "reshifter-test-cluster")
	if err != nil {
		t.Errorf("Error during restore: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(backupid + ".zip")
	_ = util.EtcdDown()
}
