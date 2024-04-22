package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// {"Size":5003804672,"Head":"volume-head-000.img","Dirty":true,"Rebuilding":false,"Error":"","Parent":"","SectorSize":512,"BackingFilePath":""}
type VolumeMeta struct {
	Size            int64
	Head            string
	Dirty           bool
	Rebuilding      bool
	Error           string
	Parent          string
	SectorSize      int
	BackingFilePath string
}

func main() {
	rootPath := os.Getenv("FIND_PATH")
	printAll := os.Getenv("PRINT_ALL")
	if rootPath == "" {
		log.Fatalln("empty rootPath!!,set it into the env FIND_PATH")
	}
	if printAll == "" {
		printAll = "1"
	}
	// rootPath := "/home/runner/work/lanactions/lanactions/k8s.io/registry.k8s.io/images/tests/"
	// TODO 拼凑docker命令
	//  docker run -v /dev:/host/dev -v /proc:/host/proc -v /var/lib/longhorn/replicas/pvc-c71e014b-fe06-463a-82f7-c9da74c55cdb-9339ab17:/volume --privileged longhornio/longhorn-engine:v1.6.1 launch-simple-longhorn pvc-c71e014b-fe06-463a-82f7-c9da74c55cdb 10737418240
	fss, err := os.ReadDir(rootPath)
	checkErr(err)
	for _, f := range fss {
		if f.IsDir() {
			// log.Println("fs.path:", f.Name())
			data, err := os.ReadFile(rootPath + f.Name() + "/volume.meta")
			checkErr(err)
			v := &VolumeMeta{}
			err = json.Unmarshal(data, v)
			checkErr(err)
			if printAll == "1" {
				log.Printf("dirty:%t,rebuilding:%t \n", v.Dirty, v.Rebuilding)
			} else {
				if !v.Dirty && !v.Rebuilding {
					log.Println("fs.path is not dirty and rebuilding:", f.Name())
					splitStrs := strings.Split(f.Name(), "-")
					replicasName := splitStrs[len(splitStrs)-1]
					pvcName := strings.ReplaceAll(f.Name(), "-"+replicasName, "")
					dockerCommand := fmt.Sprintf("docker run -v /dev:/host/dev -v /proc:/host/proc -v %s%s:/volume --privileged longhornio/longhorn-engine:v1.6.1 launch-simple-longhorn %s %d", rootPath, f.Name(), pvcName, v.Size)
					log.Println("got the command:\n", dockerCommand)
				}
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
