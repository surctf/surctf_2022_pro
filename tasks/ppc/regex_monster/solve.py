import exrex
import hashlib

regex = r'surctf_([p1l]{1}[3rk_n]{2}[so1]){3}'

flag_hash = '5f1824397c92c1860e87c987c3684f40'

reg_list = list(exrex.generate(regex))

for match in reg_list:
	if hashlib.md5(match.encode('utf-8')).hexdigest() == flag_hash:
		print(match)
		break
