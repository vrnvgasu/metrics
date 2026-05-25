# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Профилирование памяти (iter17)

Оптимизация `pkg/compress.GzipCompress`: заменён `gzip.NewWriter` на `sync.Pool` для переиспользования `gzip.Writer` и `bytes.Buffer`.

Результат бенчмарка:

| | ns/op | B/op | allocs/op |
|---|---|---|---|
| До | 115 561 | 814 746 | 21 |
| После | 10 450 | 456 | 1 |

Сравнение профилей (`pprof -top -diff_base=profiles/base.pprof profiles/result.pprof`):

```
File: server
Build ID: d01aadaa2e53793abdf3120ef97d63c3dc6f7eba
Type: inuse_space
Time: 2026-05-11 23:34:31 MSK
Showing nodes accounting for -3076.69kB, 24.48% of 12566.53kB total
Dropped 1 node (cum <= 62.83kB)
      flat  flat%   sum%        cum   cum%
   -1539kB 12.25% 12.25%    -1539kB 12.25%  runtime.allocm
 -521.05kB  4.15% 16.39%  -521.05kB  4.15%  encoding/xml.map.init.0
  516.01kB  4.11% 12.29%   516.01kB  4.11%  io.init.func1
     514kB  4.09%  8.20%      514kB  4.09%  bufio.NewReaderSize (inline)
     513kB  4.08%  4.11%      513kB  4.08%  bufio.NewWriterSize (inline)
  512.69kB  4.08% 0.035%   512.69kB  4.08%  regexp/syntax.(*compiler).inst (inline)
  512.22kB  4.08%  4.04%   512.22kB  4.08%  runtime.malg
 -512.17kB  4.08% 0.034%  -512.17kB  4.08%  net/http.Header.Clone (inline)
 -512.17kB  4.08%  4.11%  -512.17kB  4.08%  net/textproto.readMIMEHeader
 -512.07kB  4.07%  8.18%  -512.07kB  4.07%  net/url.parse
 -512.05kB  4.07% 12.26%  -512.05kB  4.07%  runtime.acquireSudog
 -512.05kB  4.07% 16.33%  -512.05kB  4.07%  sync.runtime_notifyListWait
 -512.02kB  4.07% 20.41%  -512.02kB  4.07%  github.com/go-playground/validator/v10.lazyRegexCompile (inline)
 -512.02kB  4.07% 24.48%  -512.02kB  4.07%  reflect.packEface
```
