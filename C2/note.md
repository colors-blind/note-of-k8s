## Note

### Step 1 Running a Pod in Kubernetes

create file `k8s-for-beginners-pod.yaml`

```yaml

kind: Pod
apiVersion: v1
metadata:
  name: k8s-for-beginners
spec:
  containers:
  - name: k8s-for-beginners
    image: k8s-for-beginners:v0.0.1

```

start a pod

```bash
kubectl apply -f k8s-for-beginners-pod.yaml
```

check the pod

```bash
[blue@master k8s-for-b]$ kubectl get pod -o wide | grep k8s-for-be
k8s-for-beginners        1/1     Running   0             12m   10.244.104.19    node2   <none>           <none>
[blue@master k8s-for-b]$ curl http://10.244.104.19:8080
Hello Kubernetes Beginners! Server in k8s-for-beginners
```

delete the pod

```bash
[blue@master k8s-for-b]$ kubectl delete -f k8s-for-beginners-pod.yaml 
pod "k8s-for-beginners" deleted
```

### Step 2 service spceification

```
A service to abstract the network access to you application's pods.
Use Labels, which are defined in th pod definitions, and label selectors, which are defined in the 
Service definition, to describe this relationship.
```

#### Accessing a pod via a service

descript the pod via yaml

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: k8s-for-beginners
  labels:
    tier: frontend
spec:
  containers:
  - name: k8s-for-beginners
    image: k8s-for-beginners:v0.0.1

```

apply the conf
```bash
kubectl apply -f k8s-for-beginners-pod1.yaml
```

```bash
[blue@master k8s-for-b]$ kubectl get pod --show-labels
NAME                     READY   STATUS    RESTARTS      AGE    LABELS
k8s-for-beginners        1/1     Running   0             100s   tier=frontend
```

descript the service vai yaml

```yaml
kind: Service
apiVersion: v1
metadata:
  name: k8s-for-beginners
spec:
  selector:
    tier: frontend
  type: NodePort
  ports:
  - port: 80
    targetPort: 8080
```

apply the service conf
```bash
[blue@master k8s-for-b]$ kubectl  apply -f k8s-for-beginners-svc.yaml

```

get service 
```bash
[blue@master k8s-for-b]$ kubectl get service | grep beginners
k8s-for-beginners   NodePort       10.110.178.235   <none>        80:31710/TCP     3m40s

```
use 10.110.178.235:80 we can access the service

```bash
[blue@master k8s-for-b]$ curl http://10.110.178.235:80
Hello Kubernetes Beginners! Server in k8s-for-beginners

```

Port 31710 is exposed on every node, so we can access by `curl http://nodeip:31710`

test the service with curl

```bash
[blue@node2 ~]$ curl http://192.168.0.12:31710
Hello Kubernetes Beginners! Server in k8s-for-beginners
[blue@node2 ~]$ curl http://192.168.0.11:31710
Hello Kubernetes Beginners! Server in k8s-for-beginners
```

In the node 192.168.0.11 no docker k8s-for-beginners running, but can be accessed.


There are 3 layers of traffice transitions:

```
The first layer is from the external user to the machine IP at the auto-generated
random port (3XXXX).

The second layer is from the random port (3XXXX) to the Service IP (10.X.X.X) at
port  80 .

The third layer is from the Service IP (10.X.X.X) ultimately to the pod IP at
port  8080.
```

### Step 3 scaling a k8s applcation


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: k8sforbeginner
spec:
  replicas: 3
  selector:
    matchLebels:
      tier: simplehttp
  template:
    metadata:
      labels:
        tier: simplehttp
    spce:
      containers:
      - name: k8sforbeginner
        image: k8s-for-beginners:v0.0.1

```

apply a deploy

```bash
[blue@master k8s-for-b]$ kubectl apply -f k8s-for-beginners-deploy.yaml 
deployment.apps/k8sforbeginner created

```

show the deploy

```bash
[blue@master k8s-for-b]$ kubectl get deploy
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
k8sforbeginner   3/3     3            3           102s
nginx            1/1     1            1           43d
web              3/3     3            3           37d

