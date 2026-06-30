#!/usr/bin/env bash
set -e

SERVERS="$HOME/Stalcraft-Server-Blocker-linux/Servers.json"
OUT="stalcraft.nft"

cat > "$OUT" <<'EOF'
table inet stalcraft {
    chain output {
        type filter hook output priority 0;
        policy accept;
EOF

jq -r '
.pools[]
| select(.region == "RU")
| select(.name != "NSK1" and .name != "EKB")
| .tunnels[]
| .address
| split(":")[0]
' "$SERVERS" | sort -u | while read -r ip; do
    echo "        ip daddr $ip drop" >> "$OUT"
done

cat >> "$OUT" <<'EOF'
    }
}
EOF

echo "Файл $OUT создан."
