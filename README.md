# Stalzone-Server-Blocker-linux

скриптяра блочит все сервера мск и оставляет только екб и нск

## как запустить
сгенерирует ip адреса для блокировки
```
./generate.sh
```
запушит адреса в nft
```
sudo nft -f stalcraft.nft
```
