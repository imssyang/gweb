
class OptionKey:
    def __init__(self, option_key: str):
        self.literal = option_key

    def dumps(self, indent: int) -> str:
        return self.literal
