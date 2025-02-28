## Описание сервиса
* Сервис позволяет залогиниться (создается новый пользователь, если ранее не существовал);
* Доступ ко всем методам кроме `api/auth` не авторизован пользователям без JWT токена;
* При создании пользователя начисляется 1000 монет, за которые можно купить мерч (в миграции предусмотрено заполнение бд мерча), также можно отправить монеты другому пользователю;
* При запросе `api/info` пользователь получит сгруппированную информацию о перемещении его монеток формата:
```json
{
  "inventory": [
    {
      "type": "book",
      "quantity": 3
    },
    {
      "type": "cup",
      "quantity": 3
    },
    {
      "type": "powerbank",
      "quantity": 1
    }
  ],
  "coins": "190",
  "coinHistory": {
    "received": [
      {
        "user": "pasha",
        "amount": "200"
      }
    ],
    "sent": [
      {
        "user": "marina",
        "amount": "300"
      },
      {
        "user": "pasha",
        "amount": "300"
      }
    ]
  }
}
```
* В проекте предусмотрен HTTP сервер (**8080**) и gRPC сервер (**8081**);
* Мерч при запуске сервиса выгружается в кэш, что снижает нагрузку на бд и увеличивает производительность;
* Покрытие сервиса тестами составляет **45%**
* Реализовано E2E тестирование всех эндпоинтов

***
## Как пользоваться
1. Склонировать репозиторий
2. Находясь в корневом каталоге проекта поднять докер: `docker-compose up -d` или `docker-compose up`
<br>
Готово, бд и сервис подняты :)

***
Чтобы проверить работу сервиса (используя Postman) через gRPC:
<br>
### Адрес: `localhost:8081`
1. Создать новый запрос;
2. Выбрать протокол запроса gRPC;
3. Далее, при выборе метода, нужно выбрать **Use server reflection**;
4. Отправить запрос методом `Authenticate`, введя сообщение в формате 
```json
{
    "password": "pass",
    "username": "user"
}
```
<br>

Или можно воспользоваться кнопкой `Use example message` (также относится ко всем последующим методам);
5. Скопировать токен без кавычек;
<br>
6. Теперь можно вызывать любые методы, передавая в metadata значение `authorization: your_token`;

***
Чтобы проверить работу сервиса (используя Postman) через HTTP:
<br>
### Адрес:`localhost:8080/api`
1. Создать новый запрос;
2. Выбрать протокол HTTP;
3. Вызвать метод `/auth` с глаголом запроса **POST**, ввести данные в формате:
```json
{
    "password": "pass",
    "username": "user"
}
```
4. Скопировать токен без кавычек;
5. Теперь можно вызывать любые методы, передавая в headers значение `authorization: your_token`;

### Все доступные для вызова HTTP запросы описаны и полностью соответствуют OpenAPI схеме из ТЗ.

***
### Мои контакты:
**VK:** https://vk.com/jabrimuhi
<br>
**Telegram:** https://t.me/JabriMuhi
<br>
**EMail:** jabri.muhi@yandex.ru
