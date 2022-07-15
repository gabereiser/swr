import random
import time

def main():
    alpha = "abcdefghijklmnopqrstuvwxyz"
    ret = ""
    for i in range(26):
        ret += alpha[random.randint(0,25)]
        time.sleep(0.1)
    ret += ret.upper()
    print(ret)

if __name__ == "__main__":
    main()
