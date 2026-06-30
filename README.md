# Stalzone-Server-Blocker-linux

скриптяра блочит все сервера мск и оставляет только екб и нск

## как запустить
```
git clone https://github.com/sh1zicus/Stalzone-Server-Blocker-linux
```
```
cd Stalzone-Server-Blocker-linux
```
сгенерирует ip адреса для блокировки
```
chmod +x generate.sh

./generate.sh
```
запушит адреса в nft
```
sudo nft -f stalcraft.nft
```
