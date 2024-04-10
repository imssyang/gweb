from pycmd import Command

def dumps(cmd, indent):
    try:
        c = Command(cmd)
        r = c.dumps(indent)
        return r
    except Exception as e:
        print(f"PyException: {e}\n")
        return None
