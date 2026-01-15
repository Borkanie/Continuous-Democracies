from datetime import datetime

def file_log(*args, sep=' ', end='\n', alsoPrint=False):
    """
    Log messages to console only (stdout) with a timestamp.
    No file logging to prevent large log files on disk.
    Docker console logs are temporary and size-limited.
    """
    try:
        timestamp = datetime.utcnow().isoformat() + 'Z'
        message = sep.join(str(a) for a in args)
        print(f"{timestamp} {message}", end=end)
    except Exception:
        # best-effort: avoid raising from logging
        pass
