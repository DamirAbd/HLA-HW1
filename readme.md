## Домашнее задание для курса Otus Highload Arhitect

### Особенности


Все доступы в открытом виде в коде.

- есть возможность авторизации, регистрации, получения анкет по ID;
- golang database/sql защищает от SQL-инъекций;
- пароли пользователей хранятся в хешированном виде;
- монолит с прямым обращением к базе без ORM.

База - одна таблица с полями:
*  userID serial NOT NULL,
*  ID VARCHAR(255) NOT NULL,
*  FirstName VARCHAR(255) NOT NULL,
*  SecondName VARCHAR(255) NOT NULL,
*  BirthDate DATE NOT NULL,
*  Biography VARCHAR(255) NULL,
*  City VARCHAR(255) NULL,
*  Password VARCHAR(255) NOT NULL,
*  CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

PrimaryKey на userID.

### Запуск и работа

Нужен образ postgres:12.19-bullseye
Запуск: docker compose up, слушает порт 8080.

3 метода:
- /login
- /user/register
- /user/get/{id}

[POSTMAN collection](OTUS-HighLoadArch-HW1.postman_collection.json)