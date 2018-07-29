# emseco

Extrēmus是基于DAG(有向无环图)+Sharding(分表扩容)结构的强大基础架构。Extrēmus的核心是建立一个网络，可以支撑的安全、高效和去信任的生态系统。为了满足日益增长的交易吞吐量，Extrēmus必须确保网络中的所有节点都能一致地确认交易。为此，共识算法的设计是最重要的部分

## 运行生成蓝图和构造蓝图顺序测试用例

生成蓝图节点测试，先引入依赖，再运行测试用例

引入依赖的以太坊程序

    go get github.com/ethereum/go-ethereum/common

运行测试程序

    go test -v emseco/core/
