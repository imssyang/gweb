import ast
import io
import json
import pprint
from typing import Any
from typing import Optional


class BaseText:
    def is_json(self, data: str) -> bool:
        try:
            json.loads(s=data)
            return True
        except:
            return False

    def json2ast(self, data: str) -> Any:
        try:
            return json.loads(s=data)
        except:
            return None

    def ast2json(self, data: Any) -> Optional[str]:
        try:
            return json.dumps(data, ensure_ascii=False)
        except:
            return None

    def str2ast(self, data: str) -> Any:
        try:
            return ast.literal_eval(data)
        except:
            return None


class JsonText(BaseText):
    def __init__(self, data: Any):
        super().__init__()
        if isinstance(data, str):
            d = self.json2ast(data)
            if d is None:
                d = self.str2ast(data)
                if d is None:
                    self.data = None
                    return
            self.data = d
        else:
            self.data = data

    def dumps(self, indent: int) -> Optional[str]:
        try:
            indent = None if indent <= 0 else indent
            r = json.dumps(
                self.data,
                ensure_ascii=False,
                indent=indent,
            )
            return r if self.data else None
        except:
            return None

class AstText(BaseText):
    def __init__(self, data: Any):
        super().__init__()
        if isinstance(data, str):
            if self.is_json(data):
                self.data = data
            else:
                d = self.str2ast(data)
                if d is None:
                    self.data = None
                    return
                self.data = self.ast2json(d)
        else:
            self.data = self.ast2json(data)

    def dumps(self, indent: int) -> Optional[str]:
        try:
            if indent:
                d = json.loads(s=self.data)
                with io.StringIO() as buf:
                    pp = pprint.PrettyPrinter(
                        indent=indent,
                        compact=False,
                        stream=buf,
                    )
                    pp.pprint(d)
                    return buf.getvalue()
            else:
                d = json.loads(s=self.data)
                return str(d)
        except:
            return None
