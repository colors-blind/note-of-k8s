## Note

### Step 1 build basic http server

```golang
package main

import (
        "fmt"
        "log"
        "net/http"
        "os"
)

func main() {
        http.HandleFunc("/", handler)
        log.Println("try to listend 8080 port...")
        log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request from %s", r.RemoteAddr)
        hostName, err := os.Hostname()
        if err != nil {
                log.Fatal("Hostname() error")
                os.Exit(1)
        }
        content := "Hello Kubernetes Beginners! Server in " + hostName
        fmt.Fprintln(w, content)
}

```

build with:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8s-for-beginners
```

### Step 2 create Dockerfile

```Dockerfile
FROM alpine:3.10
COPY k8s-for-beginners /
CMD ["/k8s-for-beginners"]
```

### Step 3 build image

```bash
docker build -t k8s-for-beginners:v0.0.1 .

Sending build context to Docker daemon  6.073MB
Step 1/3 : FROM alpine:3.10
3.10: Pulling from library/alpine
396c31837116: Pull complete 
Digest: sha256:451eee8bedcb2f029756dc3e9d73bab0e7943c1ac55cff3a4861c52a0fdd3e98
Status: Downloaded newer image for alpine:3.10
 ---> e7b300aee9f9
Step 2/3 : COPY k8s-for-beginners /
 ---> dfbab45b339b
Step 3/3 : CMD ["/k8s-for-beginners"]
 ---> Running in 21a685603ad3
Removing intermediate container 21a685603ad3
 ---> 6458bbc016a0
Successfully built 6458bbc016a0
Successfully tagged k8s-for-beginners:v0.0.1

```

### Step 4 run the container and test

```bash

[blue@master k8s-for-b]$ docker run -p 8080:8080 -d k8s-for-beginners:v0.0.1
76c93e0601ef277ae359335dcd2c00a9df44e40349d4e1f6cb69712b1d9f0c89
[blue@master k8s-for-b]$ curl http://localhost:8080
Hello Kubernetes Beginners! Server in 76c93e0601ef

[blue@master k8s-for-b]$ docker logs -f 76c93e0601ef
2021/10/02 02:44:21 try to listend 8080 port...
2021/10/02 02:45:42 Request from 172.17.0.1:34414
2021/10/02 02:45:45 Request from 172.17.0.1:34438
2021/10/02 02:45:46 Request from 172.17.0.1:34452
2021/10/02 02:47:37 Request from 172.17.0.1:35296
2021/10/02 02:48:22 Request from 192.168.0.19:52300


```

### Step 5 aboute namespace

```bash

[blue@master k8s-for-b]$ sudo ls -ls /proc/$(ps aux | grep k8s-for-beginners | grep -v color | awk {'print $2'})/ns
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 cgroup -> 'cgroup:[4026532746]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 ipc -> 'ipc:[4026532671]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 mnt -> 'mnt:[4026532669]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:44 net -> 'net:[4026532674]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 pid -> 'pid:[4026532672]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 pid_for_children -> 'pid:[4026532672]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 time -> 'time:[4026531834]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 time_for_children -> 'time:[4026531834]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 user -> 'user:[4026531837]'
0 lrwxrwxrwx 1 root root 0 10月  2 10:50 uts -> 'uts:[4026532670]'


/ # [blue@master k8s-for-b]$ sudo docker exec -it 76c93e0601ef  sh 
/ # ps aux
PID   USER     TIME  COMMAND
    1 root      0:00 /k8s-for-beginners
   34 root      0:00 sh
   41 root      0:00 ps aux

```

###  Step 6  Joining a Container to the Network Namespace of Another Container

get container id:

```bash
[blue@master k8s-for-b]$ docker ps| grep k8s-f | awk  {'print $1'}
76c93e0601ef
```

run another container with `--net` join k8s-for-beginners container network namespace

```
[blue@master k8s-for-b]$ docker run -it --net container:76c93e0601ef nicolaka/netshoot
                    dP            dP                           dP   
                    88            88                           88   
88d888b. .d8888b. d8888P .d8888b. 88d888b. .d8888b. .d8888b. d8888P 
88'  `88 88ooood8   88   Y8ooooo. 88'  `88 88'  `88 88'  `88   88   
88    88 88.  ...   88         88 88    88 88.  .88 88.  .88   88   
dP    dP `88888P'   dP   `88888P' dP    dP `88888P' `88888P'   dP   
                                                                    
Welcome to Netshoot! (github.com/nicolaka/netshoot)
                                                                


 76c93e0601ef  ~  curl localhost:8080
Hello Kubernetes Beginners! Server in 76c93e0601ef

 76c93e0601ef  ~  
```

check /proc/pid/ns network namespace

