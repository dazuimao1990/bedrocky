# Bedrocky

<!-- 向 SRE 角色转化的练手开发项目。 -->

基岩(bedrock)是陆地表层中的坚硬岩层。一般多被土层覆盖，埋藏深度不一，少则数米到数十米，多则数百米。由沉积岩、变质岩、岩浆岩中的一种或数种岩类组成，可作大型建筑工程的地基。

为了项目名字的唯一性，基于 bedrock 造了一个词。这个项目旨在创建一种类似 ansible 的机制，能够让交付人员方便的在 Linux 环境中快速部署一些基础设施类软件。工作过程中尽量减少对环境的依赖，专注于离线场景的应用。

# 主要特性

bedrocky:
- [x] 从 Yaml 配置文件中读取目标主机集群，包括IP、端口、密码/密钥、角色。
- [x] 从 Yaml 配置文件中读取 Role 细节，根据模块的不同，批量执行任务。
- [x] 引入 shell 模块作为最基础的能力，面向多个主机批量执行 shell 命令。
<!-- - [ ] 仿照 ansible-playbook 设计 job 模块，执行 containerd docker k8s/k3s etcd nfs 等基础服务的批量安装
	- 解决各服务的优化配置：
		- etcd：加自动压缩
		- containerd docker 引用国内源，日志自动压缩，支持导入预置镜像
		- k8s 证书100年，可扩展指令参数
		- k3s 支持导入预置镜像 -->
<!-- - [ ] 引入 k8s 模块
	- staticpod 子模块，规划批量启动静态pod
	- daemonset 子模块，规划批量启动守护进程 -->

# 设计思路

## 核心抽象

1. TaskChan: 定义一个任务队列，利用 golang 语言的并行优势，通过 ssh 协议批量对远端主机进行操作。数据类型 `chan TaskSpec`
2. TaskSpec: 定义一次任务中涉及到的所有远程主机（host）信息，以及针对每个主机即将执行的任务集（role），结构体字段包括 ip port username password/publickey role。数据类型 `struct`
3. role: 定义一个角色（role），包括角色的名字（如 `installdocker`），和角色所承担的工作列表（steps）。
4. step: 工作列表中的每一个步骤，根据 `module` 的不同，处理的方式也不一样。
5. module: 一个接口，任何实现了该接口的动作都可以在 `step` 中集成进来。

## 实现的顺序

- [x] 基于 cobra、viper 框架构建命令行基础，实现基本的命令行参数传递与调用
- [x] 实现最基础的能力，通过 ssh 操作远端主机
- [x] 实现基于 yaml 类型配置文件解析来输出配置
- [x] 定义 moudle 的基础规范，并将 shell 实现
- [x] 将多个 shell step 操作合并成为一个 role ，实现顺序执行
- [x] 实现 taskchan 队列，利用 goroutine 实现并发多个 host 的任务执行
- [ ] 开始扩充 moudle 的不同实现类型，如 copy 、systemd 

## 编译

拉取代码编译即可：

```bash
git clone https://github.com/dazuimao1990/bedrocky.git
cd bedrocky
go build -o bedrocky *.go
```
## 示例

代码中提供了样例配置文件 `bedrocky.yaml`

使用以下命令来运行：

```bash
./bedrocky play -f bedrocky.yaml
```

输出类似：
```
2023/01/04 18:19:50 [SUCCESS] 主机 127.0.0.1 任务 hostname 执行完成. 执行结果: MacBook-Pro-FrozenFire.local
2023/01/04 18:19:50 [SUCCESS] 主机 192.168.2.59 任务 hostname 执行完成. 执行结果: MacBook-Pro-FrozenFire.local
2023/01/04 18:19:50 [SUCCESS] 主机 127.0.0.1 任务 uname -v 执行完成. 执行结果: Darwin Kernel Version 21.3.0: Wed Jan  5 21:37:58 PST 2022; root:xnu-8019.80.24~20/RELEASE_ARM64_T8101
2023/01/04 18:19:50 [SUCCESS] 主机 192.168.2.59 任务 uname -v 执行完成. 执行结果: Darwin Kernel Version 21.3.0: Wed Jan  5 21:37:58 PST 2022; root:xnu-8019.80.24~20/RELEASE_ARM64_T8101
2023/01/04 18:19:50 [SUCCESS] 主机 127.0.0.1 任务 arch 执行完成. 执行结果: arm64
2023/01/04 18:19:50 [SUCCESS] 主机 192.168.2.59 任务 arch 执行完成. 执行结果: arm64
```

命令也接纳命令行参数输入，方便基于一条命令面向某个远程主机批量执行命令：

```bash
./bedrocky run --host=127.0.0.1 -P someword -u root -c "hostname,arch" -c "uname -r"
```

输出类似：
```
2023/01/04 18:23:18 [SUCCESS] 主机 127.0.0.1 任务 hostname 执行完成. 执行结果: MacBook-Pro-FrozenFire.local
2023/01/04 18:23:18 [SUCCESS] 主机 127.0.0.1 任务 arch 执行完成. 执行结果: arm64
2023/01/04 18:23:18 [SUCCESS] 主机 127.0.0.1 任务 uname -r 执行完成. 执行结果: 21.3.0
```

## bedrocky.yaml 格式

### hosts 

主机列表中包含了一组远程主机信息，主要包括 ssh 连接所需信息，以及即将执行的任务角色，单个主机的全量信息包括：

```yaml
- host: 127.0.0.1
  username: root
  port: 33
  password: yyyyyyyy
  publickey: ~/.ssh/other.pub
  role:
  - checksystem
```

其中 `host` `role` 字段不可省略。

### port

全局定义的 ssh 端口号。

### username

全局定义的 ssh 用户名。

### password

全局定义的 ssh 密码。

### publickey

全局定义的 ssh 公钥。

### rolename

此处定义了一个 role 的详情，shell 示例：

```yaml
checksystem:
  - name: prepare
    module: shell
    spec:
      - hostname
      - uname -v
  - name: checkarch
    module: shell
    spec:
      - arch 
```

## 主机信息优先级

bedrocky 中提供三种方式定义连接一个远程主机所需要的用户名、密码/密钥、端口信息：

- 命令行参数：bedrocky 将命令行参数的规定置为最高优先级。
- bedrocky.yaml 中在 `hosts` 列表中为每个主机定义的用户名、密码/密钥、端口为第二优先级。
- bedrocky.yaml 中全局定义的用户名、密码/密钥、端口为最低优先级

当你无法确定时，使用以下命令，来了解并行任务中如何连接远程主机：

```bash
./bedrocky hosts show -f bedrocky.yaml --password 123asd --port 34
Using config file: bedrocky.yaml
序号    主机            端口      用户    密码                    公钥
1       127.0.0.1       34      root    MTIzYXNk                ~/.ssh/id_rsa.pub
2       192.168.0.1     34      root    MTIzYXNk                ~/.ssh/id_rsa.pub
```