import re
from typing import (
    AnyStr,
    List,
    Tuple,
    Pattern,
)
from .base import Decorator
from .option import Option


class Command:
    def __init__(self, cmd: str):
        self.cmd = self.adjust_cmd(cmd) #cmd.replace("\\\n", "")
        self.program, self.remainder = self.get_program(self.cmd)
        self.options = self.get_options(self.remainder)

    def adjust_cmd(self, cmd: str) -> str:
        return cmd.replace("\\\n", "")

    @classmethod
    @Decorator.cache
    def compile_program(cls) -> Pattern[AnyStr]:
        program = r"\s*(?P<program>.*?)(?=\s+-{1,2})"
        remainder = r"(?P<remainder>.*)"
        pattern = fr"{program}{remainder}"
        return re.compile(pattern)

    def get_program(self, cmd) -> Tuple[str, str]:
        pattern = self.compile_program()
        m = pattern.match(cmd)
        if not m:
            return None, None
        program = m.group("program")
        remainder = m.group("remainder")
        if program and program[0] == '-':
            remainder = ' ' + program + remainder
            program = None
        return program, remainder

    def get_options(self, remainder: str) -> Tuple[List[Option], str]:
        options = []
        contexts, remainder = self.get_option_contexts(remainder)
        for context in contexts:
            key, delimiter, value = self.get_middle_option(context)
            if key:
                options.append(Option(key, delimiter, value))
        key, delimiter, value = self.get_last_option(remainder)
        if key:
            options.append(Option(key, delimiter, value))
        return options

    def get_option_contexts(self, remainder: str) -> Tuple[List[str], str]:
        contexts = []
        while True:
            option, remainder_tmp = self.next_option_context(remainder)
            if not option:
                break

            contexts.append(option)
            remainder = remainder_tmp
        return contexts, remainder

    @classmethod
    @Decorator.cache
    def compile_option_context(cls) -> Pattern[AnyStr]:
        option = r"\s+(?P<option>-{1,2}\S+.*?)(?=\s+-{1,2}\S+)"
        remainder = r"(?P<remainder>.*)"
        pattern = fr"{option}{remainder}"
        return re.compile(pattern)

    def next_option_context(self, remainder: str) -> Tuple[str, str]:
        pattern = self.compile_option_context()
        m = pattern.match(remainder)
        if not m:
            return None, None
        option = m.group("option")
        remainder = m.group("remainder")
        return option, remainder

    @classmethod
    @Decorator.cache
    def compile_middle_option(cls) -> Pattern[AnyStr]:
        key = r"\s*(?P<key>-{1,2}[^=\s]+.*?)(?=[=\s]*)"
        delimiter = r"(?P<delimiter>[=\s]*)"
        value = r"(?P<value>.*)"
        pattern = fr"{key}{delimiter}{value}"
        return re.compile(pattern)

    def get_middle_option(self, context: str) -> Tuple[str, str, str]:
        pattern = self.compile_middle_option()
        m = pattern.match(context)
        if not m:
            return None
        key = m.group("key")
        delimiter = m.group("delimiter")
        value = m.group("value")
        return key, delimiter, value

    @classmethod
    @Decorator.cache
    def compile_last_option(cls) -> Pattern[AnyStr]:
        key = r"\s+(?P<key>-{1,2}\S+.*?)(?=\s*)"
        delimiter = r"(?P<delimiter>\s*)"
        value = r"(?P<value>.*)(?=\s*)"
        pattern = fr"{key}{delimiter}{value}"
        return re.compile(pattern)

    def get_last_option(self, remainder: str) -> Tuple[str, str, str]:
        pattern = self.compile_last_option()
        m = pattern.match(remainder)
        if not m:
            return None, None, None
        key = m.group("key")
        delimiter = m.group("delimiter")
        value = m.group("value")
        return key, delimiter, value

    def dumps(self, indent: int) -> str:
        cmd = str()
        if self.program:
            cmd += self.program
            if indent:
                cmd += ' \\'
        for i, option in enumerate(self.options):
            tail_flag = bool(i == len(self.options) - 1)
            cmd += option.dumps(indent, tail_flag)
        return cmd
