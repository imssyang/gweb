import unittest
from formatify import Command


class PycmdTestCase(unittest.TestCase):
    FFPROBE_CMD = 'ffprobe -hide_banner -select_streams v:0 -of json 123.mp4'

    def test_command(self):
        self.assertEqual(
            Command(self.FFPROBE_CMD).dumps(0),
            'ffprobe -hide_banner -select_streams v:0 -of json 123.mp4',
            'incorrect return value',
        )

        self.assertEqual(
            Command(self.FFPROBE_CMD).dumps(2),
            'ffprobe \\\n  -hide_banner \\\n  -select_streams v:0 \\\n  -of json 123.mp4',
            'incorrect return value',
        )

