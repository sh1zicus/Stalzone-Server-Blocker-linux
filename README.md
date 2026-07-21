# Stalzone Server Blocker

TUI-приложение для управления блокировкой серверов Stalzone через nftables. Автоматически применяет правила при запуске игры и снимает при остановке.

## Возможности

- Интерактивный TUI для выбора серверов
- Пинг серверов с цветовой индикацией
- Автоматическое применение nft-правил при обнаружении `stalzone.exe`
- Автоматическая очистка правил при остановке игры
- Поиск по серверам
- systemd сервис для работы в фоне

## Требования

- Linux (nftables)
- Go 1.21+
- `pgrep` (обычно есть в `procps`)
- `jq` (для скриптов)
- Права на управление nftables (`CAP_NET_ADMIN`)

## Установка

### Из исходников

```bash
git clone https://github.com/sh1zicus/Stalzone-Server-Blocker-linux.git
cd stalzone-server-blocker-linux
make build
```

### Права на nftables

Бинарник нуждается в правах `CAP_NET_ADMIN` для управления nftables:

```bash
# Вариант 1: через setcap (рекомендуется)
sudo setcap cap_net_admin+ep ./stalzone-blocker

# Вариант 2: запуск через sudo
sudo ./stalzone-blocker
```

## Использование

### TUI (интерактивный режим)

```bash
./stalzone-blocker
```

#### Управление

| Клавиша | Действие |
|---------|----------|
| `↑` `↓` | Навигация |
| `←` `→` | Свернуть/развернуть пул |
| `␣` | Включить/выключить сервер |
| `Enter` | Включить/выключить весь пул |
| `R` | Обновить пинг |
| `A` | Применить правила |
| `D` | Сбросить правила |
| `/` | Поиск |
| `q` | Выход |

#### Цвета пинга

- 🟢 Зелёный: ≤30 ms
- 🟡 Жёлтый: ≤60 ms
- 🟠 Оранжевый: ≤90 ms
- 🔴 Красный: >90 ms
- ⚪ Серый: таймаут

### Daemon (фоновый режим)

Daemon следит за процессом `stalzone.exe` и автоматически управляет nft-правилами:

```bash
./stalzone-blocker --daemon
```

#### Логика работы

1. При обнаружении `stalzone.exe` → задержка 2 сек → применение правил
2. При остановке `stalzone.exe` → немедленная очистка правил

## Systemd сервис

### Установка

```bash
# Скопировать файл сервиса
cp stalzone-blocker.service ~/.config/systemd/user/

# Перезагрузить конфигурацию
systemctl --user daemon-reload

# Включить автозапуск
systemctl --user enable stalzone-blocker.service

# Запустить сейчас
systemctl --user start stalzone-blocker.service
```

### Управление

```bash
# Статус
systemctl --user status stalzone-blocker.service

# Остановить
systemctl --user stop stalzone-blocker.service

# Перезапустить
systemctl --user restart stalzone-blocker.service

# Логи
journalctl --user -u stalzone-blocker -f

# Логи за последний час
journalctl --user -u stalzone-blocker --since "1 hour ago"
```

### Настройка пути

Отредактируйте `stalzone-blocker.service`:

```ini
[Service]
ExecStart=%h/Projects/stalzone-server-blocker/stalzone-blocker --daemon
```

Замените путь на ваш实际ный путь к бинарнику.

## Конфигурация

Конфигурация хранится в `~/.config/stalzone-blocker/config.json`:

```json
{
    "selected": ["SERVER-1", "SERVER-2"],
    "theme": "dark",
    "sort": "ping",
    "window": {
        "width": 120,
        "height": 35
    }
}
```

### Поля

| Поле | Описание |
|------|----------|
| `selected` | Список включённых серверов |
| `theme` | Тема оформления |
| `sort` | Сортировка серверов |
| `window` | Размер окна TUI |

## Структура проекта

```
stalzone-blocker/
├── cmd/
│   └── stalzone-blocker/
│       └── main.go          # Точка входа
├── internal/
│   ├── app/
│   │   └── app.go           # Инициализация TUI
│   ├── config/
│   │   ├── config.go        # Загрузка/сохранение конфига
│   │   └── selection.go     # Управление выбором серверов
│   ├── daemon/
│   │   └── daemon.go        # Фоновый мониторинг
│   ├── logger/
│   │   └── logger.go        # Логирование
│   ├── model/
│   │   ├── config.go        # Модель конфигурации
│   │   ├── pool.go          # Модель пула/туннеля
│   │   ├── servers_json.go  # Структура Servers.json
│   │   └── state.go         # Состояние TUI
│   ├── nft/
│   │   ├── apply.go         # Применение правил
│   │   ├── check.go         # Проверка прав
│   │   └── reset.go         # Сброс правил
│   ├── ping/
│   │   └── worker.go        # Пинг серверов
│   ├── servers/
│   │   └── loader.go        # Загрузка Servers.json
│   ├── tui/
│   │   ├── model.go         # Модель TUI
│   │   ├── style.go         # Стили Lipgloss
│   │   ├── update.go        # Обработка сообщений
│   │   └── view.go          # Отрисовка
│   └── util/
│       ├── files.go         # Поиск файлов
│       └── paths.go         # Пути к файлам
├── data/
│   └── Servers.json         # Список серверов
├── go.mod
├── go.sum
├── Makefile
└── stalzone-blocker.service # Systemd сервис
```

## Разработка

### Сборка

```bash
# Обычная сборка
make build

# Сборка + установка прав
make setcap

# Запуск
make run

# Очистка
make clean
```

### Структура данных

Серверы хранятся в `data/Servers.json`:

```json
{
    "mode": "roxy",
    "pools": [
        {
            "name": "GAME-EU",
            "region": "EU",
            "tunnels": [
                {
                    "name": "GAME-EU-2",
                    "address": "79.127.241.67:29450"
                }
            ]
        }
    ]
}
```

### Добавление нового сервера

1. Отредактируйте `data/Servers.json`
2. Добавьте запись в соответствующий пул
3. Перезапустите приложение или нажмите `R` для обновления пинга

### Тестирование nft правил

```bash
# Проверить текущие правила
sudo nft list table inet stalzone_blocker

# Ручное применение
./stalzone-blocker
# Нажмите A

# Ручной сброс
# Нажмите D
```

## Troubleshooting

### "нет прав на управление nftables"

```bash
sudo setcap cap_net_admin+ep ./stalzone-blocker
```

### Сервис не запускается

```bash
# Проверить логи
journalctl --user -u stalzone-blocker -n 50

# Проверить путь в сервисе
cat ~/.config/systemd/user/stalzone-blocker.service
```

### pgrep не находит процесс

Убедитесь, что `stalzone.exe` запущен и процесс доступен через:

```bash
pgrep -f "stalzone.exe"
```

## Лицензия

MIT
