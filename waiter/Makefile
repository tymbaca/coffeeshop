export REGISTRY=registry.tymbaca

build:
	docker buildx create --use
	docker buildx build --platform linux/arm64,linux/amd64 -t coffeeshop/waiter:latest -f ./kube/Dockerfile .
	docker tag coffeeshop/waiter $(REGISTRY)/coffeeshop/waiter:latest
	docker push $(REGISTRY)/coffeeshop/waiter:latest

apply:
	kubectl apply -f kube/deployment.yaml
	kubectl apply -f kube/service.yaml

delete:
	kubectl delete -f kube/deployment.yaml
	kubectl delete -f kube/service.yaml
