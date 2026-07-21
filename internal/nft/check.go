package nft

import (
	"fmt"

	"github.com/google/nftables"
)

const (
	TableName = "stalzone_blocker"
	ChainName = "output"
)

func CheckPermissions() error {

	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf(
			"нет прав на управление nftables (CAP_NET_ADMIN).\n"+
				"Выполните:\n"+
				"  sudo setcap cap_net_admin+ep <бинарник>\n"+
				"или запустите через sudo",
		)
	}

	// Пробуем прочитать таблицы — если нет прав, ошибка вылезет здесь.
	_, err = conn.ListTables()
	if err != nil {
		return fmt.Errorf(
			"нет прав на управление nftables (CAP_NET_ADMIN).\n"+
				"Выполните:\n"+
				"  sudo setcap cap_net_admin+ep <бинарник>\n"+
				"или запустите через sudo",
		)
	}

	return nil
}
