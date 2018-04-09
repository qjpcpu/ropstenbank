# 以太坊ropsten网络自动充币

参数:

* `--eth`, 以太坊ropsten节点,建议连接自建节点,否则ropsten共有节点容易被封ip
* `--to`, 转出到你的地址
* `--trans`, 发币+转账或者仅转账

# examples

```
# 发eth同时转账
ropstenbank --eth http://localhost:18545 --to 0x86bb2d1c849bb37ea160ABDaa7C4e722e38364A0
# 仅转账(适用于将到账延迟的账户转出)
ropstenbank --eth http://localhost:18545 --to 0x86bb2d1c849bb37ea160ABDaa7C4e722e38364A0 --trans transfer
```

# attention

请仅在测试网络使用
