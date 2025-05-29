# Balance Checker

Программа для проверки баланса через API yarurf.ru

## Требования
- Go 1.16+
- Файл конфигурации `config.yaml`

## Установка
```bash
go get github.com/sirupsen/logrus gopkg.in/yaml.v3
go build -o balance-checker
```
## Конфигурация

Создайте файл config.yaml:
```yaml
user:
  login: "ваш_логин"
  password: "ваш_пароль"
templates:
  high_balance: "Баланс: %.2f (всё хорошо)"
  low_balance: "Внимание! Баланс: %.2f (ниже нормы)"

trigger: 100.00
```

## Использование
###  Базовый запуск
```bash
./balance-checker
```
### С указанием конфига
```bash
./balance-checker -config /path/to/config.yaml
```
## Примеры форматов вывода

  - Если balance > trigger → high_balance шаблон
  - Если balance ≤ trigger → low_balance шаблон

    Баланс: 150.00 (всё хорошо)  # При balance > 100
    
    Внимание! Баланс: 80.00 (ниже нормы)  # При balance ≤ 100

## Флаги

    -config - путь к файлу конфигурации (по умолчанию config.yaml)

## Логика работы
- Проверяет валидность текущей сессии по cookies
- Если сессия недействительна - выполняет авторизацию
- Делает запрос баланса
- Выводит результат в указанном формате
- Завершает работу