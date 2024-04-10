import ast
import io
import json
import pprint


class BaseText:
    def is_json(self, data):
        try:
            json.loads(s=data)
            return True
        except Exception as e:
            return False

    def json2ast(self, data, pivotal = False):
        try:
            return json.loads(s=data)
        except json.decoder.JSONDecodeError as e:
            return None

    def ast2json(self, data):
        try:
            return json.dumps(data, ensure_ascii=False)
        except Exception as e:
            return None

    def str2ast(self, data, pivotal = False):
        try:
            return ast.literal_eval(data)
        except SyntaxError as e:
            return None


class JsonText(BaseText):
    def __init__(self, data):
        super().__init__()
        if isinstance(data, str):
            d = self.json2ast(data)
            if d is None:
                d = self.str2ast(data, True)
                if d is None:
                    self.data = None
                    return
            self.data = d
        else:
            self.data = data

    def dumps(self, indent: int):
        try:
            return json.dumps(
                self.data,
                ensure_ascii=False,
                indent=None if indent == 0 else indent,
            )
        except Exception as e:
            print(f"[PYTHON] JsonText dumps: {e}")
            return None

class AstText(BaseText):
    def __init__(self, data):
        super().__init__()
        if isinstance(data, str):
            if self.is_json(data):
                self.data = data
            else:
                d = self.str2ast(data, True)
                if d is None:
                    self.data = None
                    return
                self.data = self.ast2json(d)
        else:
            self.data = self.ast2json(data)

    def dumps(self, indent: int):
        try:
            if indent:
                d = json.loads(s=self.data)
                with io.StringIO() as buf:
                    pp = pprint.PrettyPrinter(
                        indent=None if indent == 0 else indent,
                        compact=False,
                        stream=buf,
                    )
                    pp.pprint(d)
                    return buf.getvalue()
            else:
                d = json.loads(s=self.data)
                return str(d)
        except Exception as e:
            print(f"[PYTHON] AstText dumps: {e}")
            return None
