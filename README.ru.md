# Pocket Star

> Преврати старый компьютер в личный сервер без белого IP. Твой дом, доступный откуда угодно.

У тебя есть старый ПК, который хочется использовать как сервер (файлопомойка, локальные нейросети, дев-среда), но провайдер не даёт белый IP. Арендовать VPS с мощностями твоего старого железа — дорого.

Вместо этого ты снимаешь **самый дешёвый VPS** (~$3/мес) — он работает хабом с белым IP. Твой сервер и ноутбук подключаются к нему через AmneziaWG, образуя виртуальную локальную сеть. Всё, сервер доступен откуда угодно.

```
VPS (Хаб)         Старый ПК (Сервер)   Ноутбук (Клиент)
10.0.0.1 ──────── 10.0.0.2             10.0.0.3
    │                   │                   │
    └── AmneziaWG ──────┴───────────────────┘
    Docker Compose      Docker Compose       Нативный WG GUI
                                             или pstar CLI
```

## Быстрый старт

### 1. На VPS (хаб)

```bash
curl -sSL https://github.com/anomalyco/pocket-star/raw/main/scripts/quick-start.sh | bash

pstar hub init --docker --pass
# → Создаёт: hub.conf, docker-compose.yml, зашифрованные конфиги для узлов
# → Придумай мастер-пароль (он защищает конфиги узлов)

docker compose -f pstar-data/docker-compose.yml up -d
```

### 2. На старом ПК (сервер)

```bash
pstar node join --hub <IP_ТВОЕГО_VPS>:51820 --docker --name server
# → Введи мастер-пароль (тот же, что задал на хабе)

docker compose -f pstar-data/docker-compose.yml up -d
```

**Без Docker?** Скопируй `pstar-data/clients/server.conf` с хаба — импортируй в WireGuard / AmneziaWG GUI.

### 3. На ноутбуке (клиент)

```bash
pstar node join --hub <IP_ТВОЕГО_VPS>:51820 --name work-laptop
# Введи мастер-пароль → подключено.
```

**Или через GUI:** скопируй `work-laptop.conf` → импорт в WireGuard/Amnezia → Connect.

### 4. Проверка

```bash
ping 10.0.0.2               # Пингуем сервер
ssh user@10.0.0.2           # SSH
scp file user@10.0.0.2:~/  # Файлы туда-сюда
```

## Команды

| Команда | Описание |
|---------|----------|
| `pstar hub init --docker --pass` | Инициализация хаба на VPS |
| `pstar hub status` | Список подключённых узлов |
| `pstar node join --hub X` | Подключить эту машину как узел |
| `pstar node status` | Статус VPN-соединения |
| `pstar module list` | Список модулей |
| `pstar decrypt <файл>` | Расшифровать конфиг мастер-паролем |
| `pstar doctor` | Диагностика системы |
| `pstar version` | Версия |

## Мастер-пароль

Конфиги узлов шифруются **AES-256-GCM** мастер-паролем. Это защищает приватные ключи при передаче через SCP или другие каналы.

```bash
# На хабе: конфиги сохраняются как *.conf.enc
# На узле: запусти `pstar decrypt server.conf.enc` и введи пароль
# Или просто: `pstar node join` — пароль запросится автоматически
```

## Модульная система

Всё строится как подключаемые модули. Ты можешь написать свой, реализовав один интерфейс:

```go
type Module interface {
    Name() string
    Priority() int
    Dependencies() []string
    Init(*Context) error
    Start() error
    Stop() error
    Status() Status
}
```

### Планируемые модули

- **core** — VPN + маршрутизация (MVP)
- **sftp** — файлообменник через SFTP
- **port-fwd** — проброс портов через хаб
- **dashboard** — веб-интерфейс

## Архитектура

```
pocket-star/
├── cmd/pstar/           # Точка входа
├── pkg/
│   ├── cli/             # CLI команды (hub, node, module, doctor)
│   ├── config/          # AES-256-GCM шифрование
│   ├── vpn/             # Генерация ключей и конфигов AmneziaWG
│   ├── module/          # Интерфейс модулей + реестр
│   └── docker/          # Генерация docker-compose
├── modules/
│   └── core/            # core-модуль (VPN)
├── images/
│   ├── hub/Dockerfile
│   └── node/Dockerfile
└── scripts/
    ├── quick-start.sh   # curl → bash установщик
    └── bootstrap.sh     # Полная автоматическая настройка
```

## Разработка

```bash
# Сборка
go build -o pstar ./cmd/pstar

# Тесты
go test ./...

# Кросс-компиляция
GOOS=linux GOARCH=arm64 go build -o pstar-arm64 ./cmd/pstar
```

## Лицензия

MIT
