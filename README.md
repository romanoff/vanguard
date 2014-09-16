Vanguard
========

Vanguard is a docker orchestration tool for multihost environment. It relies on etcd for sharing data between hosts and weave to create network between docker containers. Host machines that will be orchestrated should have both etcd and weave up and running.

Starting with Etcd
------------------

First of all download latest release of [etcd](https://github.com/coreos/etcd/releases) and put binary files to folder specified in paths (/bin/ for example). This has to be performed on each machine in your cluster. After that you have to connect etcd into cluster.

Let's say you have 2 machines that you want to connect with following ip addresses: 173.194.115.66, 173.194.115.69.

Here is how you can start etcd on first machine (with 173.194.115.66 ip address)

```
etcd -peer-addr 173.194.115.66:7001 -addr 173.194.115.66:4001 -bind-addr 0.0.0.0
```

Now you can start etcd on second server and connect etcd to first host

```
etcd -peer-addr 173.194.115.69:7001 -addr 173.194.115.69:4001 -peers 173.194.115.66:7001 -bind-addr 0.0.0.0
```

You can check that etcd works correctly by getting cluster machines through curl:

```
curl -L http://127.0.0.1:4001/v2/machines
```

Should return 2 items.

Starting weave
--------------

Download and put it to paths on each host:

```
sudo wget -O /bin/ https://raw.githubusercontent.com/zettio/weave/master/weaver/weave
sudo chmod a+x /bin/weave
```

For the next step [docker](http://docs.docker.com/installation/) should be installed.

On first host (with 173.194.115.66 ip address) :

```
sudo weave launch 10.0.0.1/16
```

After that start it on the second host and connect them:

```
sudo weave launch 10.0.0.2/16 173.194.115.66
```

Starting vanguard agent
-----------------------

On each host start vanguard agent:

```
sudo vanguard agent
```

Vanguard agent performs tasks like starting and stopping containers.