```

scale more replicas

```bash

[blue@master k8s-for-b]$ kubectl scale deploy k8sforbeginner --replicas=5
deployment.apps/k8sforbeginner scaled
[blue@master k8s-for-b]$ kubectl get deploy
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
k8sforbeginner   5/5     5            5           11m
nginx            1/1     1            1           43d
web              3/3     3            3           37d

```

scale less replicas

```bash
[blue@master k8s-for-b]$ kubectl scale deploy k8sforbeginner --replicas=1
deployment.apps/k8sforbeginner scaled
[blue@master k8s-for-b]$ kubectl get deploy
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
k8sforbeginner   1/1     1            1           13m
nginx            1/1     1            1           43d
web              3/3     3            3           37d
[blue@master k8s-for-b]$ kubectl get pod
NAME                              READY   STATUS    RESTARTS        AGE
k8s-for-beginners                 1/1     Running   0               5h22m
k8sforbeginner-57d94db799-k427r   1/1     Running   0               3m33s
nginx-6799fc88d8-d6w82            1/1     Running   2 (6h29m ago)   35d
web-5bfc9bc56d-qd5g7              1/1     Running   3 (6h29m ago)   35d
web-5bfc9bc56d-vdpqz              1/1     Running   3 (6h29m ago)   35d
web-5bfc9bc56d-xllx8              1/1     Running   3 (6h29m ago)   35d

```

scale to 10 replicas

```bash
[blue@master k8s-for-b]$ kubectl scale deploy k8sforbeginner --replicas=10
deployment.apps/k8sforbeginner scaled
[blue@master k8s-for-b]$ kubectl get pod
NAME                              READY   STATUS              RESTARTS        AGE
k8s-for-beginners                 1/1     Running             0               5h28m
k8sforbeginner-57d94db799-4tm4k   0/1     ContainerCreating   0               1s
k8sforbeginner-57d94db799-579k5   0/1     ContainerCreating   0               1s
k8sforbeginner-57d94db799-8kslx   1/1     Running             0               4m43s
k8sforbeginner-57d94db799-bqvn8   0/1     Pending             0               1s
k8sforbeginner-57d94db799-dbsb5   0/1     ContainerCreating   0               2s
k8sforbeginner-57d94db799-dx5g7   0/1     ContainerCreating   0               1s
k8sforbeginner-57d94db799-k427r   1/1     Running             0               9m21s
k8sforbeginner-57d94db799-mn42n   1/1     Running             0               4m43s
k8sforbeginner-57d94db799-nxkw8   1/1     Running             0               4m43s
k8sforbeginner-57d94db799-sjzx4   1/1     Running             0               4m43s
nginx-6799fc88d8-d6w82            1/1     Running             2 (6h35m ago)   35d
web-5bfc9bc56d-qd5g7              1/1     Running             3 (6h35m ago)   35d
web-5bfc9bc56d-vdpqz              1/1     Running             3 (6h35m ago)   35d
web-5bfc9bc56d-xllx8              1/1     Running             3 (6h35m ago)   35d
```


delete a pod will restart a new one auto

```bash
[blue@master ~]$ kubectl get pod
NAME                              READY   STATUS    RESTARTS        AGE
k8s-for-beginners                 1/1     Running   0               5h29m
k8sforbeginner-57d94db799-4tm4k   1/1     Running   0               80s
k8sforbeginner-57d94db799-579k5   1/1     Running   0               80s
k8sforbeginner-57d94db799-8kslx   1/1     Running   0               6m2s
k8sforbeginner-57d94db799-bqvn8   1/1     Running   0               80s
k8sforbeginner-57d94db799-dbsb5   1/1     Running   0               81s
k8sforbeginner-57d94db799-dx5g7   1/1     Running   0               80s
k8sforbeginner-57d94db799-k427r   1/1     Running   0               10m
k8sforbeginner-57d94db799-mn42n   1/1     Running   0               6m2s
k8sforbeginner-57d94db799-nxkw8   1/1     Running   0               6m2s
k8sforbeginner-57d94db799-sjzx4   1/1     Running   0               6m2s
nginx-6799fc88d8-d6w82            1/1     Running   2 (6h36m ago)   35d
web-5bfc9bc56d-qd5g7              1/1     Running   3 (6h36m ago)   35d
web-5bfc9bc56d-vdpqz              1/1     Running   3 (6h36m ago)   35d
web-5bfc9bc56d-xllx8              1/1     Running   3 (6h36m ago)   35d
[blue@master ~]$ kubectl delete pod k8sforbeginner-57d94db799-4tm4k
pod "k8sforbeginner-57d94db799-4tm4k" deleted

