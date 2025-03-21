## Overview

The repository of loyalty system service

## Endpoints

Method | Path                                | Description                                   |                                                                         
---    |-------------------------------------|------------------------------------------------
GET    | `/health`                           | Health page                                   |
GET    | `/metrics`                          | Страница с метриками                          |
GET    | `/v1/loyalty/list`                  | Получение всех видов лояльностей              |
GET    | `/v1/loyalty/get/{id}`              | Получение лояльности                          |
GET    | `/v1/loyalty/get-for-user/{userId}` | Получение всех лояльностей по пользователю    |