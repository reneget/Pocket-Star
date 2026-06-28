# My Pretty Star ⭐

**My Pretty Star** — site-to-site VPN orchestrator. Превращает старый компьютер в личный сервер без белого IP.

## Проблема

У вас есть старый ПК, который вы хотите использовать как сервер (файловое хранилище, нейросети, и т.д.), но провайдер не даёт белый IP. Аренда VPS с мощностями вашего ПК стоит дорого.

## Решение

1. **Снимаете самый дешёвый VPS** (от $3/мес) — он будет хабом с белым IP
2. **Ваш сервер** (старый ПК) подключается к хаб через AmneziaWG
3. **Ваш рабочий ноутбук** подключается к тому же хабу
4. Всё — вы в одной виртуальной локальной сети. Сервер доступен как `10.0.0.2`

```
VPS (Хаб)        Старый ПК (Сервер)    Ноутбук (Клиент)
10.0.0.1 ─────── 10.0.0.2             10.0.0.3
    │                  │                    │
    └── AmneziaWG ─────┴────────────────────┘
```

## Быстрый старт

### 1. На VPS (хаб)

```bash
# Установка
curl -sSL https://github.com/anomalyco/my-pretty-star/raw/main/scripts/quick-start.sh | bash

# Инициализация хаба с мастер-паролем (+ Docker)
mps hub init --docker --pass
# → Создаст: hub.conf, docker-compose.yml, клиентские конфиги
# → Введи мастер-пароль (он защищает конфиги узлов)

# Запуск
docker compose -f mps-data/docker-compose.yml up -d
```

### 2. На старом ПК (сервер)

```bash
# Получить конфиг с хаба
mps node join --hub <IP_ВАШЕГО_VPS>:51820 --docker --name server
# → Введи мастер-пароль (тот же, что на хабе)

# Запуск
docker compose -f mps-data/docker-compose.yml up -d
```

**Или без Docker** — скопируй `mps-data/clients/server.conf` с хаба и импортируй в WireGuard/AmneziaWG GUI.

### 3. На ноутбуке (клиент)

```bash
mps node join --hub <IP_ВАШЕГО_VPS>:51820 --name work-laptop
# Введи мастер-пароль
# → Подключено!
```

**Или через GUI:** скопируй `work-laptop.conf` → импорт в WireGuard/Amnezia → Connect.

### 4. Проверка

```bash
ping 10.0.0.2               # Пингуем сервер
ssh user@10.0.0.2           # SSH на сервер
scp file user@10.0.0.2:~/  # Копируем файл
```

## Команды

| Команда | Описание |
|---------|----------|
| `mps hub init --docker --pass` | Инициализация хаба на VPS |
| `mps hub status` | Список подключённых узлов |
| `mps node join --hub X` | Подключение узла (автоматически) |
| `mps node status` | Статус VPN-соединения |
| `mps module list` | Список модулей |
| `mps decrypt <file>` | Расшифровать конфиг мастер-паролем |
| `mps doctor` | Диагностика системы |
| `mps version` | Версия |

## Master-пароль

Конфиги узлов шифруются AES-256-GCM мастер-паролем. Это защищает ключи доступа при передаче через SCP или другие каналы.

## Архитектура

Проект построен на **модульной системе**. Любой новый функционал (SFTP, проброс портов, веб-панель) добавляется как отдельный модуль.

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

## Разработка

```bash
# Сборка
go build -o mps ./cmd/mps

# Тесты
go test ./...

# Сборка для другой платформы
GOOS=linux GOARCH=arm64 go build -o mps-arm64 ./cmd/mps
```

## Лицензия

MIT
