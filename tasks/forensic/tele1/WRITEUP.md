Скачиваем прикреплённый архив, открываем, видим исходники на Go. Бегло просматриваем их, понимаем что это бот, который байтил людей на отправку флагов. В файле constant.go находим интересную вещь:  
![Screenshot 2022-11-25 at 17 48 04](https://user-images.githubusercontent.com/24609869/204008804-29863eec-3852-405c-b594-6fea8f56a1a0.png)

ID какого-то телеграм чата, пытаемся понять, где его использовали:  
![Screenshot 2022-11-25 at 17 48 59](https://user-images.githubusercontent.com/24609869/204008944-d2bb326e-f11e-471a-ae78-61e09eba9611.png)  
Оказывается в этот чат, отправлялись все присланные боту флаги. Первая зацепка есть, теперь нужно понять, как получить информацию об этом чате. Пошаримся в исходниками чуть внимательнее и найдем папку `.git`:  
![Screenshot 2022-11-25 at 17 51 26](https://user-images.githubusercontent.com/24609869/204010025-52bfcb79-89ef-41a7-832a-4183f30f0f49.png)  
Понимаем, что это гит репозиторий, попробуем порыться в комитах. Для райтапа я буду использовать Github Desktop, но таск можно было решить и через консольную утилиту git.  

Открываем репозиторий в клиенте Github Desktop, заходим в `History` и видим несколько коммитов:  
![Screenshot 2022-11-25 at 17 54 41](https://user-images.githubusercontent.com/24609869/204010451-ede96cfc-7160-44b8-9ef8-05a684508ffe.png)  

Пробегаемся по всем коммитам и в самом первом находим файл token, который далее был удален. В этом файле находим токен телеграм бота:  
![Screenshot 2022-11-25 at 17 55 56](https://user-images.githubusercontent.com/24609869/204010667-37c90df8-176c-4803-874b-ac21dc700b50.png)

## Tele1
Вспоминаем, что для решение **tele1** нам нужно найти юзернейм бота. Бежим читать документацию к Telegram Bot API. Находим метод [/getMe](https://core.telegram.org/bots/api#getme):  
![Screenshot 2022-11-25 at 17 57 51](https://user-images.githubusercontent.com/24609869/204010989-1b4a7269-fbfb-4866-a0c1-1d767cfcda9d.png)  
Пробуем отправить на него запрос. Правильный запрос к API выглядит примерно так:
`https://api.telegram.org/bot<bot_token>/<method>?<query>`  
Где, `<bot_token>` найденный нами токен, `<method>` - вызываемый метод и `<query>` - параметры запроса.  

Отправим запрос к методу [/getMe](https://core.telegram.org/bots/api#getme):  
`curl https://api.telegram.org/bot5631770980:AAExlmzAZ1fbMuWwHd5Oa7oHmrtnvdRMuB8/getMe`  
Получим:  
![Screenshot 2022-11-25 at 18 02 29](https://user-images.githubusercontent.com/24609869/204011749-d91d5ce4-cd11-430a-b496-9b29f9ceef3a.png)

Видим в ответе username бота: `fe34fas23d3_never_find_this_bot`, добавляем к нему `surctf_`, и получаем флаг для tele1.  
`[tele1] flag: surctf_fe34fas23d3_you_never_find_this_bot`

## Tele2
Читаем, что необходимо для решения **tele2**, а это юзернейм создателя телеграм бота. Можем попробовать использовать `anonymouse1337` из истории коммитов в репозиторий, но это оказывается не правильно.  

Вспоминаем про найденный в исходниках ID телеграм чата, в который бот отправлял все найденные флаги. Читаем API, ищем методы который позволили бы нам получить больше информации об этом чате. Находим метод [/getChatAdministrators](https://core.telegram.org/bots/api#getchatadministrators) позволяющий получить администраторов чата:  
![Screenshot 2022-11-25 at 18 10 22](https://user-images.githubusercontent.com/24609869/204013207-4c23b9cc-18af-40ec-9eed-90fd4776e002.png)  

Отправим запрос к [/getChatAdministrators](https://core.telegram.org/bots/api#getchatadministrators):  
`curl https://api.telegram.org/bot5631770980:AAExlmzAZ1fbMuWwHd5Oa7oHmrtnvdRMuB8/getChatAdministrators?chat_id=-875674057`  
Получим:  
![Screenshot 2022-11-25 at 18 12 33](https://user-images.githubusercontent.com/24609869/204013600-f193c422-f20d-4df5-9141-8dc93732aa23.png)  

Видим юзернейм `c00l_h4ck3r_ko1ya`, добавляем к нему подпись `surctf_` и получаем флаг.
`[tele2] flag: surctf_fe34fas23d3_you_never_find_this_bot`


