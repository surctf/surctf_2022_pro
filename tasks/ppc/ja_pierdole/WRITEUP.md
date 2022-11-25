Подключаемся, видим:

![Screenshot 2022-11-25 at 13 51 40](https://user-images.githubusercontent.com/24609869/203967726-f1280927-05fb-4dd8-b7bb-1e63a2281b9e.png).   
Какой-то текст на польском и строка из математических операторов и чисел. Кидаем в переводчик, понимаем, что нужно вычислить значение какого-то выражения, видимо это строчка из операторов и чисел. Гуглим что-то типа - `Польша математика выражение`, и находим страницу с википедии о [польской записи](https://ru.wikipedia.org/wiki/Польская_запись). Если вам повезет, то первым примером вы получите выражение в [обратной польской записи](https://ru.wikipedia.org/wiki/Обратная_польская_запись), что является тем же самым, но перевернутым. На обоих страницах в википедии есть ссылки на прямую и обратную записи, так что найти их не трудно.  

Читаем, разбираемся как она работает и пробуем вычислить данный пример, разберу вариант со скриншота:   
Имеем `+ - 20 5 + 10 13`, решение:
  1. Читаем слева-направо, первый оператор `+`, значит мы должны сложить 2 выражения `- 20 5` и `+ 10 13`
  2.  `- 20 5` равно `20 - 5 = 15`
  3.  `+ 10 13` равно `10 + 13 = 23`
  4.  Подставляем в выражение из первого шага, получаем `+ 15 23 = 15 + 23 = 38`

Отправляем, в ответ получаем:  
![Screenshot 2022-11-25 at 14 46 55](https://user-images.githubusercontent.com/24609869/203979300-a24f77d7-16ec-471a-9da6-96cf5f43607f.png)

Видим, что сервис предложил нам новое выражение, причем в [обратной польской записи](https://ru.wikipedia.org/wiki/Обратная_польская_запись), можем, конечно попробовать руками порешать, но так-как мы не знаем как много там примеров, то лучше написать скрипт, который всё сделает за нас.

[Скрипт](./solve.py):
```python
from pwn import *

ops = {
  "+": (lambda a, b: a + b),
  "-": (lambda a, b: a - b),
  "*": (lambda a, b: a * b),
}

# Вычисление Постфиксной записи (Обратная польская)
def evalPostfix(expression):
  tokens = expression.split()
  stack = []

  for token in tokens:
    if token in ops:
      arg2 = stack.pop()
      arg1 = stack.pop()
      result = ops[token](arg1, arg2)
      stack.append(result)
    else:
      stack.append(int(token))

  return stack.pop()

# Вычисление Префиксной записи (Прямая польская)
def evalPrefix(expression):
  tokens = expression.split()[::-1] # Отличие от evalPostfix #1, переворачиваем выражение
  stack = []

  for token in tokens:
    if token in ops:
      arg2 = stack.pop()
      arg1 = stack.pop()
      result = ops[token](arg2, arg1) # Отличие от evalPostfix #2, меняем операнды местами
      stack.append(result)
    else:
      stack.append(int(token))

  return stack.pop()


ip = "185.104.115.19"
r = remote(ip, 8877)
r.recvuntil("całkowite!\n".encode("utf-8"))

try:
	while True:
		s = r.recvline().decode("utf-8")
		print(s)
		exp = s[s.index("]") + 2:]
    # Проверяем, если первый символ - число, то запись постфиксная, иначе префиксная
		if exp[0].isnumeric():
			result = evalPostfix(exp)
		else:
			result = evalPrefix(exp)

		r.sendline(str(result).encode("utf-8"))
except:
	r.interactive()
finally:
	r.close()
```
`flag: surctf_polish_notation_invented_by_bobr`
