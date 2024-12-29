upgrade:
	helm upgrade --namespace gitlab --timeout 120s --set gitlabUrl="https://gitlab.tymbaca" --set runnerToken="glrt-t1_efFbiWgtsJGNzZHrQxva" --values values.yaml gitlab-runner gitlab/gitlab-runner

reinstall:
	helm uninstall -n gitlab gitlab-runner
	helm install --namespace gitlab --timeout 120s --set gitlabUrl="https://gitlab.tymbaca" --set runnerToken="glrt-t1_efFbiWgtsJGNzZHrQxva" --values values.yaml gitlab-runner gitlab/gitlab-runner
