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



## Profiles diff:

$ go tool pprof -top -diff_base=profiles/base.pprof profiles/result2.pprof
File: handler.test
Build ID: d813951dfeb6409dc10cb49ff3682f8c973b88f5
Type: alloc_space
Time: 2026-05-14 22:49:29 MSK
Showing nodes accounting for -31.89MB, 0.13% of 23732.10MB total
Dropped 2 nodes (cum <= 118.66MB)
      flat  flat%   sum%        cum   cum%
  310.21MB  1.31%  1.31%   310.21MB  1.31%  bufio.NewReaderSize (inline)
  -95.38MB   0.4%  0.91%  -424.08MB  1.79%  github.com/FeshLig/metrcollector/internal/handler.(*RootHandler).RootPage
  -90.01MB  0.38%  0.53%  -129.16MB  0.54%  bytes.(*Buffer).grow
   89.01MB  0.38%   0.9%    89.01MB  0.38%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).Updates
     -85MB  0.36%  0.54%  -593.13MB  2.50%  github.com/gin-gonic/gin.(*Context).String (inline)
     -82MB  0.35%   0.2%      -82MB  0.35%  net/http/httptest.NewRecorder (inline)
   68.03MB  0.29%  0.48%    68.03MB  0.29%  encoding/json.(*Decoder).refill
   63.52MB  0.27%  0.75%    63.52MB  0.27%  net/http.(*Request).WithContext (inline)
  -63.13MB  0.27%  0.49%   -65.19MB  0.27%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).SnapshotCounterMetrics
  -59.74MB  0.25%  0.23%   -79.24MB  0.33%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).SnapshotGaugeMetrics
  -51.50MB  0.22% 0.017%   -51.50MB  0.22%  bytes.(*Buffer).String (inline)
   51.50MB  0.22%  0.23%    58.01MB  0.24%  encoding/json.Marshal
  -44.01MB  0.19% 0.048%   -44.01MB  0.19%  net/http.Header.Clone (inline)
   44.01MB  0.19%  0.23%    44.01MB  0.19%  net/url.parse
   43.51MB  0.18%  0.42%    43.51MB  0.18%  reflect.growslice
  -39.50MB  0.17%  0.25%  -253.04MB  1.07%  github.com/FeshLig/metrcollector/internal/handler.(*ValueJSONHandler).UpdateFromJSON
  -39.15MB  0.16% 0.086%   -39.15MB  0.16%  bytes.growSlice
   38.01MB  0.16%  0.25%   116.02MB  0.49%  net/http.readRequest
  -37.50MB  0.16% 0.088%   -37.50MB  0.16%  html/template.htmlReplacer
  -33.50MB  0.14% 0.053%     -127MB  0.54%  reflect.Value.call
   30.50MB  0.13% 0.075%    30.50MB  0.13%  net/textproto.readMIMEHeader
     -28MB  0.12% 0.043%   -41.50MB  0.17%  reflect.MakeSlice
   22.50MB 0.095% 0.052%    30.50MB  0.13%  encoding/json.(*decodeState).literalStore
  -19.50MB 0.082%  0.03%   -19.50MB 0.082%  github.com/FeshLig/metrcollector/internal/repository.(*MemStorage).SnapshotGauges
   16.50MB  0.07% 0.039%   524.63MB  2.21%  github.com/gin-gonic/gin.(*Context).JSON (inline)
      15MB 0.063%   0.1%       15MB 0.063%  encoding/json.NewDecoder (inline)
      15MB 0.063%  0.17%   360.09MB  1.52%  github.com/FeshLig/metrcollector/internal/handler.(*UpdatesHandler).Updates
     -14MB 0.059%  0.11%      -24MB   0.1%  github.com/FeshLig/metrcollector/internal/handler.(*UpdateURLHandler).UpdateFromURL
     -14MB 0.059% 0.048%      -14MB 0.059%  net.IP.String
  -13.62MB 0.057% 0.0095%   -15.62MB 0.066%  internal/fmtsort.Sort
  -13.50MB 0.057% 0.066%   -13.50MB 0.057%  reflect.unsafe_NewArray
      10MB 0.042% 0.024%       10MB 0.042%  net/textproto.canonicalMIMEHeaderKey
   -9.50MB  0.04% 0.064%   483.75MB  2.04%  net/http/httptest.NewRequestWithContext
    8.50MB 0.036% 0.028%   -48.52MB   0.2%  github.com/FeshLig/metrcollector/internal/handler.(*UpdateJSONHandler).UpdateFromJSON
   -7.50MB 0.032%  0.06%    -7.50MB 0.032%  encoding/json.(*scanner).pushParseState
   -7.01MB  0.03%  0.09%    -7.01MB  0.03%  net/textproto.MIMEHeader.Set (inline)
    5.50MB 0.023% 0.066%     5.50MB 0.023%  io.NopCloser (inline)
      -5MB 0.021% 0.087%       -5MB 0.021%  reflect.unsafe_New
      -4MB 0.017%   0.1%    26.50MB  0.11%  encoding/json.(*decodeState).object
       4MB 0.017% 0.087%        4MB 0.017%  net/textproto.(*Reader).ReadLine (inline)
      -4MB 0.017%   0.1%       -4MB 0.017%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).Get
       4MB 0.017% 0.087%        4MB 0.017%  internal/strconv.FormatInt
   -3.50MB 0.015%   0.1%    -3.50MB 0.015%  fmt.Sprint
       3MB 0.013%  0.09%        3MB 0.013%  bytes.NewReader (inline)
   -2.05MB 0.0086% 0.098%    -2.05MB 0.0086%  github.com/FeshLig/metrcollector/internal/repository.(*MemStorage).SnapshotCounters
      -2MB 0.0084%  0.11%       -2MB 0.0084%  strings.NewReader (inline)
   -1.55MB 0.0065%  0.11%    -1.55MB 0.0065%  regexp.(*bitState).reset
   -1.50MB 0.0063%  0.12%    -1.50MB 0.0063%  github.com/gin-gonic/gin/render.writeContentType
   -1.50MB 0.0063%  0.13%    -1.50MB 0.0063%  net/textproto.NewReader (inline)
    1.50MB 0.0063%  0.12%     1.50MB 0.0063%  internal/strconv.FormatFloat (inline)
   -0.51MB 0.0022%  0.12%    -0.51MB 0.0022%  github.com/FeshLig/metrcollector/internal/repository.(*MemStorage).AddCounter
   -0.51MB 0.0022%  0.12%    -0.51MB 0.0022%  golang.org/x/net/html.map.init.1
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  github.com/go-playground/validator/v10.map.init.10
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  sync.(*Pool).pinSlow
    0.50MB 0.0021%  0.13%     0.50MB 0.0021%  github.com/ugorji/go/codec.helperDecDriverMsgpackBytes.fastpathDList
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  vendor/golang.org/x/net/http2/hpack.init
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  fmt.init.func1
   -0.50MB 0.0021%  0.13%       -1MB 0.0042%  encoding/json.newEncodeState
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  sync.(*poolChain).pushHead
    0.50MB 0.0021%  0.13%     0.50MB 0.0021%  github.com/gin-gonic/gin.SetMode
   -0.50MB 0.0021%  0.13%    -0.50MB 0.0021%  github.com/gin-gonic/gin.(*Engine).allocateContext (inline)

