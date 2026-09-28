from datetime import datetime

def log(*args, sep=' ', end='\n'):
    """
    Log messages to console only (stdout) with a timestamp.
    Docker console logs are temporary and size-limited.
    """
    try:
        timestamp = datetime.utcnow().isoformat() + 'Z'
        message = sep.join(str(a) for a in args)
        print(f"{timestamp} {message}", end=end)
    except Exception:
        # best-effort: avoid raising from logging
        pass
