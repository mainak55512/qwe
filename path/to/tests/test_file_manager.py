# qwe/tests/test_file_manager.py
import unittest
from unittest.mock import patch
from qwe.file_manager import get_file_id

class TestFileManager(unittest.TestCase):
    def test_get_file_id(self):
        # Test with a file path
        file_path = './some.txt'
        expected_file_id = 'some_unique_id'
        with patch('hashlib.sha256') as mock_sha256:
            mock_sha256.return_value.hexdigest.return_value = expected_file_id
            self.assertEqual(get_file_id(file_path), expected_file_id)
        
        # Test with a different file path
        file_path = 'some.txt'
        expected_file_id = 'some_unique_id'
        with patch('hashlib.sha256') as mock_sha256:
            mock_sha256.return_value.hexdigest.return_value = expected_file_id
            self.assertEqual(get_file_id(file_path), expected_file_id)