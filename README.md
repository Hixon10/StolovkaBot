# StolovkaBot

* Уведомляет вас в Телеграмме, когда ваша еда готова. Речь идёт о Яндекс столовой в БЦ Бенуа (СПБ).
* Данные берутся из http://kotikicanteen.ru/orders

## Запуск:
* docker build -t stolovkabot .
* docker run --restart=always -itd --name stolovkabot stolovkabot:latest