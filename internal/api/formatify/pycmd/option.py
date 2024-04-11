from .key import OptionKey
from .value import OptionValue
from .ffmpeg import FFmpegGuess, FFmpegValue


class Option:
    def __init__(self, key: str, delimiter: str, value: str):
        self.raw_key = key
        self.raw_value = value
        self.delimiter = delimiter
        self.key = OptionKey(key)
        self.value = self.parse_value(key, value)

    def parse_value(self, key: str, value: str) -> OptionValue:
        guess = FFmpegGuess(key, value)
        if guess.ok:
            return FFmpegValue(guess, value)
        return OptionValue(value)

    def dumps(self, indent: int, tail_flag: bool) -> str:
        key = self.key.dumps(indent)
        value = self.value.dumps(indent)
        if indent:
            prefix = '\n' + ' ' * indent
            suffix = '' if tail_flag else ' \\'
        else:
            prefix = ' '
            suffix = ''
        return f'{prefix}{key}{self.delimiter}{value}{suffix}'