$ go tool pprof -top -diff_base=profiles/service_base.pprof profiles/service_result2.pprof
File: service.test
Build ID: c0f6140fbd5e41700dff34b53a4970f5e0509b6d
Type: alloc_space
Time: 2026-05-14 23:02:26 MSK
Showing nodes accounting for -1598.48MB, 23.71% of 6740.95MB total
Dropped 50 nodes (cum <= 33.70MB)
      flat  flat%   sum%        cum   cum%
-1785.39MB 26.49% 26.49% -1726.37MB 25.61%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).SnapshotGaugeMetrics
  180.39MB  2.68% 23.81%   179.89MB  2.67%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).Updates
   59.02MB  0.88% 22.93%    59.02MB  0.88%  github.com/FeshLig/metrcollector/internal/repository.(*MemStorage).SnapshotGauges
  -52.50MB  0.78% 23.71%   -52.50MB  0.78%  github.com/FeshLig/metrcollector/internal/service.(*MetricServiceImpl).Get
         0     0% 23.71%        8MB  0.12%  github.com/FeshLig/metrcollector/internal/service.BenchmarkMetricService_GetCounter
         0     0% 23.71%   -60.50MB   0.9%  github.com/FeshLig/metrcollector/internal/service.BenchmarkMetricService_GetGauge
         0     0% 23.71% -1726.37MB 25.61%  github.com/FeshLig/metrcollector/internal/service.BenchmarkMetricService_SnapshotGaugeMetrics
         0     0% 23.71%   179.89MB  2.67%  github.com/FeshLig/metrcollector/internal/service.BenchmarkMetricService_Updates
         0     0% 23.71% -1598.98MB 23.72%  testing.(*B).launch
         0     0% 23.71% -1598.48MB 23.71%  testing.(*B).runN
