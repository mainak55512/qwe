# qwe/file_manager.py
import hashlib
import os

def get_file_id(file_path):
    """
    Get the file ID based on the file's content or a unique identifier.
    
    Args:
    file_path (str): The path to the file.
    
    Returns:
    str: A unique file ID.
    """
    # Get the absolute path of the file
    absolute_path = os.path.abspath(file_path)
    
    # Get the file's content
    with open(absolute_path, 'rb') as file:
        content = file.read()
    
    # Generate a unique identifier using the file's content
    file_id = hashlib.sha256(content).hexdigest()
    
    return file_id