```bash

[blue@master ~]$ docker ps | grep net
104e0f44cf12   nicolaka/netshoot                                   "zsh"                    2 minutes ago       Up 2 minutes                                                   condescending_shockley
[blue@master ~]$ docker inspect --format '{{.State.Pid}}' 104e0f44cf12
135668

[blue@master ~]$ sudo ls -la  /proc/135668/ns
[sudo] blue 的密码：
总用量 0
dr-x--x--x 2 root root 0 10月  2 11:23 .
dr-xr-xr-x 9 root root 0 10月  2 11:20 ..
lrwxrwxrwx 1 root root 0 10月  2 11:23 cgroup -> 'cgroup:[4026532751]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 ipc -> 'ipc:[4026532749]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 mnt -> 'mnt:[4026532747]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 net -> 'net:[4026532674]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 pid -> 'pid:[4026532750]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 pid_for_children -> 'pid:[4026532750]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 time -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 time_for_children -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 10月  2 11:23 uts -> 'uts:[4026532748]'

[blue@master k8s-for-b]$  docker inspect --format '{{.State.Pid}}'  76c93e0601ef
98529
[blue@master k8s-for-b]$ sudo ls -la  /proc/98529/ns
[sudo] blue 的密码：
总用量 0
dr-x--x--x 2 root root 0 10月  2 11:24 .
dr-xr-xr-x 9 root root 0 10月  2 11:13 ..
lrwxrwxrwx 1 root root 0 10月  2 11:25 cgroup -> 'cgroup:[4026532746]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 ipc -> 'ipc:[4026532671]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 mnt -> 'mnt:[4026532669]'
lrwxrwxrwx 1 root root 0 10月  2 11:24 net -> 'net:[4026532674]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 pid -> 'pid:[4026532672]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 pid_for_children -> 'pid:[4026532672]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 time -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 time_for_children -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 10月  2 11:25 uts -> 'uts:[4026532670]'

```

the two containers net is the same:

```bash
 187 lrwxrwxrwx 1 root root 0 10月  2 11:25 user -> 'user:[4026531837]'
```

### Setp 7 look at cgroups

start a new container

```bash

$ sudo docker run -p 8080:8080 -d k8s-for-beginners:v0.0.1                                                                    
55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4

```

find container's memory limit

```bash

$ sudo find /sys/fs | grep 55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4 | grep mem
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cgroup.procs
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.use_hierarchy
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.tcp.usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.soft_limit_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.force_empty
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.pressure_level
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.move_charge_at_immigrate
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.tcp.max_usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.max_usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.oom_control
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.stat
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.slabinfo
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.limit_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.swappiness
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.numa_stat
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.failcnt
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.max_usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/tasks
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.failcnt
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cgroup.event_control
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.tcp.failcnt
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.limit_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/notify_on_release
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.usage_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.tcp.limit_in_bytes
/sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cgroup.clone_children
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.memory_pressure
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.memory_migrate
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.mem_exclusive
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.memory_spread_slab
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.effective_mems
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.mems
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.mem_hardwall
/sys/fs/cgroup/cpuset/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/cpuset.memory_spread_page
```

limit_in_bytes is the bytes of memory the container can use.

```bash
$ cat /sys/fs/cgroup/memory/docker/55048906a24e6c91d654596fab4c6da6c1bcca4edbe52d676569999571ce2ad4/memory.kmem.limit_in_bytes
9223372036854771712
```

The value  9223372036854771712 is the largest positive signed integer (2 63 – 1) in a 64-bit system,
which means unlimited memory can be used by this container.

### Step 8 add memory limit code

add the code and build as an image

```golang
package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	var longStrs []string
	times := 50
	for i := 1; i <= times; i++ {
		fmt.Printf("===========%d\n===========", i)
		longStrs = append(longStrs, buildString(1000000, byte(i)))
	}
	time.Sleep(3600)
}

func buildString(n int, b byte) string {
	var builder strings.Builder
	builder.Grow(n)
	for i := 0; i < n; i++ {
		builder.WriteByte(b)
	}
	return builder.String()
}

```

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o memory_limit &&  sudo docker build -t memorylimit:v0.0.1 .
```

run in local and run docker use `--memory 20m`

```bash

$ ./memory_limit 
===========1
======================2
======================3
======================4
======================5
======================6
======================7
======================8
======================9
======================10
======================11
======================12
======================13
======================14
======================15
======================16
======================17
======================18
======================19
======================20
======================21
======================22
======================23
======================24
======================25
======================26
======================27
======================28
======================29
======================30
======================31
======================32
======================33
======================34
======================35
======================36
======================37
======================38
======================39
======================40
======================41
======================42
======================43
======================44
======================45
======================46
======================47
======================48
======================49
======================50
===========%

$ sudo docker run --memory=20m --memory-swap=20m memorylimit:v0.0.1 
WARNING: Your kernel does not support swap limit capabilities or the cgroup is not mounted. Memory limited without swap.
===========1
======================2
======================3
======================4
======================5
======================6
======================7
======================8
======================9
======================10
======================11
======================12
======================13
======================14
======================15
======================16
======================17
======================18
======================19
===========%            

```

the container memory limit is 20MB, so it was killed.

use `sudo dmesg` check it out:

```
[31442.526170] Memory cgroup out of memory: Killed process 29475 (memory_limit) total-vm:703436kB, anon-rss:19348kB, file-rss:1024kB, shmem-rss:0kB, UID:0 pgtables:116kB oom_score_adj:0

```

