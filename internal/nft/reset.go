package nft

import (
	"fmt"

	"github.com/google/nftables"
)

func Reset() error {

	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("nftables: нет прав или ошибка подключения: %w", err)
	}

	conn.DelTable(&nftables.Table{
		Family: nftables.TableFamilyINet,
		Name:   TableName,
	})

	// Игнорируем ошибку — таблицы может не быть.
	conn.Flush()

	return nil
}
