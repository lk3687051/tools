1. 赋予genhost 可执行权限
   chmod 777 genhost

2. 修改aiops密码 hosts
   ansible_password=Aiops@123

3. 运行 genhost 生成hosts 文件
   ./genhost

4. 运行ansible脚本修复
   ansible-playbook all -i hosts  reconfig.yml -M ./modules/

生成可执行文件
CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -ldflags="-w -s" -o ansible/reconfig/modules/sensu-config-amd cmd/reconfig/main.go
CGO_ENABLED=0 GOOS="linux" GOARCH="arm64" go build -ldflags="-w -s" -o ansible/reconfig/modules/sensu-config-arm cmd/reconfig/main.go
CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -ldflags="-w -s" -o ansible/reconfig/genhost cmd/allhost/main.go
