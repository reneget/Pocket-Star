# ⭐ Pocket Star

<div align="center">

<img src="pkg/icon/logo.svg" width="120" alt="Pocket Star Logo"/>

[![Readme in English](https://img.shields.io/badge/README-English-blue.svg)](README.md)
![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)
![Go](https://img.shields.io/badge/go-1.25+-00ADD8.svg)
![AmneziaWG](https://img.shields.io/badge/vpn-AmneziaWG-green.svg)
![Docker](https://img.shields.io/badge/docker-ready-blue.svg)
![License](https://img.shields.io/badge/license-MIT-yellow.svg)

**Преврати старый компьютер в личный сервер без белого IP. Твой дом, доступный откуда угодно.**

[📖 Документация](#документация) • [🚀 Быстрый старт](#быстрый-старт) • [🏗️ Архитектура](#архитектура) • [🛠️ Технологический стек](#технологический-стек) • [📦 Модули](#модули)

</div>

---

## 📋 Содержание

- [О проекте](#о-проекте)
- [Архитектура](#архитектура)
- [Технологический стек](#технологический-стек)
- [Быстрый старт](#быстрый-старт)
- [Команды](#команды)
- [Модули](#модули)
- [Мастер-пароль](#мастер-пароль)
- [Разработка](#разработка)
- [Документация](#документация)

---

## 🎯 О проекте

У тебя есть старый ПК, который хочется использовать как сервер (файлопомойка, локальные нейросети, дев-среда), но провайдер не даёт белый IP. Арендовать VPS с мощностями твоего старого железа — дорого.

Вместо этого ты снимаешь **самый дешёвый VPS** (~$3/мес) — он работает хабом с белым IP. Твой сервер и ноутбук подключаются к нему через AmneziaWG, образуя виртуальную локальную сеть. Всё, сервер доступен откуда угодно.

### Основные возможности

✨ **Для пользователей:**
- 🔌 Site-to-site VPN за минуты — старый ПК и ноутбук в одной сети
- 🔒 AES-256-GCM шифрование конфигов мастер-паролем
- 🐳 Docker Compose для хаба и сервера
- 🖥️ Поддержка нативных WG-клиентов — импорт .conf и коннект
- 🔧 Модульная архитектура — расширяй плагинами

---

## 🏗️ Архитектура

### Схема сети

```mermaid
graph TB
    subgraph "Интернет"
        VPS[VPS Хаб<br/>Белый IP<br/>10.0.0.1]
    end
    
    subgraph "Домашняя сеть"
        SRV[Старый ПК / Сервер<br/>10.0.0.2]
    end
    
    subgraph "Откуда угодно"
        LAP[Ноутбук / Клиент<br/>10.0.0.3]
    end
    
    VPS <-->|AmneziaWG<br/>Туннель| SRV
    VPS <-->|AmneziaWG<br/>Туннель| LAP
    SRV <-.->|SSH / Ping / SCP<br/>через VPN| LAP
    
    style VPS fill:#4ecdc4,color:#fff
    style SRV fill:#ff6b6b,color:#fff
    style LAP fill:#45b7d1,color:#fff
```

### Структура проекта

```mermaid
graph LR
    subgraph "CLI"
        PSTAR[pstar<br/>Точка входа]
    end
    
    subgraph "Пакеты ядра"
        CLI[CLI Команды]
        VPN[VPN Движок<br/>AmneziaWG]
        CRYPT[Крипто<br/>AES-256-GCM]
        MOD[Система Модулей]
        DOCK[Docker<br/>Compose]
    end
    
    subgraph "Модули"
        CORE[Core Module]
        SFTP[SFTP Module<br/>Планируется]
        FWD[Port Forward<br/>Планируется]
        TUI[TUI Module<br/>Планируется]
        MEDIA[Media Storage<br/>Планируется]
        CLOUD[Cloud Sync<br/>Планируется]
        MON[Monitoring<br/>Планируется]
    end
    
    subgraph "Дистрибуция"
        HUB_IMG[Docker Hub Image]
        NODE_IMG[Docker Node Image]
        SCRIPTS[Setup Scripts]
    end
    
    PSTAR --> CLI
    PSTAR --> VPN
    PSTAR --> CRYPT
    PSTAR --> MOD
    PSTAR --> DOCK
    MOD --> CORE
    MOD --> SFTP
    MOD --> FWD
    MOD --> TUI
    MOD --> MEDIA
    MOD --> CLOUD
    MOD --> MON
    DOCK --> HUB_IMG
    DOCK --> NODE_IMG
    CLI --> SCRIPTS
    
    style PSTAR fill:#4ecdc4,color:#fff
    style CLI fill:#96ceb4,color:#fff
    style VPN fill:#ff9ff3,color:#fff
    style CRYPT fill:#feca57,color:#000
    style MOD fill:#45b7d1,color:#fff
    style CORE fill:#ff6b6b,color:#fff
    style SFTP fill:#ff6b6b,color:#fff
    style FWD fill:#ff6b6b,color:#fff
    style TUI fill:#ff6b6b,color:#fff
    style MEDIA fill:#ff6b6b,color:#fff
    style CLOUD fill:#ff6b6b,color:#fff
    style MON fill:#ff6b6b,color:#fff
```

### Последовательность подключения

```mermaid
sequenceDiagram
    participant VPS as 🖥️ VPS Хаб
    participant SRV as 🖧 Старый ПК
    participant LAP as 💻 Ноутбук
    
    Note over VPS,LAP: Шаг 1: Инициализация хаба
    VPS->>VPS: pstar hub init --docker --pass
    VPS->>VPS: Генерация ключей, конфигов, compose
    VPS->>VPS: docker compose up -d
    
    Note over VPS,LAP: Шаг 2: Сервер подключается
    SRV->>VPS: pstar node join --hub IP:51820
    VPS-->>SRV: Назначен IP 10.0.0.2
    SRV->>SRV: docker compose up -d
    
    Note over VPS,LAP: Шаг 3: Клиент подключается
    LAP->>VPS: pstar node join --hub IP:51820
    VPS-->>LAP: Назначен IP 10.0.0.3
    
    Note over VPS,LAP: Шаг 4: Проверка связи
    LAP->>SRV: ping 10.0.0.2
    SRV-->>LAP: ✅ pong
    LAP->>SRV: ssh user@10.0.0.2
    SRV-->>LAP: ✅ Подключено
```

---

## 🛠️ Технологический стек

### Ядро
- **Go 1.25+** — один бинарник, кроссплатформенная компиляция
- **AmneziaWG** — VPN с обфускацией, совместим с WireGuard
- **cobra** — CLI-фреймворк

### Безопасность
- **Curve25519** — криптографический обмен ключами
- **AES-256-GCM** — шифрование конфигов мастер-паролем

### Инфраструктура
- **Docker & Docker Compose** — контейнеризация хаба и сервера
- **Alpine Linux** — базовый образ для Docker

### VPN протоколы
- **AmneziaWG** — VPN по умолчанию с защитой от DPI
- **WireGuard** — полная совместимость, импорт .conf в любой клиент

---

## 🚀 Быстрый старт

### Предварительные требования
- Docker и Docker Compose (для хаба и сервера)
- Дешёвый VPS с белым IP (любой провайдер)
- Linux / macOS на клиенте

### Установка

```bash
curl -sSL https://github.com/reneget/Pocket-Star/raw/main/scripts/quick-start.sh | bash
```

### 1. Инициализация хаба на VPS

```bash
pstar hub init --docker --pass
# → Создаёт: hub.conf, docker-compose.yml, зашифрованные конфиги узлов
# → Придумай мастер-пароль (он защищает конфиги узлов)

docker compose -f pstar-data/docker-compose.yml up -d
```

### 2. Подключение сервера (старый ПК)

```bash
pstar node join --hub <IP_ТВОЕГО_VPS>:51820 --docker --name server
# → Введи мастер-пароль (тот же, что задал на хабе)

docker compose -f pstar-data/docker-compose.yml up -d
```

**Без Docker?** Скопируй `pstar-data/clients/server.conf` с хаба — импортируй в WireGuard / AmneziaWG GUI.

### 3. Подключение клиента (ноутбук)

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

---

## ⌨️ Команды

| Команда | Описание |
|---------|----------|
| `pstar hub init --docker --pass` | Инициализация хаба на VPS |
| `pstar hub status` | Список подключённых узлов |
| `pstar node join --hub X [--docker]` | Подключить эту машину как узел |
| `pstar node status` | Статус VPN-соединения |
| `pstar module list` | Список модулей |
| `pstar module enable <name>` | Включить модуль |
| `pstar module disable <name>` | Отключить модуль |
| `pstar decrypt <файл>` | Расшифровать конфиг мастер-паролем |
| `pstar doctor` | Диагностика системы |
| `pstar version` | Версия |

---

## 📦 Модули

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

### Доступные

- **core** — VPN + маршрутизация (MVP) ✅
- **tui** — терминальный интерфейс (вдохновлён opencode) ✅

### Планируемые

- **sftp** — файлообменник через SFTP
- **port-fwd** — проброс портов через хаб
- **dashboard** — веб-интерфейс
- **media** — медиа-сервер (доступ с телефона/TV; на базе Jellyfin/Immich)
- **cloud** — облачное хранилище и синхронизация (как OneDrive; на базе Nextcloud)
- **monitoring** — мониторинг состояния системы (Pulse или Uptime Kuma)

---

## 🔒 Мастер-пароль

Конфиги узлов шифруются **AES-256-GCM** мастер-паролем. Это защищает приватные ключи при передаче через SCP или другие каналы.

```bash
# На хабе: конфиги сохраняются как *.conf.enc
# На узле: запусти `pstar decrypt server.conf.enc` и введи пароль
# Или просто: `pstar node join` — пароль запросится автоматически
```

---

## 🔧 Разработка

### Сборка

```bash
go build -o pstar ./cmd/pstar
```

### Тесты

```bash
go test ./...
```

### Кросс-компиляция

```bash
GOOS=linux GOARCH=arm64 go build -o pstar-arm64 ./cmd/pstar
GOOS=linux GOARCH=amd64 go build -o pstar-linux ./cmd/pstar
GOOS=darwin GOARCH=amd64 go build -o pstar-macos ./cmd/pstar
```

### Структура проекта

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

---

## 📊 Статистика проекта

```mermaid
pie title Распределение кода по пакетам
    "VPN Engine (AmneziaWG)" : 30
    "CLI Commands" : 25
    "Crypto (AES-256-GCM)" : 15
    "Module System" : 15
    "Docker / Scripts" : 10
    "Tests" : 5
```

---

## 📖 Документация

- [README.en.md](./README.md) — английская версия

---

## 🤝 Вклад в проект

Мы приветствуем вклад в развитие проекта! Пожалуйста:

1. Создайте форк проекта
2. Создайте ветку для новой функции (`git checkout -b feature/AmazingFeature`)
3. Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`)
4. Отправьте в ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

---

## 📝 Лицензия

Этот проект лицензирован под MIT License — см. файл [LICENSE](LICENSE) для деталей.

---

## 👥 Авторы

- **[reneget](https://github.com/reneget)**

---

## ⭐ История звёзд

<div align="center">

[![Star History Chart](https://api.star-history.com/svg?repos=reneget/Pocket-Star&type=Date)](https://star-history.com/#reneget/Pocket-Star&Date)

</div>

---

<div align="center">

**Сделано с ❤️ чтобы дать вторую жизнь старым ПК**

⭐ Если проект был полезен, поставьте звезду!

</div>
