export REGISTRY=localhost:5000

build:
	docker buildx build --platform linux/amd64 -t $(REGISTRY)/coffeeshop/waiter:latest -f ./kube/Dockerfile .
	docker push $(REGISTRY)/coffeeshop/waiter:latest

apply:
	kubectl apply -f kube/deployment.yaml
	kubectl apply -f kube/service.yaml

delete:
	kubectl delete -f kube/deployment.yaml
	kubectl delete -f kube/service.yaml
