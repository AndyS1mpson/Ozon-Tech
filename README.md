## Реализация проекта в рамках курса "Go-middle разработчик" от Ozon Tech


# Описание системы

1. LOMS (Logistics and Order Management System) - сервис, отвечающий за учет заказов и логистику.
2. Checkout - сервис, отвечающий за корзину пользователя и оформление заказа.
3. Notifications - сервис, отвечающий за отправку уведомлений.
4. ProductService - внешний сервис, который предоставляет информацию о товарах.

### Путь покупки товаров
* Checkout.addToCart  
    * добавляем в корзину и проверяем, что есть в наличии)
* Можем удалять из корзины
* Можем получать список товаров корзины
    * название и цена тянутся из ProductService.get_product
* Приобретаем товары через Checkout.purchase
    * идем в LOMS.createOrder и создаем заказ
    * У заказа статус new
    * LOMS резервирует нужное количество единиц товара
    * Если не удалось зарезервить, заказ падает в статус failed
    * Если удалось, падаем в статус awaiting payment
* Оплачиваем заказ
    * Вызываем LOMS.orderPayed
    * Резервы переходят в списание товара со склада
    * Заказ идет в статус payed
* Можно отменить заказ до оплаты
    * Вызываем LOMS.cancelOrder
    * Все резервирования по заказу отменяются, товары снова доступны другим пользователям
    * Заказ переходит в статус cancelled
    * LOMS должен сам отменять заказы по таймауту, если не оплатили в течение 10 минут


## Схема баз данных
<img width="1253" alt="image" src="https://github.com/deerc-dev/openedx-admin/assets/45228812/beb54696-c907-41e5-ab69-6302a6d53801">

# Локальная разработка

1. В папку deployments необходимо добавить .env файл по аналогии с example.env, в котором указаны параметры подключения к базам данных микросервисом 
2. В папки checkout, loms, notifications необходимо добавить файлы config.yaml по аналогии с config.example.yaml
3. В папку certs необходимо добавить свой SSL сертификат
4. Запустить окружение для мониторинга логов и подождать пока все поднимется:
> make run-log-env
5. Если в предыдущем шаге все сервисы поднялись без ошибок, в отдельном терминале поднять окружение сервисов:
> make run-services


## Создание миграций
1. Чтобы выполнить миграцию баз данных, необходимо установить библиотеку goose:
> go install github.com/pressly/goose/v3/cmd/goose@latest

2. Перейти в папку migration соответствующего микросервиса и выполнить команду:
> goose create \*имя миграции\* sql

3. В терминале импортировать соответствующую переменную окружения:
> CH_POSTGRES_URL  # connection string для сервиса checkout  
> LOMS_POSTGRES_URL # connection string для сервиса loms  
> NOTIF_POSTGRES_URL # connection string для сервиса notifications
4. выполнить комманду:
> make migrate

## Мониторинг логов
1. Необходимо через браузер зайти в Graylog (по умолчанию порт 7555)
2. Нажать на "Systems/Inputs", перейти в "Inputs"
3. Выбрать "GELF TCP" и нажать на кнопку "Launch new input"
4. Включить "Global"
5. В поле "Title" вставить желаемое имя инпута (например, checkout или loms)
6. Для checkout оставить порт 12201, для loms изменить на 12202, для notifications изменить на 12203

## Мониторинг трейсов
Jaeger доступен по порту 16686

## Мониторинг метрик в Prometheus + Grafana
1. Перейти на http://localhost:3000
2. Создать новый Data source:
    a. в поле Prometheus server URL указать http://localhost:9090,
    b. scrape interval указать 5s (как в prometheus.yml)
3. Создать новый Dashboard, в качестве data source указать Prometheus
4. Выбрать интересующую метрику
