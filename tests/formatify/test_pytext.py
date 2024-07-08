import unittest
from formatify import JsonText, AstText


class PytextTestCase(unittest.TestCase):
    SAMPLE_JSON = '{"a":true,"b":[1,2,3]}'
    SAMPLE_DICT = "{'a': True, 'b': [1, 2, 3]}"
    SAMPLE_ESCAPE = '"{\\"a\\": true, \\"b\\": [1, 2, 3]}"'
    FFPROBE_CMD = 'ffprobe -hide_banner -select_streams v:0 -show_packets -of json 123.mp4'

    def test_jsontext(self):
        self.assertEqual(
            JsonText(self.SAMPLE_JSON).dumps(0),
            '{"a": true, "b": [1, 2, 3]}',
            'incorrect return value',
        )

        self.assertEqual(
            JsonText(self.SAMPLE_DICT).dumps(0),
            '{"a": true, "b": [1, 2, 3]}',
            'incorrect return value',
        )

        self.assertEqual(
            JsonText(self.SAMPLE_JSON).dumps(0, True),
            '"{\\"a\\": true, \\"b\\": [1, 2, 3]}"',
            'incorrect return value',
        )

        self.assertEqual(
            JsonText(self.SAMPLE_DICT).dumps(0, True),
            '"{\\"a\\": true, \\"b\\": [1, 2, 3]}"',
            'incorrect return value',
        )

        self.assertEqual(
            JsonText(self.SAMPLE_ESCAPE).dumps(0),
            '{"a": true, "b": [1, 2, 3]}',
            'incorrect return value',
        )

        self.assertEqual(
            JsonText(self.FFPROBE_CMD).dumps(0),
            None,
            'incorrect return value',
        )

    def test_asttext(self):
        self.assertEqual(
            AstText(self.SAMPLE_JSON).dumps(0),
            "{'a': True, 'b': [1, 2, 3]}",
            'incorrect return value',
        )
