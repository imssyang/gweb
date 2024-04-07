from pycmd import Command

def getInteger():
    print('Python function getInteger() called')
    c = 100*50/30
    return c

def formatCommand():
    try:
        c = Command("/sss/test -a 1 -b 2 -c 3 adsf")
        return c.dumps(2)
    except Exception as e:
        print(f"PyException: {e}\n")
        return None

def formatCommand2(data):
    try:
        c = Command(data)
        r = c.dumps(2)
        print("Python: ", r)
        return r
    except Exception as e:
        print(f"PyException: {e}\n")
        return None

#print(formatCommand())
#print(formatCommand2("/sss/test -a 1 -b 2 -c 3 adsf"))