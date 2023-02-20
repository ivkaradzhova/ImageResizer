package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

func getUserUidGid() (string, string) {
	user, err := user.Lookup("karadzhovai")
	if err != nil {
		return "false", "false"
	}
	return user.Uid, user.Gid
}

func chmodExample() {
	os.Chmod("/Users/karadzhovai/index.php", 0777)
}

//	func recursiveOwnership(rootDir string, mode fs.FileMode, uid int, gid int) error {
//		return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
//			if info.IsDir() {
//				fmt.Println(path, info)
//				os.Chmod(path, mode)
//				fmt.Println(os.Chown(path, uid, gid))
//			}
//
//			return nil
//		})
//
// }
var (
	likewiseService = "lwsmd"
	vmafdUsername   = "vmafdd-user"
)

func restartLikewiseServiceService() error {
	cmd := exec.Command("/usr/bin/service-control", "--restart", "lwsmd")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start services")
	}

	return nil
}

func updateVmafdOwnership() error {

	vmafdUser, err := user.Lookup(vmafdUsername)
	if err != nil {
		return err
	}
	// #nosec G302
	if err := os.Chmod("/var/lib/vmware/vmafdd_data", 0700); err != nil {
		return err
	}
	if err := os.Chmod("/var/lib/vmware/vmafdd_data/machine-ssl.crt", 0600); err != nil {
		return err
	}
	if err := os.Chmod("/var/lib/vmware/vmafdd_data/machine-ssl.key", 0600); err != nil {
		return err
	}

	uid, err := strconv.Atoi(vmafdUser.Uid)
	if err != nil {
		return err
	}
	gid, _ := strconv.Atoi(vmafdUser.Gid)
	if err != nil {
		return err
	}

	if err := os.Chown("/var/lib/vmware/vmafdd_data", uid, gid); err != nil {
		return err
	}
	if err := os.Chown("/var/lib/vmware/vmafdd_data/machine-ssl.crt", uid, gid); err != nil {
		return err
	}
	if err := os.Chown("/var/lib/vmware/vmafdd_data/machine-ssl.key", uid, gid); err != nil {
		return err
	}

	return recursiveOwnership("/etc/ssl/certs", fs.FileMode(0775), 0, gid)
}

func recursiveOwnership(rootDir string, mode fs.FileMode, uid int, gid int) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		if err := os.Chown(path, uid, gid); err != nil {
			return err
		}
		return nil
	})
}

func main() {
	restartLikewiseServiceService()

	updateVmafdOwnership()

	//uid, gid := getUserUidGid()
	//fmt.Println(uid, " ", gid)
	//
	//chmodExample()
	//// fmt.Println(os.Chown("/Users/karadzhovai/test_go_file_change/test_folder", 20, 502))
	//// recursiveOwnership("/Users/karadzhovai/test_go_file_change/test_folder", 0777, 0, 0)

	st := time.Now()
	elapsed := time.Since(st)
	b := time.Duration(11000000000)
	a := elapsed + 50000000000000
	fmt.Println(b)
	fmt.Println(a + b)
}
