# qwe/tests/test_cli.py
import unittest
from unittest.mock import patch
from qwe.cli import handle_file

class TestCLI(unittest.TestCase):
    def test_handle_file(self):
        # Test with a file path
        file_path = './some.txt'
        with patch('qwe.file_manager.get_file_id') as mock_get_file_id:
            mock_get_file_id.return_value = 'some_unique_id'
            handle_file(file_path)
        
        # Test with a different file path
        file_path = 'some.txt'
        with patch('qwe.file_manager.get_file_id') as mock_get_file_id:
            mock_get_file_id.return_value = 'some_unique_id'
            handle_file(file_path)