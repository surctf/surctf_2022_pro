from pwn import *

ops = {
  "+": (lambda a, b: a + b),
  "-": (lambda a, b: a - b),
  "*": (lambda a, b: a * b),
}

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

def evalPrefix(expression):
  tokens = expression.split()[::-1]
  stack = []

  for token in tokens:
    if token in ops:
      arg2 = stack.pop()
      arg1 = stack.pop()
      result = ops[token](arg2, arg1)
      stack.append(result)
    else:
      stack.append(int(token))

  return stack.pop()


ip = "0.0.0.0"
r = remote(ip, 8877)
r.recvuntil("całkowite!\n".encode("utf-8"))

try:
	while True:
		s = r.recvline().decode("utf-8")
		print(s)
		exp = s[s.index("]") + 2:]
		if exp[0].isnumeric():
			result = evalPostfix(exp)
		else:
			result = evalPrefix(exp)

		r.sendline(str(result).encode("utf-8"))
except:
	r.interactive()
finally:
	r.close()