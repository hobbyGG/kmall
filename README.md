# kmall

## 项目相关问题

### 常见的项目文件管理方式有哪些？

- monorepo
- mutirepo+submodule

### 微服务项目中，如何管理pb文件？

proto文件相同，且protoc也要相同。

直接将项目放在github上或其他公用代码库上可以吗？有什么缺点

应该是可行的，只要能获取到proto文件就可以，但是不好保证protoc的版本一致性，并且有些项目不开源，对于没有代码仓库的公司可能办不到。

可以使用git-submodule模式，将proto文件单独抽离出来放在代码库中

### kratos如何使用consul?

kratos的轻量化思想使得kratos通过第三方插件就可以支持各种主流或者公司内部的注册中心，

kratos默认将grpc服务注册到9000端口号上，所以grpc客户端也可以直接写死9000端口号，不需要使用consul。

对于服务端。使用consul需要在Newapp中添加kratos.Registrar选项，该选项可以接入各种注册中心，包括consul、etcd、nacos等。我们需要在服务启动时注册RPC服务，所以可以将创建registrar的代码存放在server.go中。在NewRegistrar中我们使用kratos提供的consul插件，使用consul api提供的new方法得到与consul交互的client，配置好该client的地址和端口后，用consul.New包装这个客户端就可以得到一个registrar。使用依赖注入的方式将reg放入Newapp中。

对于客户端。将grpc的client放在data层的结构体中以供业务调用。client通过consul自动获取。客户端通过dial并配置grpc.WithDiscovery(discovery)获取服务，dial的endpoint在使用consul时就可以用discovery:///[服务名]自动解析。自动解析需要与consul的连接，连接的建立也需要先创建一个consul client与client进行交互，这与server端一致，因为consul.New返回的Registry实现了discover和registrar两个接口，所以代码相同。