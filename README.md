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
