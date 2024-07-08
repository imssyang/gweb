import traceback
from typing import Optional
from formatify import Command
from formatify import AstText, JsonText


def dumps(mode: str, data: str, indent: int, has_escape: bool) -> Optional[str]:
    try:
        if mode == "json":
            return JsonText(data).dumps(indent, has_escape)
        elif mode == "python":
            return AstText(data).dumps(indent)
        elif mode == "command":
            return Command(data).dumps(indent)
        else:
            print(f"[PYFMT] Unsupport {mode} mode.")
            return None
    except:
        tb_info = traceback.format_exc()
        print(f"[PYFMT] dumps> {tb_info}")
        return None
