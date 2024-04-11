import os
import re
from typing import (
    AnyStr,
    List,
    Tuple,
    Pattern,
)
from .base import Decorator
from .value import OptionValue
from .section import Section


class FFmpegGuess:
    def __init__(self, key: str, value: str):
        self.is_filter = self.__guess(key, value, ['filter_complex'])
        self.is_codec = self.__guess(key, value, ['-params', '_params'])

    def __guess(self, key: str, value: str, flags: List[str]) -> bool:
        for flag in flags:
            if flag in key.lower() or flag in value.lower():
                return True
        return False

    @property
    def ok(self) -> bool:
        return self.is_filter or self.is_codec


class FFmpegFilterSection(Section):
    def __init__(self, section: str):
        super().__init__(section)

    def parse(self, section: str) -> List[str]:
        items = []
        caches = []
        segments = re.split(r"(['\";])", section)
        for segment in segments:
            caches.append(segment)
            if ';' in segment:
                items.append(''.join(caches).strip())
                caches.clear()
        if caches:
            items.append(''.join(caches).strip())
        return items

    def dumps(self, indent: int) -> str:
        r = str()
        for i, item in enumerate(self.items):
            if indent:
                head_flag = bool(i == 0)
                tail_flag = bool(i == len(self.items) - 1)
                prefix = '' if head_flag else '\n' + ' ' * indent * 2
                suffix = '' if tail_flag else ' \\'
            else:
                prefix = suffix = ''
            r += f'{prefix}{item}{suffix}'
        return r


class FFmpegCodecSection(Section):
    def __init__(self, section: str):
        super().__init__(section)

    def parse(self, section: str) -> List[str]:
        items = []
        caches = []
        segments = re.split(r"(['\":])", section)
        for segment in segments:
            caches.append(segment)
            if ':' in segment:
                items.append(''.join(caches).strip())
                caches.clear()
        if caches:
            items.append(''.join(caches).strip())
        return items

    def dumps(self, indent: int) -> str:
        r = str()
        for i, item in enumerate(self.items):
            if indent:
                head_flag = bool(i == 0)
                tail_flag = bool(i == len(self.items) - 1)
                prefix = '' if head_flag else '\n'
                suffix = '' if tail_flag else '\\'
            else:
                prefix = suffix = ''
            r += f'{prefix}{item}{suffix}'
        return r


class FFmpegValue(OptionValue):
    def __init__(self, guess: FFmpegGuess, value: str):
        super().__init__(value)
        self.__update_sections(guess)

    def __update_sections(self, guess: FFmpegGuess):
        for i, section in enumerate(self.sections):
            if guess.is_filter:
                self.sections[i] = FFmpegFilterSection(section)
            elif guess.is_codec:
                self.sections[i] = FFmpegCodecSection(section)