[blue@master ~]$ 
[blue@master ~]$ kubectl get pod
NAME                              READY   STATUS    RESTARTS        AGE
k8s-for-beginners                 1/1     Running   0               5h29m
k8sforbeginner-57d94db799-579k5   1/1     Running   0               105s
k8sforbeginner-57d94db799-8kslx   1/1     Running   0               6m27s
k8sforbeginner-57d94db799-bqvn8   1/1     Running   0               105s
k8sforbeginner-57d94db799-ctghl   1/1     Running   0               7s
k8sforbeginner-57d94db799-dbsb5   1/1     Running   0               106s
k8sforbeginner-57d94db799-dx5g7   1/1     Running   0               105s
k8sforbeginner-57d94db799-k427r   1/1     Running   0               11m
k8sforbeginner-57d94db799-mn42n   1/1     Running   0               6m27s
k8sforbeginner-57d94db799-nxkw8   1/1     Running   0               6m27s
k8sforbeginner-57d94db799-sjzx4   1/1     Running   0               6m27s
nginx-6799fc88d8-d6w82            1/1     Running   2 (6h37m ago)   35d
web-5bfc9bc56d-qd5g7              1/1     Running   3 (6h37m ago)   35d
web-5bfc9bc56d-vdpqz              1/1     Running   3 (6h37m ago)   35d
web-5bfc9bc56d-xllx8              1/1     Running   3 (6h37m ago)   35d
```

describe pod

```bash
[blue@master ~]$ kubectl describe pod k8sforbeginner-57d94db799-579k5
Name:         k8sforbeginner-57d94db799-579k5
Namespace:    default
Priority:     0
Node:         node1/192.168.0.12
Start Time:   Tue, 05 Oct 2021 15:52:27 +0800
Labels:       pod-template-hash=57d94db799
              tier=simplehttp
Annotations:  cni.projectcalico.org/containerID: cb605e5ef4d62545b10fc2929c1c37662dc072eb63a0549c83391e1fa408851a
              cni.projectcalico.org/podIP: 10.244.166.177/32
              cni.projectcalico.org/podIPs: 10.244.166.177/32
Status:       Running
IP:           10.244.166.177
IPs:
  IP:           10.244.166.177
Controlled By:  ReplicaSet/k8sforbeginner-57d94db799
Containers:
  k8sforbeginner:
    Container ID:   docker://8f6a6681bff88aa8d6bc2b855116b8abf3c0a3108d841d23912a9afbafe5239f
    Image:          k8s-for-beginners:v0.0.1
    Image ID:       docker://sha256:6cdee0a8bb610d3d48224e73e8c8ababfdd0dd765b4505c546196af23349a25e
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Tue, 05 Oct 2021 15:52:30 +0800
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-wwg9l (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  kube-api-access-wwg9l:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  17m   default-scheduler  Successfully assigned default/k8sforbeginner-57d94db799-579k5 to node1
  Normal  Pulled     17m   kubelet            Container image "k8s-for-beginners:v0.0.1" already present on machine
  Normal  Created    17m   kubelet            Created container k8sforbeginner
  Normal  Started    17m   kubelet            Started container k8sforbeginner

```


delete the deploy

```bash
[blue@master ~]$ kubectl delete deploy k8sforbeginner
deployment.apps "k8sforbeginner" deleted

```

