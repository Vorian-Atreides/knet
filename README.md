# knet

Kubernetes services will always do a round robin load balancing between your pods,
which is fine in most of the time but isn't enough in many other.

Knet is a simple library, taking advantage of the Kubernetes' headless services
to give you more tools to build your services.

## Resolver

Resolver helps to retrieve Pod's hostnames from the Kubernetes DNS.

## Consistent hashing

Will help you to define a consistent routing from a given resource to a given
destination. The hashring package does not require to run in Kubernetes and
can be imported and used in any environment.

The package does not force you to use any third party library
allowing you to bring your own Tree datastructure,
which only need to be wrapped to respect the interface.

Example with `github.com/emirpasic/gods/trees/redblacktree`:
```
type Node struct {
	*redblacktree.Node
}

func (n Node) Value() hashring.Node {
	return n.Node.Value.(hashring.Node)
}

func (n Node) Key() uint64 {
	return n.Node.Key.(uint64)
}

func (n Node) Left() hashring.Node {
	if n.Node.Left == nil {
		return nil
	}
	return Node{n.Node.Left}
}

func (n Node) Right() hashring.Node {
	if n.Node.Right == nil {
		return nil
	}
	return Node{n.Node.Right}
}

type Tree struct {
	tree *redblacktree.Tree
}

func (t Tree) Root() hashring.Node {
	return Node{t.tree.Root}
}

func (t Tree) Put(key uint64, n hashring.Node) {
	t.tree.Put(key, n)
}

func (t Tree) Remove(key uint64) {
	t.tree.Remove(key)
}

func (t Tree) Get(key uint64) hashring.Node {
	if v, ok := t.tree.Get(key); ok {
		return v.(hashring.Node)
	}
	return nil
}
```

## Headless service

To help you prototype your application, you can use the below instructions to
setup a dummy application with its headless service.

First, create the deployment of your service:
> $ kubectl create deployment hello-node --image=gcr.io/hello-minikube-zero-install/hello-node

Then create the headless service with the selector on the previous deployment:
> $ kubectl expose deployment hello-node --cluster-ip=None --port=8080

