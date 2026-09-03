import os
import re
from typing import (
    AnyStr,
    Tuple,
    Pattern,
)
from .base import Decorator
from .section import Section


class OptionValue:
    def __init__(self, value: str):
        self.sections = []
        self.tail_path = None
        if value:
            prefix, quote, suffix = self.parse_on_quote(value)
            if quote:
                if prefix:
                    self.sections.append(prefix)

                prefix, self.tail_path = self.parse_on_tail_path(quote)
                if self.tail_path:
                    if prefix:
                        self.sections.append(prefix)
                else:
                    self.sections.append(quote)
                    if suffix:
                        prefix, self.tail_path = self.parse_on_tail_path(suffix)
                        if self.tail_path:
                            if prefix:
                                self.sections.append(prefix)
                        else:
                            self.sections.append(suffix)
            else:
                prefix, self.tail_path = self.parse_on_tail_path(value)
                if prefix:
                    self.sections.append(prefix)

    @classmethod
    @Decorator.cache
    def compile_on_quote(cls) -> Pattern[AnyStr]:
        prefix = r"(?P<prefix>.*?)(?=['\"]+)"
        quote = r"(?P<quote>['\"]+.*['\"]+)"
        suffix = r"(?P<suffix>.*)"
        pattern = fr"{prefix}{quote}{suffix}"
        return re.compile(pattern)

    def parse_on_quote(self, value) -> Tuple[str, str, str]:
        pattern = self.compile_on_quote()
        m = pattern.match(value)
        if not m:
            return value, None, None
        prefix = m.group("prefix")
        quote = m.group("quote")
        suffix = m.group("suffix")
        return prefix, quote, suffix

    @classmethod
    @Decorator.cache
    def compile_on_tail_path(cls) -> Pattern[AnyStr]:
        separator = re.escape(os.path.sep)
        flag = separator + r'+'
        prefix = fr"(?P<prefix>.*?)(?=\s+{separator})"
        path = fr"\s+(?P<path>['\"]?{separator}.*)"
        pattern = fr"{prefix}{path}"
        return re.compile(pattern)

    def parse_on_tail_path(self, value: str) -> Tuple[str, str]:
        pattern = self.compile_on_tail_path()
        m = pattern.match(value)
        if not m:
            return value, None
        prefix = m.group("prefix")
        path = m.group("path")
        return prefix, path

    def dumps(self, indent: int) -> str:
        r = str()
        for section in self.sections:
            if isinstance(section, str):
                r += section
            elif isinstance(section, Section):
                r += section.dumps(indent)
        if self.tail_path:
            if indent:
                suffix = ' \\'
                prefix = '\n' + ' ' * indent
            else:
                suffix = ''
                prefix = ' '
            r += f'{suffix}{prefix}{self.tail_path}'
        return r
