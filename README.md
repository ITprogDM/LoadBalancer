# Load Balancer

## 📌 Описание проекта

Балансировщик нагрузки с реализованным алгоритмом round-robin.
Сервис был написан с использованием Clean Architecture, что позволяет легко расширять функционал и тестировать сервис. Также реализован Graceful Shutdown для корректного завершения работы.

## 🛠️ Используемые технологии

**PostgreSQL** (в качестве Базы Данных)

**Docker** (для запуска сервиса)

**golang-migrate/migrate** (для миграций БД)

**pgx** (драйвер для работы с PostgreSQL)

## 🔧 Для запуска сервиса необходимо:

В качестве конфига используется json-файл, находящийся в корне проекта в папке pkg/config - соответственно его нужно заполнить под себя.
Заполнить .env файл. Указать путь к конфиг файлу PATH_CONFIG. Если вы меняете в .env файле PG_URL, то для миграций также нужно поменять в Makefile раздел migrate.

Запустить сервис:
```
make start
```
Применить миграции БД:
```
make migrate
```
Остановить сервис:
```
make stop
```
В проекте в разделе pkg/utils есть два сервера (backend1.go, backend2.go) для локального тестирования балансировщика. 

Запустить нагрузку Apache Bench с помощью Docker:
```
make ab-test
```
Прикладываю результаты AB-тестирования:
![image](https://github.com/user-attachments/assets/571552b2-919d-4b88-98cb-164c41808554)

Concurrency Level:      1000
Time taken for tests:   7.709 seconds
Complete requests:      5000
Failed requests:        0
Total transferred:      840000 bytes
HTML transferred:       255000 bytes
Requests per second:    648.57 [#/sec] (mean)
Time per request:       1541.861 [ms] (mean)
Time per request:       1.542 [ms] (mean, across all concurrent requests)
Transfer rate:          106.41 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        1  437  97.3    475     559
Processing:    57 1080 169.5   1089    1406
Waiting:       10  774 162.3    754    1263
Total:         58 1517 142.4   1557    1773
