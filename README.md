```
yum install wget -y && wget https://golang.org/dl/go1.17.linux-amd64.tar.gz && sudo tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz && echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && echo 'export GOPATH=$HOME/go' >> ~/.bashrc && source ~/.bashrc && curl -OL https://raw.githubusercontent.com/SPEED4G-Official/Proxy-Vncloud/refs/heads/main/ServiceProxy.go && go build ServiceProxy.go
```
```
bash <(curl -Ls https://raw.githubusercontent.com/SPEED4G-Official/Proxy-Vncloud/refs/heads/main/install.sh)
```
