from pytext import JsonText, AstText


if __name__ == "__main__":
    j = '{"a":True,"b":[1,2,3]}'
    #j = '{"name": "Alice", "age": {"a": 1, "b": 2}, "phone": "abc,}'
    print(JsonText(j).pretty())
    print(JsonText(j).compact())
    print(AstText(j).pretty())
    print(AstText(j).compact())
