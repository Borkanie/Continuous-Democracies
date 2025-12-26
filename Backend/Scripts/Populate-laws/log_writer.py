import os
from datetime import datetime

OUTPUT_PATH = os.path.join(os.getcwd(), "output.txt")

def file_log(*args, sep=' ', end='\n', alsoPrint = False):
    """Append a message to the output file with a timestamp."""
    try:
        if alsoPrint:
            print(sep.join(str(a) for a in args))
        with open(OUTPUT_PATH, 'a', encoding='utf-8') as f:
            timestamp = datetime.utcnow().isoformat() + 'Z'
            f.write(timestamp + ' ')
            f.write(sep.join(str(a) for a in args))
            f.write(end)
    except Exception:
        # best-effort: avoid raising from logging
        pass
