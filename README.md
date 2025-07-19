[qtun]: https://github.com/shadowsocks/qtun
[quic-tun]: https://github.com/kungze/quic-tun

# qt

I am qt(pronounced "Q-Tee"), not Qt(pronounced "cute").
This page is only available in Chinese. Please use AI to translate it into the language you need.

## 故事

稣经常在弱网环境下访问一些云原生服务。

> 比如蹭商店的不安全网络……

听说 KcpTun 加速效果不错，稣就测试了一番，结果发现它对多 IPv6 地址的支持有些问题，详见 [#1005](https://github.com/xtaci/kcptun/issues/1005)。

> Server 有多个 IPv6 地址，假设为 A、B，存在一种可能客户端 C 连接的是 A 地址，然而 Server 使用 B 来回复 C，那么 C 处的防火墙有可能丢弃来自 B 的回复。

稣想起 QUIC 应该也能做出效果不错的 Tunnel，并且应该对多 IPv6 地址支持很完善，于是找了一下，果然已经有人做好了，比如：

- [qtun][qtun] Yet another SIP003 plugin based on IETF-QUIC

- [quic-tun][quic-tun] A fast and security tunnel based on QUIC, make you can access remote TCP/UNIX application like a local application. 一个快速且安全的 TCP 隧道工具，能加速弱网环境下（如网络有丢包）TCP 的转发性能。

其中，[qtun][qtun] 是 Rust 写的，目前 QUIC 库可能还是 C、C++、Go 的最成熟，所以放弃。而另一个 [quic-tun][quic-tun] 已经存档，作者投入另一个更强的项目去了。

稣认为写这类不赚钱的工具，还是 Go 语言更合适，排除使用 C、C++，所以打算参考第二个~~自己写一个~~，名字就叫做 qt。

算了，稣太忙，大家贡献代码吧！

## 设计

### 1. 操作系统兼容性

- qt-server 主要运行于 Debian、OpenWRT。

- qt-client 主要运行于 Debian、macOS、OpenWRT、Windows。

### 2. 安全性

- 证书必须是自签名证书，既安全又容易获得。如果是从第三方申请获得的，那安全性存疑，毕竟越少人知道越好。像“[qtun][qtun] will look for TLS certificates signed by acme.sh by default.”这样的，是又麻烦又不一定安全。

- 证书的密钥长度需要达标，比如 RSA 至少需要 4096bit，而 ECC 至少 386bit。

- 客户端必须验证服务端，防止有人伪装服务端。如果真正的连接目标（Tunnel 后面的服务端）使用明文协议，那么有人伪造 Tunnel 服务端，是可能拿到明文信息的。

### 3. 命令行接口

- 采用 Boost.Program_options 格式。即选项名前面有两个 -，名字全小写，单词以 - 分隔，比如 --remote-port，可选短名为 -p。

- 支持配置文件，JSON 或 TOML 格式。通常来说，命令行的优先级高于配置文件。
