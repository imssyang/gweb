from typing import List


class Section:
    def __init__(self, section: str):
        self.items = self.parse(section)

    def parse(self, section: str) -> List[str]:
        raise NotImplementedError()

    def dumps(self, indent: int) -> str:
        raise NotImplementedError()
