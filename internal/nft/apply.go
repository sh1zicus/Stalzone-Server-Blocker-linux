package nft

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

func Apply(pools []model.Pool) error {

	// Фаза 1: удаление старой таблицы (если есть).
	deleteTable()

	// Фаза 2: создание новой таблицы с правилами.
	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("nftables: нет прав или ошибка подключения: %w", err)
	}

	table := conn.AddTable(&nftables.Table{
		Family: nftables.TableFamilyINet,
		Name:   TableName,
	})

	chain := conn.AddChain(&nftables.Chain{
		Name:     ChainName,
		Table:    table,
		Type:     nftables.ChainTypeFilter,
		Hooknum:  nftables.ChainHookOutput,
		Priority: nftables.ChainPriorityFilter,
	})

	blocked := blockedIPs(pools)

	for ipStr := range blocked {

		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}

		ip4 := ip.To4()
		if ip4 == nil {
			continue
		}

		conn.AddRule(&nftables.Rule{
			Table: table,
			Chain: chain,
			Exprs: []expr.Any{
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseNetworkHeader,
					Offset:       16,
					Len:          4,
				},
				&expr.Cmp{
					Op:       expr.CmpOpEq,
					Register: 1,
					Data:     ip4,
				},
				&expr.Verdict{
					Kind: expr.VerdictDrop,
				},
			},
		})
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("nftables flush: %w", err)
	}

	return nil
}

func deleteTable() {
	conn, err := nftables.New()
	if err != nil {
		return
	}

	conn.DelTable(&nftables.Table{
		Family: nftables.TableFamilyINet,
		Name:   TableName,
	})

	// Игнорируем ошибку — таблицы может не быть.
	conn.Flush()
}

func blockedIPs(pools []model.Pool) map[string]struct{} {

	allowed := make(map[string]struct{})
	blocked := make(map[string]struct{})

	for _, pool := range pools {
		for _, t := range pool.Tunnels {
			ip := t.Address
			if i := strings.Index(ip, ":"); i >= 0 {
				ip = ip[:i]
			}
			if t.Selected {
				allowed[ip] = struct{}{}
			} else {
				blocked[ip] = struct{}{}
			}
		}
	}

	for ip := range allowed {
		delete(blocked, ip)
	}

	return blocked
}
