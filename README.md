# qt

I am qt(pronounced "Q-Tee"), not Qt(pronounced "cute").
This page is only available in Chinese. Please use AI to translate it into the language you need.

## 故事

稣经常在一些弱网环境下访问一些云原生服务。听说 KcpTun 加速效果不错，稣就测试了一番，结果发现它对多 IPv6 地址的支持有些问题，详见 [#1005](https://github.com/xtaci/kcptun/issues/1005)。

> Server 有多个 IPv6 地址，假设为 A、B，存在一种可能客户端 C 连接的是 A 地址，然而 Server 使用 B 来回复 C，那么 C 处的防火墙有可能丢弃来自 B 的回复。

稣想起 QUIC 应该也能做出效果不错的 Tunnel，并且应该对多 IPv6 地址支持很完善，于是找了一下，果然已经有人做好了，比如：

- [qtun](https://github.com/shadowsocks/qtun) Yet another SIP003 plugin based on IETF-QUIC

- [quic-tun](https://github.com/kungze/quic-tun) About
A fast and security tunnel based on QUIC, make you can access remote TCP/UNIX application like a local application. 一个快速且安全的 TCP 隧道工具，能加速弱网环境下（如网络有丢包）TCP 的转发性能。

其中，第一个是 Rust 写的，目前 QUIC 库可能还是 C、C++、Go 的最成熟，所以放弃。第二个已经存档，作者投入另一个更强的项目去了。

稣认为写这类不赚钱的工具，还是 Go 语言更合适，排除使用 C、C++，所以打算参考第二个~~自己写一个~~，名字就叫做 qt。

算了，稣太忙，大家贡献代码吧！
