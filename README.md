<h3 align="center">Yandex Distributed Calculator 2</h3>

  <p align="center">
    Распределённый калькулятор для вычисления выражений, вторая версия
    <br />
    <br />
    <a href="https://github.com/Uikola/distributedCalculator2">View Demo</a>
    ·
    <a href="https://t.me/uikola">Report Bug</a>
  </p>

<!-- ABOUT THE PROJECT -->
## О Проекте
Вторая версия распределённого калькулятора(ссылка на первую версию: https://github.com/Uikola/distributedCalculator). Новая версия работает в разы быстрее благодаря коммуникации сервисов через grpc. Появился фронтенд, поэтому можно будет опробовать проект тыкая кнопочки)(в прошлой версии был свагер). В разы упрощен запуск приложения.

### Использованные Технологии

- PostgresSQL
- grpc
- Chi router
- Docker
- Golang(1.22.0)

<!-- GETTING STARTED -->
## Начало Работы

Чтобы запустить приложение следуйте следующим шагам.

### Установка

1. Клонируйте репозиторий
   ```sh
   git clone https://github.com/Uikola/distributedCalculator2.git
   ```

2. Запустите docker-compose:
   ```sh
   docker-compose up -d
   ```

3. В корневой папке запустите go mod tidy
   ```sh
   go mod tidy
   ```

4. Последовательно запустите оркестратор и фронтенд
   ```sh
   go run orchestrator/cmd/app/main.go
   go run frontend/cmd/main.go
   ```

5. Приложение готово к использованию и распологается на http://localhost:8000/html. Для запуска вычислительного ресурса используйте go run calculator/cmd/app/main.go. Вы можете запустить их сколько угодно!
   ```sh
   go run calculator/cmd/app/main.go
   ```

6. Интерфейс предоставляет следующий функцонал. Регистарция аккаунта и вход в аккаунт. На гланой странице вы можете просматривать список выражений с результатами, а также добавлять новые(если вычислительных ресурсов недостаточно, то выводится ошибка сервера). В меню вычислительые ресуры вы можете просматривать все вычислительные ресурсы в формате: "Адрес": "Задача, выполняемая на нём". В меню операции вы можете просматривать и менять время операций.

<!-- CONTACT -->
## Contact(Если возникли вопросы)

Yuri - [@telegram](https://t.me/uikola) - ugulaev806@yandex.ru

Project Link: [https://github.com/Uikola/distributedCalculator2](https://github.com/Uikola/distributedCalculator2)