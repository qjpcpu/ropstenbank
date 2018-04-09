package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"net/http"
	"time"
)

var pwd = "123456"

var (
	eth_node string
	to_addr  string
	do_trans bool
)

func init() {
	flag.StringVar(&eth_node, "eth", "", "")
	flag.StringVar(&to_addr, "to", "", "")
	flag.BoolVar(&do_trans, "trans", false, "")
	flag.Parse()
}

func main() {
	if eth_node == "" || to_addr == "" {
		fmt.Println("参数缺失")
		return
	}
	ks := keystore.NewKeyStore("./bank-keystore", 2, 1)
	max := 100
	if total := len(ks.Accounts()); total < max {
		for i := total; i < max; i++ {
			ks.NewAccount(pwd)
		}
		fmt.Printf("创建%d个账户\n", max)
	}

	if !do_trans {
		for i, account := range ks.Accounts() {
			fmt.Printf("%v.发送rosten网络eth到%s\n", i, account.Address.Hex())
			http.Get(fmt.Sprintf("http://faucet.ropsten.be:3001/donate/%s", account.Address.Hex()))
		}

		fmt.Println("发送完成,等待3分钟eth到账....")
		time.Sleep(1 * time.Minute)
	}
	fmt.Printf("开始转账到 %s\n", to_addr)

	conn, err := ethclient.Dial(eth_node)
	if err != nil {
		fmt.Println("连接eth节点失败:", err)
		return
	}

	threshold := big.NewInt(10000000000000000) // 0.01 ETH
	networkId, _ := conn.NetworkID(context.Background())
	var nonce uint64
	for _, account := range ks.Accounts() {
		balance, err := conn.BalanceAt(context.Background(), account.Address, nil)
		if err != nil || balance.Cmp(threshold) < 0 {
			fmt.Printf("账户%s余额不足0.01eth\n", account.Address.Hex())
			continue
		}
		fmt.Printf("账户%s余额为%.3feth\n", account.Address.Hex(), asfloat(balance))
		amount := new(big.Int).Div(new(big.Int).Mul(balance, big.NewInt(999)), big.NewInt(1000))
		tx := types.NewTransaction(nonce, common.HexToAddress(to_addr), amount, 21000, big.NewInt(20000000000), nil)
		ks.Unlock(account, pwd)
		signed, err := ks.SignTx(account, tx, networkId)
		if err != nil {
			fmt.Println("sign fail:", err)
			continue
		}
		ks.Lock(account.Address)
		if err = conn.SendTransaction(context.Background(), signed); err != nil {
			fmt.Printf("%s => %s 转账%.3feth失败:%v\n", account.Address.Hex(), to_addr, asfloat(amount), err)
		} else {
			fmt.Printf("%s => %s 转账%.3feth, tx:%s\n", account.Address.Hex(), to_addr, asfloat(amount), tx.Hash().Hex())
		}
	}
}

func asfloat(num *big.Int) float64 {
	one_eth := big.NewFloat(1000000000000000000)
	f, _ := new(big.Float).Quo(new(big.Float).SetInt(num), one_eth).Float64()
	return f
}
