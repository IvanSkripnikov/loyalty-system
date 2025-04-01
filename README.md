## Overview

The repository of loyalty system service

## Endpoints

Method | Path                                   | Description                                     +|                                                                         
------ |----------------------------------------|--------------------------------------------------
GET    | `/health`                              | Health page                                     |
GET    | `/metrics`                             | Страница с метриками                            |
GET    | `/v1/loyalty/list`                     | Получение всех видов лояльностей                |
GET    | `/v1/loyalty/get/{id}`                 | Получение лояльности                            |
GET    | `/v1/loyalty/get-for-user/{userId}`    | Получение всех лояльностей по пользователю      |
PUT    | `/v1/loyalty/apply-for-order`          | Применение доступных лояльностей к заказу       |
POST   | `/v1/loyalty/create`                   | Создание новой лояльности                       |
PUT    | `/v1/loyalty/update`                   | Изменение лояльности                            |
DELETE | `/v1/loyalty/remove/{loyaltyId}`       | Удаление лояльности                             |
DELETE | `/v1/loyalty/remove-for-user/{userId}` | Удаление лояльности у конкретного пользователя  |
DELETE | `/v1/loyalty/remove-certificate`       | Деактивация подарочного сертификата             |
GET    | `/v1/loyalty/configuration/list`       | Получение список настроек системы               |
PUT    | `/v1/loyalty/configuration/update`     | Обновление настроек системы                     |
------ |----------------------------------------|-------------------------------------------------|
GET    | `/test/run-loyalty`                    | Тестовая эмуляция крона применения лояльностей  |
GET    | `/test/remove-loyalty`                 | Тестовая эмуляция крона деактивации лояльностей |