## Step 1 create `k8s-pageview-deploy.yaml`

First, create pageview deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-pageview
spec:
  replicas: 1
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend
    spce:
      containers:
      - name: k8s-pageview
        image: pageview:v0.0.1 

```


## Step 2 create redis deployment `k8s-redis-deploy.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-redis
spec:
  replicas: 1
  selector:
    metchLabels:
      tier: backend
  template:
    metadata:
      labels:
        tier: backend
    spec:
      containers:
      - name: k8s-redis
        image: redis

```

## Step 3 make the redis as a service with name as db `k8s-redis-internal-service.yaml`

```yaml
kind: Service
apiVersion: v1
metadata:
  name: db
spec:
  selector:
    tier: backend
  ports:
  - port: 6379
    targetPort: 6379

```

## Step 4 apply the three yaml file:

```bash
[blue@master pageview]$ kubectl apply -f k8s-redis-deploy.yaml
deployment.apps/k8s-redis created
[blue@master pageview]$ kubectl apply -f k8s-redis-internal-service.yaml 
service/db created
[blue@master pageview]$ kubectl apply -f k8s-pageview-deploy.yaml 
error: error validating "k8s-pageview-deploy.yaml": error validating data: ValidationError(Deployment.spec.template): unknown field "spce" in io.k8s.api.core.v1.PodTemplateSpec; if you choose to ignore these errors, turn validation off with --validate=false
[blue@master pageview]$ vim k8s-pageview-deploy.yaml 
[blue@master pageview]$ kubectl apply -f k8s-pageview-deploy.yaml 
deployment.apps/k8s-pageview created

```


## Step 5 check redis service

```bash
[blue@master pageview]$ kc get svc db
NAME   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
db     ClusterIP   10.102.130.157   <none>        6379/TCP   2m58s
[blue@master pageview]$ telnet 10.102.130.157 6379
Trying 10.102.130.157...
Connected to 10.102.130.157.
Escape character is '^]'.

```

## Step 6 apply the pageview app

k8s-pageview-service.yaml

```yaml
kind: Service
apiVersion: v1
metadata:
  name: k8s-pageview
spce:
  selector:
    tier: frontend
  type: NodePort
  ports:
  - port: 80
    targetPort: 8080 

```

```bash
[blue@master pageview]$ kc apply -f k8s-pageview-service.yaml 
service/k8s-pageview created
```


## Step 7 check the service

```
[blue@master pageview]$ curl http://192.168.0.12:32467
Hello, you're visitor #5.
[blue@master pageview]$ curl http://192.168.0.12:32467
Hello, you're visitor #6.
[blue@master pageview]$ curl http://192.168.0.12:32467
Hello, you're visitor #7.
[blue@master pageview]$ curl http://192.168.0.12:32467
Hello, you're visitor #8.
[blue@master pageview]$ curl http://192.168.0.13:32467
Hello, you're visitor #9.
[blue@master pageview]$ curl http://192.168.0.13:32467
Hello, you're visitor #10.
[blue@master pageview]$ curl http://192.168.0.13:32467
Hello, you're visitor #11.
[blue@master pageview]$ curl http://192.168.0.13:32467
Hello, you're visitor #12.
```
