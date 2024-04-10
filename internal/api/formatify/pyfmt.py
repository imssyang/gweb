from pycmd import Command
from pytext import AstText, JsonText


def dumps(mode, data, indent):
    try:
        if mode == "json":
            return JsonText(data).dumps(indent)
        elif mode == "python":
            return AstText(data).dumps(indent)
        elif mode == "command":
            return Command(data).dumps(indent)
        else:
            print(f"[PYTHON] Unsupport {mode} mode.")
            return None
    except Exception as e:
        print(f"[PYTHON] dumps: {e}")
        return None
