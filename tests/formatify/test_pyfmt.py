import unittest
from formatify import dumps


class PyfmtTestCase(unittest.TestCase):
    SAMPLE_JSON = '{"a":True,"b":[1,2,3]}'
    SAMPLE_DICT = "{'a': True, 'b': [1, 2, 3]}"
    FFPROBE_CMD = 'ffprobe -hide_banner -select_streams v:0 -show_packets -of json 123.mp4'

    def test_dumps(self):
        self.assertEqual(
            dumps('json', self.SAMPLE_JSON, 0),
            '{"a": true, "b": [1, 2, 3]}',
            'incorrect return value',
        )

        self.assertEqual(
            dumps('python', self.SAMPLE_DICT, 0),
            "{'a': True, 'b': [1, 2, 3]}",
            'incorrect return value',
        )

        self.assertEqual(
            dumps('python', self.FFPROBE_CMD, 0),
            None,
            'incorrect return value',
        )
